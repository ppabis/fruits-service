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

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	printIndexPage(fruits, activateSession(r), w)

}

func SetFruit(w http.ResponseWriter, r *http.Request) {
	// Sets a fruit
	// path = PUT /fruit, POST /fruit
	if r.Method != "PUT" && r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("only PUT or POST is allowed"))
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

	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusFound)
	w.Write([]byte("fruit set"))

}
