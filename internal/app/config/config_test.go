package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name           string
		envRunAddr     string
		envBaseAddr    string
		flagRunAddr    string
		flagBaseAddr   string
		expectRunAddr  string
		expectBaseAddr string
	}{
		{
			name:           "testing env var",
			envRunAddr:     ":8282",
			envBaseAddr:    "http://localhost:8383",
			flagRunAddr:    "",
			flagBaseAddr:   "",
			expectRunAddr:  ":8282",
			expectBaseAddr: "http://localhost:8383",
		},
		{
			name:           "testing flags",
			envRunAddr:     "",
			envBaseAddr:    "",
			flagRunAddr:    ":8484",
			flagBaseAddr:   "http://localhost:8585",
			expectRunAddr:  ":8484",
			expectBaseAddr: "http://localhost:8585",
		},
		{
			name:           "testing env var and flags",
			envRunAddr:     ":8282",
			envBaseAddr:    "http://localhost:8383",
			flagRunAddr:    ":8484",
			flagBaseAddr:   "http://localhost:8585",
			expectRunAddr:  ":8282",
			expectBaseAddr: "http://localhost:8383",
		},
	}

	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(test.name, flag.ContinueOnError)

			var addr Addr

			t.Setenv(envServerAddress, test.envRunAddr)
			t.Setenv(envBaseURL, test.envBaseAddr)

			os.Args = []string{"cmd", "-" + flagA, test.flagRunAddr, "-" + flagB, test.flagBaseAddr}

			ParseFlags(&addr)
			assert.Equal(t, test.expectRunAddr, addr.FlagRun)
			assert.Equal(t, test.expectBaseAddr, addr.FlagBase)
		})
	}
}
