package fakturago

import (
	"io"
	"time"

	"gopkg.in/yaml.v2"
)

type BillingInfo struct {
	Number  string    `json:"number"`
	Date    time.Time `json:"date"`
	Company Company   `json:"company"`
	BillTo  Company   `json:"bill_to"`
}

type Company struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	City      string `json:"city"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
	VatNumber string `json:"vat_number"`
}

func LoadBillingInfo(r io.Reader) (BillingInfo, error) {
	info := BillingInfo{}
	err := yaml.NewDecoder(r).Decode(&info)
	return info, err
}
