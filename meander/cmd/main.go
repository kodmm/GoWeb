package main

import (
	"encoding/json"
	"net/http"

	"github.com/kodmm/GoWeb/meander"
)

func main() {
	// meander.APIKey = "TODO"
	http.HandleFunc("/journeys", func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	})
	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
