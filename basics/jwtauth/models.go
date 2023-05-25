package jwtauth

import "github.com/golang-jwt/jwt/v5"

const SecretKey = "ewvfdjknl"

type CustomClaims struct {
	User                 string `json:"user"`
	Id                   uint   `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

type User struct {
	Id       uint
	Name     string
	Email    string
	Password []byte
}
