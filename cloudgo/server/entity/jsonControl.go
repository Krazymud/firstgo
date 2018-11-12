package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	users    []user
	jsonFile *os.File
	err      error
)

type user struct {
	Username string
	Email    string
	Phone    string
	Number   string
}

func init() {
	jsonFile, err = os.Open("server/entity/store.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	bt, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(bt, &users)
}

func Register(userInfo map[string][]string) error {
	var usr user
	for k, v := range userInfo {
		switch k {
		case "username":
			usr.Username = v[0]
		case "email":
			usr.Email = v[0]
		case "number":
			usr.Number = v[0]
		case "phone":
			usr.Phone = v[0]
		}
	}
	err := CheckDup(usr)
	if err != nil {
		return err
	}
	users = append(users, usr)
	bt, _ := json.Marshal(users)
	err = ioutil.WriteFile("server/entity/store.json", bt, 0644)
	if err != nil {
		panic(err)
	}
	return nil
}

func GetUser(usrname string) user {
	for _, u := range users {
		if usrname == u.Username {
			return u
		}
	}
	return user{}
}

func CheckDup(usr user) error {
	for _, u := range users {
		if u.Username == usr.Username {
			return errors.New("用户名重复")
		}
		if u.Number == usr.Number {
			return errors.New("学号重复")
		}
		if u.Phone == usr.Phone {
			return errors.New("电话号码重复")
		}
		if u.Email == usr.Email {
			return errors.New("邮箱重复")
		}
	}
	return nil
}
