package config

import (
	"flag"
	"os"

	log "github.com/mnocard/shurl/internal/app/middleware/logger/zap"
)

var config *Config

type Config struct {
	FlagRun         string
	FlagBase        string
	FileStoragePath string
}

const (
	envServerAddress = "SERVER_ADDRESS"
	envBaseURL       = "BASE_URL"
	envFilePath      = "FILE_STORAGE_PATH"
	flagA            = "a"
	flagB            = "b"
	flagF            = "f"
	defRunAddr       = ":8080"
	defBaseAddr      = "http://localhost:8080"
	defFileStorage   = "storage.txt"
)

func parseFlags() {
	if config == nil {
		config = &Config{}
	}

	flag.StringVar(&config.FlagRun, flagA, defRunAddr, "address and port to run server")
	flag.StringVar(&config.FlagBase, flagB, defBaseAddr, "base address for short url")
	flag.StringVar(&config.FileStoragePath, flagF, defFileStorage, "file storage path")
	flag.Parse()

	if addr, ok := os.LookupEnv(envServerAddress); ok && addr != "" {
		config.FlagRun = addr
	}

	if base, ok := os.LookupEnv(envBaseURL); ok && base != "" {
		config.FlagBase = base
	}

	if filePath, ok := os.LookupEnv(envFilePath); ok && filePath != "" {
		config.FileStoragePath = filePath
	}

	sugar := log.GetLogger()
	sugar.Infow("parseFlags.", "config", config)
}

func GetConfig() *Config {
	sugar := log.GetLogger()
	if config != nil {
		sugar.Info("GetAddresses. addresses != nil")
		return config
	}

	sugar.Info("GetAddresses. addresses == nil")
	parseFlags()
	sugar.Infow("GetAddresses.", "addresses", config)
	return config
}
