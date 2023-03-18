package router

import (
	"monolith/fruits"
	"net/http"
)

func ListAllFruits(r *http.Request, w http.ResponseWriter) {
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

func SetFruit(r *http.Request, w http.ResponseWriter) {
	// Sets a fruit
	// path = PUT /fruit
	id := activateSession(r)
	if id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("not logged in"))
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
