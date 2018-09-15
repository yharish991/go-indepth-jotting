package main

import (
	"encoding/json"
	"fmt"
)

// User struct needs to be in capital because
// we will be passing this struct to the json package
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func main() {
	user := User{
		Username: "Ankur Anand",
		Password: "change it",
		Email:    "info@info.com",
	}
	dat, _ := json.Marshal(user)
	fmt.Println(string(dat))
	// {"username":"Ankur Anand","password":"change it","email":"info@info.com"}
	userUN := &User{}
	json.Unmarshal(dat, userUN)
	fmt.Println(*userUN)
}
