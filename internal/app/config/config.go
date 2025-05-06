package config

import (
	"flag"
	"os"
)

type Addr struct {
	FlagRun  string
	FlagBase string
}

const (
	envServerAddress = "SERVER_ADDRESS"
	envBaseURL       = "BASE_URL"
	flagA            = "a"
	flagB            = "b"
	defRunAddr       = ":8080"
	defBaseAddr      = "http://localhost:8080"
)

func ParseFlags(config *Addr) {
	flag.StringVar(&config.FlagRun, flagA, defRunAddr, "address and port to run server")
	flag.StringVar(&config.FlagBase, flagB, defBaseAddr, "base address for short url")
	flag.Parse()

	if addr, ok := os.LookupEnv(envServerAddress); ok && addr != "" {
		config.FlagRun = addr
	}

	if base, ok := os.LookupEnv(envBaseURL); ok && base != "" {
		config.FlagBase = base
	}
}
