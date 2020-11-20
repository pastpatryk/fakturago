package fakturago

import (
	"strings"

	"github.com/imdario/mergo"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/jung-kurt/gofpdf"
)

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

func (d *Document) BaseText(text string, textProps ...props.Text) {
	d.styledText(text, props.Text{Size: 10}, textProps...)
}

func (d *Document) Title(text string, textProps ...props.Text) {
	d.styledText(text, props.Text{
		Size:  12,
		Style: consts.Bold,
		Align: consts.Center,
	}, textProps...)
}

func (d *Document) SubTitle(text string, textProps ...props.Text) {
	d.styledText(text, props.Text{
		Size:  10,
		Style: consts.Bold,
	}, textProps...)
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
