package fakturago

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

type Invoice struct {
	doc Document
	loc *i18n.Localizer
}

func (inv *Invoice) addHeader() {
	inv.doc.Row(10, func() {
		inv.doc.Col(12, func() {
			inv.doc.Title(strings.ToUpper(t(inv.loc, "Invoice")))
		})
	})

	inv.doc.Line(10)
}

func (inv *Invoice) addCompaniesInfo() {
	inv.doc.Row(40, func() {
		inv.doc.Col(4, func() {
			inv.doc.SubTitle(strings.Title(t(inv.loc, "BillTo") + ":"))
			companyInfo := []string{
				"John Snow",
				"Tower 1",
				"12-345 Winterfell",
				"Westeros",
			}
			for i, line := range companyInfo {
				inv.doc.BaseText(line, props.Text{Top: 5 * float64(i+1)})
			}
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

	loc, err := setupLocalizer(language.English.String())
	if err != nil {
		return err
	}

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

func setupLocalizer(lang string) (*i18n.Localizer, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	err := walkFilesWithExt("locales", ".toml", func(path string) error {
		log.Debug("Loading language: ", path)
		_, err := bundle.LoadMessageFile(path)
		if err != nil {
			return errors.WithMessagef(err, "language %s", path)
		}
		return nil
	})
	if err != nil {
		log.Error("Loading failed! ", err.Error())
		return nil, err
	}

	return i18n.NewLocalizer(bundle, lang), nil
}

func walkFilesWithExt(root, ext string, walkFn func(path string) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ext {
			return nil
		}
		return walkFn(path)
	})
}

func t(loc *i18n.Localizer, key string) string {
	return loc.MustLocalize(&i18n.LocalizeConfig{MessageID: key})
}
