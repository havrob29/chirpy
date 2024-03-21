package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type User struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	ID       int    `json:"id"`
}

type DBStructure struct {
	Chirps  map[int]Chirp        `json:"chirps"`
	Users   map[int]User         `json:"users"`
	Revoked map[string]time.Time `json:"revoked"`
}

type Chirp struct {
	Body   string `json:"body"`
	ID     int    `json:"id"`
	Author int    `json:"author_id"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type returnError struct {
	Error string `json:"error"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

// Revokes token by saving it to database along with the current time in UTC
func (db *DB) CreateRevoked(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	dbStructure.Revoked[token] = time.Now().UTC()

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetRevoked() ([]string, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	tokens := make([]string, 0, len(dbStructure.Revoked))
	for token := range dbStructure.Revoked {
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, userID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:     id,
		Body:   body,
		Author: userID,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	//check if email is already registered to a user
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return User{}, errors.New("email already registered")
		}
	}

	id := len(dbStructure.Users) + 1
	//hashPasswordToSave
	hashedPassword := hashPassword(password)

	user := User{
		ID:       id,
		Email:    email,
		Password: hashedPassword,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// updates user's email and password in database
func (db *DB) UpdateUser(id int, newEmail, newPassword string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	//hashPasswordToSave
	newHashedPassword := hashPassword(newPassword)

	user := User{
		ID:       id,
		Email:    newEmail,
		Password: newHashedPassword,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

// GetUsers return all users in the database
func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}
	return users, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.CreateDB()
	}
	return err
}

func (db *DB) CreateDB() error {
	dbStructure := DBStructure{
		Chirps:  map[int]Chirp{},
		Users:   map[int]User{},
		Revoked: map[string]time.Time{},
	}
	return db.writeDB(dbStructure)
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

// deletes database.json file
func delDB() error {
	_, err := os.ReadFile("database.json")
	if err != nil {
		return errors.New("no database to delete")
	} else {
		err = os.Remove("./database.json")
		if err != nil {
			return err
		} else {
			fmt.Println("delete successful...")
		}
	}
	return nil
}
