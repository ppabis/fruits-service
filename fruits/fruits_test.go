package fruits

import (
	"monolith/users"
	"os"
	"testing"
)

func TestCreateUserAndSetFruit(t *testing.T) {
	// Creates a new user and sets their fruit
	defer os.Remove("monolith.db")
	err := users.CreateUser("foo", "barbarbar")
	if err != nil {
		t.Error(err)
	}

	err = UpdateFruit(1, "apple")
	if err != nil {
		t.Error(err)
	}

	fruits, err := GetFruits()
	if err != nil {
		t.Error(err)
	}

	if fruits["foo"] != "apple" {
		t.Error("Expected apple for user foo")
	}

}

func TestCreateUserAndSetSpecialFruitFail(t *testing.T) {
	// Creates a new user and sets their fruit
	defer os.Remove("monolith.db")
	err := users.CreateUser("foo", "barbarbar")
	if err != nil {
		t.Error(err)
	}

	err = UpdateFruit(1, "pineapple")
	if err == nil {
		t.Error("Expected error")
	}

	fruits, err := GetFruits()
	if err != nil {
		t.Error(err)
	}

	if _, ok := fruits["foo"]; ok {
		t.Error("Expected no fruit for user foo")
	}

}

func TestCreateSuperUserAndSetSpecialFruit(t *testing.T) {
	// Creates a new user and sets their fruit
	defer os.Remove("monolith.db")
	err := users.CreateUser("foo", "barbarbar")
	if err != nil {
		t.Error(err)
	}

	err = users.UpdateUserSuperStatus(1, 1)
	if err != nil {
		t.Error(err)
	}

	err = UpdateFruit(1, "pineapple")
	if err != nil {
		t.Error(err)
	}

	fruits, err := GetFruits()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Fruits: %v", fruits)

	if fruits["foo"] != "pineapple" {
		t.Error("Expected pineapple for user foo")
	}

}
