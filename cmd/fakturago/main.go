package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pastDexter/fakturago/internal/fakturago"

	log "github.com/sirupsen/logrus"
)

func main() {
	langPtr := flag.String("lang", "en", "Language used for generating invoice")
	outPtr := flag.String("out", "./invoice.pdf", "Output path")
	flag.Parse()

	log.SetLevel(log.DebugLevel)

	log.WithFields(log.Fields{
		"args":    flag.Arg(1),
		"outPtr":  *outPtr,
		"langPtr": *langPtr,
	}).Info("parsing CLI args")

	if len(flag.Args()) < 1 {
		fmt.Printf("Missing invoice defintion. Check '%s --help'\n", os.Args[0])
		os.Exit(1)
	}

	invoicePath := flag.Arg(0)

	file, err := os.Open(invoicePath)
	if err != nil {
		log.WithFields(log.Fields{"file": invoicePath}).Fatal("Unable to open file")
	}
	info, err := fakturago.LoadBillingInfo(file)
	if err != nil {
		log.Fatal("Unable to load billing info: ", err)
	}

	fakturago.Generate(info, *outPtr, *langPtr)
	fmt.Printf("Success! Invoice saved to: %s\n", *outPtr)
}
