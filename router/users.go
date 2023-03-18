package router

import (
	"monolith/users"
	"net/http"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	// Logs in a user
	// Path = POST /login
	id := activateSession(r)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("only POST is allowed"))
		return
	}

	if id != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("already logged in"))
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	cookie, err := users.Authenticate(username, password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: cookie,
	})
	w.Header().Add("Location", "/")

}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Logs out a user
	// Path = GET /logout
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("only GET is allowed"))
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	users.Logout(cookie.Value)

	w.WriteHeader(http.StatusOK)
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	})
	w.Header().Add("Location", "/")
}

func activateSession(r *http.Request) int {
	// Activates a session
	// If there's a valid cookie within headers, the user ID != 0
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0
	}

	id, err := users.GetByCookie(cookie.Value)
	if err != nil {
		return 0
	}

	return id
}
