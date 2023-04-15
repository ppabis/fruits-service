package router

import (
	"monolith/token"
	"monolith/users"
	"net/http"
)

// Gets a token
// path = GET /token
func GetToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("only GET is allowed"))
		return
	}

	id := activateSession(r)

	if id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("not logged in"))
		return
	}

	username, err := users.GetUsername(nil, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	extras := map[string]interface{}{
		"super":    users.IsUserSuper(id),
		"username": username,
	}

	token, err := token.CreateToken(id, "fruits", extras)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
