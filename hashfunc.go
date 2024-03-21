package main

import (
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	passByteArray := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passByteArray, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func comparePassword(toCompare string, savedPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(savedPass), []byte(toCompare))
}
