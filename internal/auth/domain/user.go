package domain

import "fmt"

type User struct {
	Name           string
	LastName       string
	Email          string
	Username       string
	HashedPassword string
}

func (u *User) GetFullName() string {
	return fmt.Sprintf("%s %s", u.Name, u.LastName)
}
