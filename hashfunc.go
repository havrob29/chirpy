package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	passByteArray := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passByteArray, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hash))
	return string(hash)
}
