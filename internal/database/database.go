package database

import (
	"sync"
)

type Chirp struct {
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {

}

func (db *DB) CreateChrip(body string) (Chirp, error) {

}
