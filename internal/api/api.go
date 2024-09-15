package api

import (
	"log"
	"mimir/internal/api/controllers"
	"mimir/internal/api/routes"
	"net/http"

	"github.com/gorilla/websocket"
)

// TODO: add security check (only for production use)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var clients = make(map[*websocket.Conn]bool)

var broadcast chan string

func Start(broadcastChan chan string) {
	broadcast = broadcastChan
	// router := mux.NewRouter()

	// router.HandleFunc("/sensors", getSensors).Methods("GET")
	// router.HandleFunc("/sensors/{id}", getSensor).Methods("GET")
	// router.HandleFunc("/sensors", createSensor).Methods("POST")
	// router.HandleFunc("/sensors/{id}", updateSensor).Methods("PUT")
	// router.HandleFunc("/sensors/{id}", deleteSensor).Methods("DELETE")

	// router.HandleFunc("/nodes", getNodes).Methods("GET")
	// router.HandleFunc("/nodes/{id}", getNode).Methods("GET")
	// router.HandleFunc("/nodes", createNode).Methods("POST")
	// router.HandleFunc("/nodes/{id}", updateNode).Methods("PUT")
	// router.HandleFunc("/nodes/{id}", deleteNode).Methods("DELETE")

	// router.HandleFunc("/groups", getGroups).Methods("GET")
	// router.HandleFunc("/groups/{id}", getGroup).Methods("GET")
	// router.HandleFunc("/groups", createGroup).Methods("POST")
	// router.HandleFunc("/groups/{id}", updateGroup).Methods("PUT")
	// router.HandleFunc("/groups/{id}", deleteGroup).Methods("DELETE")

	// router.HandleFunc("/triggers", getTriggers).Methods("GET")
	// router.HandleFunc("/triggers/{id}", getTrigger).Methods("GET")
	// router.HandleFunc("/triggers", createTrigger).Methods("POST")
	// router.HandleFunc("/triggers/{id}", updateTrigger).Methods("PUT")
	// router.HandleFunc("/triggers/{id}", deleteTrigger).Methods("DELETE")

	// router.HandleFunc("/processors", getProcessors).Methods("GET")
	// router.HandleFunc("/processors/{id}", getProcessor).Methods("GET")
	// router.HandleFunc("/processors", createProcessor).Methods("POST")
	// router.HandleFunc("/processors/{id}", updateProcessor).Methods("PUT")
	// router.HandleFunc("/processors/{id}", deleteProcessor).Methods("DELETE")

	// router.HandleFunc("/ws", handleConnections)

	// go handleMessages()

	router := routes.CreateRouter()
	sensorRouter := router.PathPrefix("/sensors").Subrouter()
	sensorRouter.HandleFunc("/", controllers.GetSensors).Methods("GET")
	sensorRouter.HandleFunc("/", controllers.CreateSensor).Methods("POST")
	sensorRouter.HandleFunc("/{id}", controllers.GetSensorById).Methods("GET")
	sensorRouter.HandleFunc("/{id}", controllers.UpdateSensor).Methods("PUT")
	sensorRouter.HandleFunc("/{id}", controllers.DeleteSensor).Methods("DELETE")

	// router.HandleFunc("/sensors", controllers.GetSensors).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
	log.Fatal(http.ListenAndServe(":8080", sensorRouter))
}
