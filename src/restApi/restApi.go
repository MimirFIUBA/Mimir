package restApi

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"github.com/gorilla/mux"
	mimir "mimir/src/mimir"
)

type response struct {
	Data int `json:"data"`
}

type sensorsResponse struct {
	Sensors []mimir.Sensor `json:"sensors"`
}

func getData(w http.ResponseWriter, r *http.Request) {
	var data = response{
		Data: rand.Intn(10),
	}
  
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func createSensor(w http.ResponseWriter, r *http.Request) {

	var sensor mimir.Sensor
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sensor = mimir.CreateSensor(sensor)
	
	json.NewEncoder(w).Encode(sensor)
}

func getSensors(w http.ResponseWriter, r *http.Request) {
	var sensors = sensorsResponse {
		Sensors: mimir.GetSensors(),
	}
  
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensors)
}

func Start() {
	router := mux.NewRouter()
	router.HandleFunc("/data", getData).Methods("GET")
	router.HandleFunc("/sensor", createSensor).Methods("POST")
	router.HandleFunc("/sensor", getSensors).Methods("GET")

  	log.Fatal(http.ListenAndServe(":8080", router))
}