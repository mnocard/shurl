package main

import (
	"flag"
	"os"
)

var config struct {
	flagRunAddr  string
	flagBaseAddr string
}

func parseFlags() {
	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); !ok {
		flag.StringVar(&config.flagRunAddr, "a", ":8080", "address and port to run server")
	} else {
		config.flagRunAddr = addr
	}

	if base, ok := os.LookupEnv("BASE_URL"); !ok {
		flag.StringVar(&config.flagBaseAddr, "b", "http://localhost:8080", "base address for short url")
	} else {
		config.flagBaseAddr = base
	}

	flag.Parse()
}
