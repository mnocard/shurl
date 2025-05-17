package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strconv"

	"github.com/mnocard/shurl/internal/app/hash"
	log "github.com/mnocard/shurl/internal/app/middleware/logger/zap"
)

type FileStorage struct {
	file       *os.File
	writer     *bufio.Writer
	filePath   string
	linesCount int
}

type record struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	sugar := log.GetLogger()
	reader := bufio.NewReader(file)
	count := 0
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}

		record := record{}

		err = json.Unmarshal(data, &record)
		if err != nil {
			return nil, err
		}

		if record.UUID != "0" {
			count, err = strconv.Atoi(record.UUID)
			if err != nil {
				return nil, err
			}
		}
	}

	sugar.Info("NewFileStorage success")

	return &FileStorage{
		file:       file,
		writer:     bufio.NewWriter(file),
		filePath:   filePath,
		linesCount: count,
	}, nil
}

func (f *FileStorage) readLine(hash string) (*record, error) {
	sugar := log.GetLogger()
	sugar.Infoln("readLine start, hash", hash)

	file, err := os.OpenFile(f.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			sugar.Info("readLine end of file")
			break
		}

		record := record{}

		err = json.Unmarshal(data, &record)
		if err != nil {
			sugar.Error("readLine Unmarshal error", err)
			return nil, err
		}

		if record.ShortURL == hash {
			sugar.Infow("readLine found", "record", record)
			return &record, nil
		}
	}

	sugar.Info("readLine not found")
	return nil, errors.New("not found")
}

func (f *FileStorage) writeLine(record *record) error {
	sugar := log.GetLogger()
	data, err := json.Marshal(&record)
	if err != nil {
		sugar.Error("writeLine Marshal error", err)
		return err
	}

	if _, err := f.writer.Write(data); err != nil {
		sugar.Error("writeLine p.writer.Write error", err)
		return err
	}

	if err := f.writer.WriteByte('\n'); err != nil {
		sugar.Error("writeLine p.writer.WriteByte error", err)
		return err
	}

	sugar.Infow("writeLine success", "record", record)

	return f.writer.Flush()
}

func (f *FileStorage) Close() error {
	sugar := log.GetLogger()
	sugar.Info("Close")
	return f.file.Close()
}

func (f *FileStorage) Get(hash string) (string, error) {
	record, err := f.readLine(hash)
	if err != nil {
		return "", err
	}

	return record.OriginalURL, nil
}

func (f *FileStorage) Add(u string) (string, error) {
	if u == "" {
		return "", errors.New("url is empty")
	}

	if _, err := url.ParseRequestURI(u); err != nil {
		return "", err
	}

	count := f.linesCount + 1

	h := hash.GetHash([]byte(u))
	r := record{
		UUID:        strconv.Itoa(count),
		ShortURL:    h,
		OriginalURL: u,
	}

	err := f.writeLine(&r)
	if err != nil {
		return "", err
	}

	f.linesCount = count

	return h, nil
}
