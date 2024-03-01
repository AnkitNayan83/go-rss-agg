package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload) // converts payload to binary
	if err != nil {
		log.Printf("Failed to marshal JSON response-> %v\n", payload)
		w.WriteHeader(500) //status code
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error", msg)
	}

	type errResponse struct {
		Error string `json:"error"` // by doing this key of the json object will be error -> error: {...}
	}

	respondWithJson(w, code, errResponse{
		Error: msg,
	})
}
