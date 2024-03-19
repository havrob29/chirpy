package main

import (
	"errors"
	"fmt"
	"os"
)

// deletes database.json
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
