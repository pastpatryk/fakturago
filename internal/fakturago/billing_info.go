package fakturago

import (
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v2"
)

// BillingInfo contains all data needed to issue an invoice
type BillingInfo struct {
	Number   string        `yaml:"number"`
	Date     time.Time     `yaml:"date"`
	Company  Company       `yaml:"company"`
	BillTo   Company       `yaml:"bill_to"`
	Currency Currency      `yaml:"currency"`
	NoTax    bool          `yaml:"no_tax"`
	Items    []BillingItem `yaml:"items"`
}

// Company contains info about single company
type Company struct {
	Name      string `yaml:"name"`
	Address   string `yaml:"address"`
	City      string `yaml:"city"`
	ZipCode   string `yaml:"zip_code"`
	Country   string `yaml:"country"`
	VatNumber string `yaml:"vat_number"`
}

// BillingItem is a single item on invoice
type BillingItem struct {
	Name     string            `yaml:"name"`
	NameLang map[string]string `yaml:"name_lang"`
	Amount   float64           `yaml:"amount"`
	Tax      Tax               `yaml:"tax"`
}

// Tax represents tax value for single item
type Tax float64

// String returns percentage string representation of tax value eg. 23%
func (t Tax) String() string {
	return fmt.Sprintf("%d%%", int(t*100))
}

// Value returns float64 representation of tax value
func (t Tax) Value() float64 {
	return float64(t)
}

// Apply calculates amount after applying tax
func (t Tax) Apply(amount float64) float64 {
	return amount * (1 + t.Value())
}

// Currency holds string representation of invoice currency
type Currency string

// Format returns string representation of price amount in given currency
func (c Currency) Format(amount float64) string {
	return fmt.Sprintf("%.2f %s", amount, c)
}

// LocalizedName returns BillingItem name in given language
// or default Name if translation is missing
func (b BillingItem) LocalizedName(lang string) string {
	name := b.Name
	if b.NameLang != nil && b.NameLang[lang] != "" {
		name = b.NameLang[lang]
	}
	return name
}

// FinalAmount returns BillingItem price after applying modificators like tax
func (b BillingItem) FinalAmount() float64 {
	return b.Tax.Apply(b.Amount)
}

// LoadBillingInfo reads billing data from Reader and parses it using YAML
func LoadBillingInfo(r io.Reader) (BillingInfo, error) {
	info := BillingInfo{}
	err := yaml.NewDecoder(r).Decode(&info)
	return info, err
}
