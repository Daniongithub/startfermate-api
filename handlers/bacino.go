package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Daniongithub/startfermate-api/services"
)

func Bacino(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	selected := q.Get("selectedOption")

	data, err := services.GetBacino(selected)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// mapping leggero come nel Node
	type out struct {
		Nome     string `json:"nome"`
		Palina   string `json:"palina"`
		TargetID string `json:"targetID"`
	}

	result := make([]out, len(data))

	for i, v := range data {
		result[i] = out(v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
