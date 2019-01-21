package configserver

import (
	"errors"
	"log"
	"time"
)

var NoSuchFile = errors.New("No such file")

type MemStorage struct {
	s map[string]Config
	l map[string]time.Time
}

func NewMemStorage() Storage {
	return &MemStorage{s: make(map[string]Config), l: make(map[string]time.Time)}
}

func (m *MemStorage) Get(identifier string) (Config, error) {
	log.Printf("Storage: Get(%s)", identifier)
	if v, ok := m.s[identifier]; ok {
		return v, nil
	}
	return Config{}, NoSuchFile
}

func (m *MemStorage) Set(identifier string, cfg Config) error {
	log.Printf("Storage: Set(%s)", identifier)
	cfg.Date = time.Now()
	m.s[identifier] = cfg
	m.l[identifier] = cfg.Date
	return nil
}

func (m *MemStorage) LastUpdate(identifier string) (time.Time, error) {
	log.Printf("Storage: LastUpdate(%s)", identifier)
	if v, ok := m.l[identifier]; ok {
		return v, nil
	}
	return time.Now(), NoSuchFile
}
