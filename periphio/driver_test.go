package periphio

import (
	"flag"
	"log"
	"os"
	"testing"

	"periph.io/x/host/v3"
)

var (
	drvI2C *PeriphioI2CDriver
)

func TestMain(m *testing.M) {
	flag.Parse()

	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}

	drvI2C, err = NewPeriphioI2CDriver("", 0x24)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
