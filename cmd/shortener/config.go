package main

import (
	"flag"
)

var config struct {
	flagRunAddr  string
	flagBaseAddr string
}

func parseFlags() {
	flag.StringVar(&config.flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.flagBaseAddr, "b", "http://localhost:8080", "base address for short url")
	flag.Parse()
}
