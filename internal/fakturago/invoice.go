package fakturago

import (
	"bytes"
	"html/template"
	"strings"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type Invoice struct {
	doc Document
	loc Localizer
}

const (
	headerTmpl = `{{ t "date" | title }}: {{ .Date.Format "2/01/2006" }}
{{ t "invoiceNumber" | title }}: {{ .Number }}`
	companyTmpl = `{{ .Name }}
{{ .Address }}
{{ .ZipCode }} {{ .City }}
{{ .Country }}
{{ t "vatNumber" | title }}: {{ .VatNumber }}`
)

func (inv *Invoice) addHeader() {
	info := BillingInfo{Date: time.Now(), Number: "01/11/2020"}
	headerText, _ := inv.renderTemplate(headerTmpl, info)

	inv.doc.Row(8, func() {
		inv.doc.Col(12, func() {
			inv.doc.Title(strings.ToUpper(inv.loc.T("invoice")))
		})
	})
	inv.doc.Row(12, func() {
		inv.doc.Col(12, func() {
			inv.doc.SubTitle(headerText, props.Text{Align: consts.Right})
		})
	})

	inv.doc.Line(10)
}

func (inv *Invoice) addCompaniesInfo() {
	company := Company{
		Name:      "John Snow Co.",
		Address:   "South Tower 1",
		City:      "Winterfell",
		ZipCode:   "12-345",
		Country:   "Westeros",
		VatNumber: "NS 01234567",
	}
	companyText, _ := inv.renderTemplate(companyTmpl, company)

	inv.doc.Row(40, func() {
		inv.doc.Col(4, func() {
			inv.doc.SubTitle(strings.Title(inv.loc.T("billTo") + ":"))
			inv.doc.BaseText(companyText, props.Text{Top: 5})
		})

		inv.doc.ColSpace(4)

		inv.doc.Col(4, func() {
			inv.doc.SubTitle(strings.Title(inv.loc.T("seller") + ":"))
			inv.doc.BaseText(companyText, props.Text{Top: 5})
		})
	})

	inv.doc.Line(10)
}

func (inv *Invoice) addItems() {
	headers := []string{
		strings.Title(inv.loc.T("name")),
		strings.Title(inv.loc.T("amount")),
	}
	contents := [][]string{
		{"Software development (01.11.2020 - 31.11.2020)", "35 000 PLN"},
		{"Tech support", "5 000 PLN"},
		{strings.ToUpper(inv.loc.T("total")), "40 000 PLN"},
	}
	inv.doc.DataTable(headers, contents)
}

func (inv *Invoice) renderTemplate(tmplText string, data interface{}) (string, error) {
	var (
		err error
		buf bytes.Buffer
	)
	funcMap := template.FuncMap{
		"t":     inv.loc.T,
		"title": strings.Title,
	}

	tmpl, err := template.New("").Funcs(funcMap).Parse(tmplText)
	if err != nil {
		return "", errors.Wrap(err, "parse template")
	}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", errors.Wrap(err, "execute template")
	}
	return buf.String(), nil
}

func Generate(path string) error {
	var err error

	bundle, err := loadLanguageBundle("i18n")
	if err != nil {
		return err
	}
	loc := NewLocalizer(bundle, language.Polish.String())

	doc := newDocument()

	inv := Invoice{doc, loc}

	inv.addHeader()
	inv.addCompaniesInfo()
	inv.addItems()

	err = inv.doc.OutputFileAndClose(path)

	if err != nil {
		return err
	}

	return nil
}
