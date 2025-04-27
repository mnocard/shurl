package main

import (
	"flag"
	"log"
	"os"
)

var config struct {
	flagRunAddr  string
	flagBaseAddr string
}

func parseFlags() {
	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		config.flagRunAddr = addr
		log.Print("addr, ok")
	} else {
		flag.StringVar(&config.flagRunAddr, "a", ":8080", "address and port to run server")
		log.Print("addr, !ok")
	}

	if base, ok := os.LookupEnv("BASE_URL"); ok {
		config.flagBaseAddr = base
		log.Print("base, ok")
	} else {
		flag.StringVar(&config.flagBaseAddr, "b", "http://localhost:8080", "base address for short url")
		log.Print("base, !ok")
	}

	flag.Parse()
	log.Print("addr: " + config.flagRunAddr)
	log.Print("base: " + config.flagBaseAddr)
}
