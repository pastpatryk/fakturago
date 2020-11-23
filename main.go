package main

import (
	"os"

	"github.com/pastDexter/fakturago/internal/fakturago"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	file, err := os.Open("./invoice.yaml")
	if err != nil {
		log.Fatal("Unable to open file")
	}
	info, err := fakturago.LoadBillingInfo(file)
	if err != nil {
		log.Fatal("Unable to load billing info: ", err)
	}

	fakturago.Generate(info, "./invoice.pdf")
}
