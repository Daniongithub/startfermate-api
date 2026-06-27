package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Daniongithub/startfermate-api/services"
)

func Fermata(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	param := q.Get("param")
	param2 := q.Get("param2")
	palina := q.Get("palina")
	det := q.Get("det")

	data, err := services.GetFermata(param, param2, palina, det)
	if err != nil {
		http.Error(w, "errore nel recupero dati", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
