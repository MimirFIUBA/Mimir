package restApi

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
)

type response struct {
	Data int `json:"data"`
}

func getData(w http.ResponseWriter, r *http.Request) {
	var data = response{
		Data: rand.Intn(10),
	}
  
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func Start() {
	http.Handle("/data", http.HandlerFunc(getData))
  	log.Fatal(http.ListenAndServe(":8080", nil))
}