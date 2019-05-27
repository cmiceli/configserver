package configserver

import (
	"encoding/json"
	"errors"
	"log"
	"os"
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

type FSStorage struct {
	dir string
}

func (m *FSStorage) Get(identifier string) (Config, error) {
	log.Printf("FSStorage: Get(%s)", identifier)
	var x Config
	f, err := os.Open(m.dir + identifier)
	if err != nil {
		log.Printf("FSStorage: failed to open %s with %v", identifier, err)
		return Config{}, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&x)
	return x, err
}

func (m *FSStorage) Set(identifier string, cfg Config) error {
	log.Printf("FSStorage: Set(%s)", identifier)
	cfg.Date = time.Now()
	f, err := os.OpenFile(m.dir+identifier, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("FSStorage: failed to open %s with %v", identifier, err)
		return err
	}
	defer f.Close()
	b, err := json.Marshal(cfg)
	if err != nil {
		log.Printf("FSStorage: failed to serialise %s", identifier)
		return err
	}
	_, err = f.Write(b)
	return err
}

func (m *FSStorage) LastUpdate(identifier string) (time.Time, error) {
	f, err := os.Lstat(m.dir + identifier)
	if err != nil {
		return time.Time{}, err
	}
	return f.ModTime(), nil
}

func NewFSStorage(dir string) Storage {
	return &FSStorage{dir: dir}
}
