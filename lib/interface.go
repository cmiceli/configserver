package configserver

import (
	"time"
)

type Storage interface {
	Get(identifier string) (Config, error)
	Set(identifier string, cfg Config) error
	LastUpdate(identifier string) (time.Time, error)
}

type Config struct {
	Checksum string    `json:"checksum"`
	Config   string    `json:"config"`
	Date     time.Time `json:"date"`
}
