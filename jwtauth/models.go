package jwtauth

type User struct {
	Id       uint
	Name     string
	Email    string `gorm:unique`
	Password []byte
}
