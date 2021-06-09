package main

import (
	"fmt"
	"strings"
)

// ReusedPassword is a password that is used by more than one user
type ReusedPassword struct {
	Hash  string
	Count int
	Users []string
}

// NewReusedPassword creates a new ReusedPassword
func NewReusedPassword(user string, hash string) *ReusedPassword {
	pw := ReusedPassword{}

	pw.Hash = hash
	pw.Count = 1
	pw.Users = make([]string, 1)
	pw.Users[0] = user

	return &pw
}

// Add adds a new user to the reused password
func (reusedpassword *ReusedPassword) Add(user string) {
	reusedpassword.Count++
	reusedpassword.Users = append(reusedpassword.Users, user)
}

func (reusedpassword *ReusedPassword) String() string {
	return fmt.Sprintf("%dx %s (%s)", reusedpassword.Count, reusedpassword.Hash, strings.Join(reusedpassword.Users, ","))
}
