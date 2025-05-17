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
	Uuid        string `json:"uuid"`
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

		if record.Uuid != "0" {
			count, err = strconv.Atoi(record.Uuid)
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

func (c *FileStorage) readLine(hash string) (*record, error) {
	sugar := log.GetLogger()
	sugar.Infoln("readLine start, hash", hash)

	file, err := os.OpenFile(c.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
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

func (p *FileStorage) writeLine(record *record) error {
	sugar := log.GetLogger()
	data, err := json.Marshal(&record)
	if err != nil {
		sugar.Error("writeLine Marshal error", err)
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		sugar.Error("writeLine p.writer.Write error", err)
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		sugar.Error("writeLine p.writer.WriteByte error", err)
		return err
	}

	sugar.Infow("writeLine success", "record", record)

	return p.writer.Flush()
}

func (c *FileStorage) Close() error {
	sugar := log.GetLogger()
	sugar.Info("Close")
	return c.file.Close()
}

func (s *FileStorage) Get(hash string) (string, error) {
	record, err := s.readLine(hash)
	if err != nil {
		return "", err
	}

	return record.OriginalURL, nil
}

func (s *FileStorage) Add(u string) (string, error) {
	if u == "" {
		return "", errors.New("url is empty")
	}

	if _, err := url.ParseRequestURI(u); err != nil {
		return "", err
	}

	count := s.linesCount + 1

	h := hash.GetHash([]byte(u))
	r := record{
		Uuid:        strconv.Itoa(count),
		ShortURL:    h,
		OriginalURL: u,
	}

	err := s.writeLine(&r)
	if err != nil {
		return "", err
	}

	s.linesCount = count

	return h, nil
}
