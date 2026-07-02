package handlers

import (
	"fmt"
	"net/http"
)

func Version(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, version)
	}
}
