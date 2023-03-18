package router

import (
	"monolith/users"
	"net/http"
)

func LoginUser(r *http.Request, w http.ResponseWriter) {
	// Logs in a user
	// Path = POST /login
	id := activateSession(r)

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

	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: cookie,
	})
}

func LogoutUser(r *http.Request, w http.ResponseWriter) {
	// Logs out a user
	// Path = GET /logout
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
