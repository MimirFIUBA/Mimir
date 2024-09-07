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

	router.HandleFunc("/triggers", getTriggers).Methods("GET")
	router.HandleFunc("/triggers/{id}", getTrigger).Methods("GET")
	router.HandleFunc("/triggers", createTrigger).Methods("POST")
	router.HandleFunc("/triggers/{id}", updateTrigger).Methods("PUT")
	router.HandleFunc("/triggers/{id}", deleteTrigger).Methods("DELETE")

	router.HandleFunc("/processors", getProcessors).Methods("GET")
	router.HandleFunc("/processors/{id}", getProcessor).Methods("GET")
	router.HandleFunc("/processors", createProcessor).Methods("POST")
	router.HandleFunc("/processors/{id}", updateProcessor).Methods("PUT")
	router.HandleFunc("/processors/{id}", deleteProcessor).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
