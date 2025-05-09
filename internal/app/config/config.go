package config

import (
	"flag"
	"os"

	log "github.com/mnocard/shurl/internal/app/logger/zap"
)

var addresses *Addr

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

func parseFlags() {
	if addresses == nil {
		addresses = &Addr{}
	}

	sugar := log.GetLogger()
	sugar.Info("parseFlags()")
	flag.StringVar(&addresses.FlagRun, flagA, defRunAddr, "address and port to run server")
	flag.StringVar(&addresses.FlagBase, flagB, defBaseAddr, "base address for short url")
	flag.Parse()

	sugar.Infow("parseFlags()", "addresses", addresses)

	if addr, ok := os.LookupEnv(envServerAddress); ok && addr != "" {
		addresses.FlagRun = addr
	}

	if base, ok := os.LookupEnv(envBaseURL); ok && base != "" {
		addresses.FlagBase = base
	}
	sugar.Infow("parseFlags()", "addresses", addresses)
}

func GetAddresses() *Addr {
	sugar := log.GetLogger()
	sugar.Info("GetAddresses()")
	if addresses != nil {
		sugar.Info("GetAddresses() addresses != nil")
		return addresses
	}

	sugar.Info("GetAddresses() addresses == nil")
	sugar.Infow("GetAddresses() 1", "addresses", addresses)
	parseFlags()
	sugar.Infow("GetAddresses() 2", "addresses", addresses)
	return addresses
}
