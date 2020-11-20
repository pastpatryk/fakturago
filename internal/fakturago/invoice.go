package fakturago

import (
	"strings"

	"github.com/johnfercher/maroto/pkg/props"
	"golang.org/x/text/language"
)

type Invoice struct {
	doc Document
	loc Localizer
}

func (inv *Invoice) addHeader() {
	inv.doc.Row(10, func() {
		inv.doc.Col(12, func() {
			inv.doc.Title(strings.ToUpper(inv.loc.T("Invoice")))
		})
	})

	inv.doc.Line(10)
}

func (inv *Invoice) addCompaniesInfo() {
	inv.doc.Row(40, func() {
		inv.doc.Col(4, func() {
			inv.doc.SubTitle(strings.Title(inv.loc.T("BillToo") + ":"))
			companyInfo := `John Snow
Tower 1
12-345 Winterfell
Westeros`
			inv.doc.BaseText(companyInfo, props.Text{Top: 5})
		})

		inv.doc.ColSpace(4)

		inv.doc.Col(4, func() {
			inv.doc.Text("Test", props.Text{
				Size: 10,
			})
		})
	})
}

func GenerateInvoice(path string) error {
	var err error

	bundle, err := loadLanguageBundle("locales")
	if err != nil {
		return err
	}
	loc := NewLocalizer(bundle, language.English.String())

	doc := newDocument()

	inv := Invoice{doc, loc}

	inv.addHeader()
	inv.addCompaniesInfo()

	err = inv.doc.OutputFileAndClose(path)

	if err != nil {
		return err
	}

	return nil
}
