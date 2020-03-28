package config

import (
	"os"
	"strconv"
)

const defaultPerPage = 50

type Config struct {
	DBDSN   string
	BindAdr string
	PerPage uint64
}

func LoadConfig() *Config {
	perPageRaw := os.Getenv("PER_PAGE")
	var perPage int
	if perPageRaw == "" {
		perPage = defaultPerPage
	}
	if perPage, _ = strconv.Atoi(perPageRaw); perPage == 0 {
		perPage = defaultPerPage
	}
	return &Config{
		DBDSN:   os.Getenv("DB_DSN"),
		BindAdr: os.Getenv("BIND_ADDR"),
		PerPage: uint64(perPage),
	}
}
