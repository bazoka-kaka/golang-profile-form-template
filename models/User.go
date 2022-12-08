package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type User struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Alias     string `json:"alias"`
}

func (u *User) WriteJsonData() error {
	// read file
	file, err := os.OpenFile("users.json", os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// write file
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var users []User
	err = json.Unmarshal(content, &users)
	if err != nil {
		return err
	}

	users = append(users, *u)
	fmt.Println(users)

	jsonData, err := json.Marshal(users)
	if err != nil {
		return err
	}

	err = os.WriteFile("users.json", jsonData, 0667)
	if err != nil {
		return err
	}
	return nil
}
