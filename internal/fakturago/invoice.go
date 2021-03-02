package fakturago

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/pkg/errors"
)

// Invoice holds dependencies need to construct an invoice
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
	paymentTmpl = `{{ t "bankTransfer" | title }}
{{ t "accountNumber" | title }}: {{ .AccountNumber }}
{{ t "dueDate" | title }}: {{ .DueDate.Format "2/01/2006" }}`
)

func (inv *Invoice) addHeader(info BillingInfo) {
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

func (inv *Invoice) addCompaniesInfo(info BillingInfo) {
	billToText, _ := inv.renderTemplate(companyTmpl, info.BillTo)
	companyText, _ := inv.renderTemplate(companyTmpl, info.Company)

	inv.doc.Row(40, func() {
		inv.doc.Col(4, func() {
			inv.doc.SubTitle(strings.Title(inv.loc.T("billTo") + ":"))
			inv.doc.BaseText(billToText, props.Text{Top: 5})
		})

		inv.doc.ColSpace(4)

		inv.doc.Col(4, func() {
			inv.doc.SubTitle(strings.Title(inv.loc.T("seller") + ":"))
			inv.doc.BaseText(companyText, props.Text{Top: 5})
		})
	})

	inv.doc.Line(10)
}

func (inv *Invoice) addItems(info BillingInfo) {
	var (
		headers  []string
		contents [][]string
	)

	enabledColumns := []bool{true, !info.NoTax, true}

	columnsWidth := []uint{8, 2, 2}
	if info.NoTax {
		columnsWidth = []uint{10, 2}
	}

	for _, header := range []string{"name", "tax", "amount"} {
		headers = append(headers, strings.Title(inv.loc.T(header)))
	}

	headers = filterRow(headers, enabledColumns)

	lang := inv.loc.Lang()
	total := 0.0
	for _, item := range info.Items {
		total += item.FinalAmount()

		row := []string{
			item.LocalizedName(lang),
			item.Tax.String(),
			info.Currency.Format(item.FinalAmount()),
		}
		contents = append(contents, filterRow(row, enabledColumns))
	}

	row := []string{
		strings.ToUpper(inv.loc.T("total")),
		"",
		info.Currency.Format(total),
	}

	contents = append(contents, filterRow(row, enabledColumns))

	inv.doc.DataTable(headers, contents, columnsWidth)
}

func (inv *Invoice) addPaymentInfo(info BillingInfo) {
	paymentText, _ := inv.renderTemplate(paymentTmpl, info.Payment)

	inv.doc.Row(20, func() {
		inv.doc.Col(12, func() {
			inv.doc.SubTitle(strings.ToUpper(inv.loc.T("payment")))
			inv.doc.BaseText(paymentText, props.Text{Top: 6})
		})
	})
}

func filterRow(row []string, filter []bool) []string {
	result := []string{}
	for i, enabled := range filter {
		if enabled {
			result = append(result, row[i])
		}
	}
	return result
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

// Generate creates and saves invoice to pdf based on passed BillingInfo
func Generate(info BillingInfo, path string, lang string) error {
	var err error

	bundle, err := loadLanguageBundle("i18n")
	if err != nil {
		return err
	}
	loc := NewLocalizer(bundle, lang)

	doc := newDocument()

	inv := Invoice{doc, loc}

	inv.addHeader(info)
	inv.addCompaniesInfo(info)
	inv.addItems(info)
	inv.addPaymentInfo(info)

	err = inv.doc.OutputFileAndClose(path)

	if err != nil {
		return err
	}

	return nil
}
