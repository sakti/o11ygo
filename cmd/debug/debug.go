package main

import (
	"log"

	"github.com/sakti/o11ygo/pkg/tracing"
)

func main() {
	log.Println("debug...")
	tracing.Run()
}
