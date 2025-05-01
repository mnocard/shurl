package config

import (
	"flag"
	"log"
	"os"
)

type Addr struct {
	FlagRun  string
	FlagBase string
}

func ParseFlags(config *Addr) {
	flag.StringVar(&config.FlagRun, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.FlagBase, "b", "http://localhost:8080", "base address for short url")

	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok && addr != "" {
		config.FlagRun = addr
		log.Print("addr, ok")
	}

	if base, ok := os.LookupEnv("BASE_URL"); ok && base != "" {
		config.FlagBase = base
		log.Print("base, ok")
	}

	flag.Parse()
	log.Print("addr: " + config.FlagRun)
	log.Print("base: " + config.FlagBase)
}
