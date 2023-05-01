package e2e

import (
	"monolith/users"
)

var Password = "waxe4_QDq23s-2q2qa12"

// Creates three users with usernames as their passwords.
// 1 alice - super: true
// 2 bob - super: false
// 3 charlie - super: true
func CreateUsers() error {
	err := users.CreateUser("alice", Password)
	if err != nil {
		return err
	}
	_, err = users.Authenticate("alice", Password)
	if err != nil {
		return err
	}
	err = users.CreateUser("bob", Password)
	if err != nil {
		return err
	}
	err = users.UpdateUserSuperStatus(1, 1)
	if err != nil {
		return err
	}
	err = users.CreateUser("charlie", Password)
	if err != nil {
		return err
	}
	err = users.UpdateUserSuperStatus(3, 1)
	if err != nil {
		return err
	}
	return nil
}
