package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFileStorage(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		isError  bool
	}{
		{
			name:     "new file storage success",
			filePath: "test.txt",
			isError:  false,
		},
		{
			name:     "new file storage error empty file path",
			filePath: "",
			isError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage, err := NewFileStorage(test.filePath)

			if test.isError {
				require.Error(t, err)
			} else {
				require.NotNil(t, storage)
				storage.Close()
				os.Remove(test.filePath)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		url      string
		isError  bool
	}{
		{
			name:     "add url success",
			filePath: "test.txt",
			url:      "http://localhost:8383",
			isError:  false,
		},
		{
			name:     "add error invalid url",
			filePath: "test.txt",
			url:      "asd",
			isError:  true,
		},
		{
			name:     "add error empty url",
			filePath: "test.txt",
			url:      "",
			isError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage, _ := NewFileStorage(test.filePath)
			defer func() {
				storage.Close()
				os.Remove(test.filePath)
			}()

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
		filePath    string
		addedURL    string
		expectedURL string
		url         string
		isError     bool
	}{
		{
			name:        "get url success",
			filePath:    "test.txt",
			addedURL:    "http://localhost:8383",
			expectedURL: "http://localhost:8383",
			isError:     false,
		},
		{
			name:        "get error url not found",
			filePath:    "test.txt",
			addedURL:    "",
			expectedURL: "",
			isError:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage, _ := NewFileStorage(test.filePath)
			defer func() {
				storage.Close()
				os.Remove(test.filePath)
			}()

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
