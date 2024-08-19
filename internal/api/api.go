package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	router := mux.NewRouter()

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
