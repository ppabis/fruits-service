package router

import (
	"monolith/fruits"
	"net/http"
)

func ListAllFruits(w http.ResponseWriter, r *http.Request) {
	// Lists all fruits
	// path = GET /
	fruits, err := fruits.GetFruits()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
	for k, v := range fruits {
		w.Write([]byte(k + ":\t\t\t" + v + "\n"))
	}

}

func SetFruit(w http.ResponseWriter, r *http.Request) {
	// Sets a fruit
	// path = PUT /fruit
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("only PUT is allowed"))
		return
	}

	id := activateSession(r)
	if id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("not logged in"))
		return
	}

	if r.ParseForm() != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid form"))
		return
	}

	fruit := r.Form.Get("fruit")

	err := fruits.UpdateFruit(id, fruit)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Location", "/")
	w.Write([]byte("fruit set"))

}
