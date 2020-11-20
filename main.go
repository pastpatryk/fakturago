package main

import (
	"github.com/pastDexter/fakturago/internal/fakturago"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	fakturago.GenerateInvoice("./invoice.pdf")
}
