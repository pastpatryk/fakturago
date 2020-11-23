package fakturago

import (
	"strings"

	"github.com/imdario/mergo"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/jung-kurt/gofpdf"
)

// Document wraps and extends creation of PDF document
type Document struct {
	*pdf.PdfMaroto
}

func newDocument() Document {
	pdfMaroto := pdf.NewMaroto(consts.Portrait, consts.A4).(*pdf.PdfMaroto)
	addFonts(pdfMaroto)
	return Document{pdfMaroto}
}

func addFonts(pdfMaroto *pdf.PdfMaroto) {
	fpdf := pdfMaroto.Pdf.(*gofpdf.Fpdf)
	fpdf.AddUTF8Font("Lato", "", "./assets/fonts/Lato-Regular.ttf")
	fpdf.AddUTF8Font("Lato", "B", "./assets/fonts/Lato-Bold.ttf")
}

// BaseText creates default style of text
func (d *Document) BaseText(text string, textProps ...props.Text) {
	d.styledText(text, props.Text{Size: 10}, textProps...)
}

// Title creates big headers
func (d *Document) Title(text string, textProps ...props.Text) {
	d.styledText(text, props.Text{
		Size:  12,
		Style: consts.Bold,
		Align: consts.Center,
	}, textProps...)
}

// SubTitle creates smaller headers
func (d *Document) SubTitle(text string, textProps ...props.Text) {
	d.styledText(text, props.Text{
		Size:  10,
		Style: consts.Bold,
	}, textProps...)
}

// DataTable creates styled table with passed headers and rows
func (d *Document) DataTable(headers []string, contents [][]string, columnsWidth []uint) {
	d.TableList(headers, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Family:     "Lato",
			GridSizes:  columnsWidth,
			Size:       12,
			IsUTF8Font: true,
		},
		ContentProp: props.TableListContent{
			Family:      "Lato",
			GridSizes:   columnsWidth,
			IsUTF8Font:  true,
			CellPadding: 2,
		},
		LastRowBackground: &color.Color{
			Red:   200,
			Green: 200,
			Blue:  200,
		},
		HeaderContentSpace: 3,
		Line:               true,
	})
}

func (d *Document) styledText(text string, baseProps props.Text, textProps ...props.Text) {
	defaultProps := props.Text{
		Family:     "Lato",
		IsUTF8Font: true,
	}
	textProps = append([]props.Text{defaultProps, baseProps}, textProps...)
	prop := combineProps(baseProps, textProps...)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		prop.Top += 5
		d.Text(line, prop)
	}
}

func combineProps(base props.Text, textProps ...props.Text) props.Text {
	for _, prop := range textProps {
		if err := mergo.Merge(&base, prop, mergo.WithOverride); err != nil {
			return base
		}
	}
	return base
}
