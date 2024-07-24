package api

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	router.HandleFunc("/data", getData).Methods("GET") //TODO: just to test, remove if not necessary

	router.HandleFunc("/sensors", getSensors).Methods("GET")
	router.HandleFunc("/sensors/{id}", getSensor).Methods("GET")
	router.HandleFunc("/sensors", createSensor).Methods("POST")
	router.HandleFunc("/sensors/{id}", updateSensor).Methods("PUT")
	router.HandleFunc("/sensors/{id}", deleteSensor).Methods("DELETE")

	router.HandleFunc("/nodes", getNodes).Methods("GET")
	router.HandleFunc("/nodes/{id}", getNode).Methods("GET")
	router.HandleFunc("/nodes", createNode).Methods("POST")
	router.HandleFunc("/nodes/{id}", updateNode).Methods("PUT")
	router.HandleFunc("/nodes/{id}", deleteNode).Methods("DELETE")

	router.HandleFunc("/groups", getGroups).Methods("GET")
	router.HandleFunc("/groups/{id}", getGroup).Methods("GET")
	router.HandleFunc("/groups", createGroup).Methods("POST")
	router.HandleFunc("/groups/{id}", updateGroup).Methods("PUT")
	router.HandleFunc("/groups/{id}", deleteGroup).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
