package memory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		isError bool
	}{
		{
			name:    "add url success",
			url:     "http://localhost:8383",
			isError: false,
		},
		{
			name:    "add error invalid url",
			url:     "asd",
			isError: true,
		},
		{
			name:    "add error empty url",
			url:     "",
			isError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewMemoryStorage()

			hash, err := storage.Add(test.url)
			if test.isError {
				require.Error(t, err)
			} else {
				require.NotEmpty(t, hash)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		addedURL    string
		expectedURL string
		url         string
		isError     bool
	}{
		{
			name:        "get url success",
			addedURL:    "http://localhost:8383",
			expectedURL: "http://localhost:8383",
			isError:     false,
		},
		{
			name:        "get error url not found",
			addedURL:    "",
			expectedURL: "",
			isError:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewMemoryStorage()
			hash, _ := storage.Add(test.addedURL)
			url, err := storage.Get(hash)

			if test.isError {
				require.Error(t, err)
			} else {
				require.Equal(t, test.expectedURL, url)
			}
		})
	}
}
