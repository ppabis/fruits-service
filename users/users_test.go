package users

import (
	"database/sql"
	"os"
	"testing"
)

func TestCreatePasswordTooShort(t *testing.T) {
	// Creates a new user
	defer os.Remove("monolith.db")
	err := CreateUser("foo", "bar")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestCreateUser(t *testing.T) {
	// Creates a new user
	defer os.Remove("monolith.db")
	err := CreateUser("foo", "barbarbar")
	if err != nil {
		t.Error(err)
	}

	db, err := sql.Open("sqlite3", "monolith.db")
	if err != nil {
		t.Error(err)
	}
	result := getUser(db, "foo", "barbarbar")

	if result == nil {
		t.Error("Expected user")
	}

}

func TestFailAuthenticateInexistentUser(t *testing.T) {
	// Authenticates user and returns a session cookie
	defer os.Remove("monolith.db")
	_, err := Authenticate("foo", "bar")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestFailAuthenticateExistingUser(t *testing.T) {
	// Creates, authenticates user and returns a session cookie
	defer os.Remove("monolith.db")
	err := CreateUser("foo", "barbarbar")
	if err != nil {
		t.Error(err)
	}
	_, err = Authenticate("foo", "bar")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestAuthenticateSuccess(t *testing.T) {
	// Creates, authenticates user and returns a session cookie
	defer os.Remove("monolith.db")
	CreateUser("foo", "barbarbar")

	cookie, err := Authenticate("foo", "barbarbar")
	if err != nil {
		t.Error(err)
	}

	_, err = GetByCookie(cookie)
	if err != nil {
		t.Error(err)
	}
}

func TestUsernameSingularity(t *testing.T) {
	// Tests if username is unique
	defer os.Remove("monolith.db")
	err := CreateUser("foo", "barbarbar")
	if err != nil {
		t.Error(err)
	}
	err = CreateUser("foo", "qwerrttyuio")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestUserSuper(t *testing.T) {
	defer os.Remove("monolith.db")
	err := CreateUser("foo", "barbarbar")
	if err != nil {
		t.Error(err)
	}
	err = UpdateUserSuperStatus(1, 1)
	if err != nil {
		t.Error(err)
	}
	super := IsUserSuper(1)
	if !super {
		t.Error("Expected user to be super")
	}

	err = UpdateUserSuperStatus(1, 0)
	if err != nil {
		t.Error(err)
	}
	super = IsUserSuper(1)
	if super {
		t.Error("Expected user to be not super")
	}
}

func TestLogout(t *testing.T) {
	defer os.Remove("monolith.db")
	err := CreateUser("foo", "barbarbar")
	if err != nil {
		t.Error(err)
	}
	cookie, err := Authenticate("foo", "barbarbar")
	if err != nil {
		t.Error(err)
	}

	_, err = GetByCookie(cookie)
	if err != nil {
		t.Error(err)
	}

	Logout(cookie)
	if err != nil {
		t.Error(err)
	}

	_, err = GetByCookie(cookie)
	if err == nil {
		t.Error("Expected error")
	}
}
