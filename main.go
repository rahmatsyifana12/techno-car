package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"techno-car/dto"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// dummy database
var cars []dto.Car

func getCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(cars) == 0 {
		http.Error(w, "No car is available", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(cars)
}

func createCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newCar dto.Car
	if err := json.NewDecoder(r.Body).Decode(&newCar); err != nil {
		http.Error(w, "Error createCar input", http.StatusBadRequest)
		return
	}

	newCar.ID = uuid.NewString()
	newCar.CalculatePrice()

	cars = append(cars, newCar)

	json.NewEncoder(w).Encode(cars)
}

func updateCarById(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")

    var updatedCar dto.Car
    if err := json.NewDecoder(r.Body).Decode(&updatedCar); err != nil {
        http.Error(w, "Error updateCarById input", http.StatusBadRequest)
        return
    }

    params := mux.Vars(r)
    inputCarId := params["id"]

    for i, item := range cars {
        if item.ID == inputCarId {
            // create copy of prev cars, with inputCar remove
            cars = append(cars[:i], cars[i+1:]...)
            // siapin mobil[0] -> mobil[index-1]
            // siapin juga mobil[index+1] -> mobil[lastIndex]

            updatedCar.ID = item.ID
            updatedCar.CalculatePrice()
            cars = append(cars, updatedCar)

            json.NewEncoder(w).Encode(updatedCar)
            return
        }
    }

    http.Error(w, "No car with the given ID was found!", http.StatusNotFound)
}

func deleteCarById(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")

    params := mux.Vars(r)
    inputCarId := params["id"]

    for i, item := range cars {
        if item.ID == inputCarId {
            cars = append(cars[:i], cars[i+1:]...)

            json.NewEncoder(w).Encode(item)
            return
        }
    }

    http.Error(w, "No car with the given ID was found!", http.StatusNotFound)
}

func main() {
	// init router
	router := mux.NewRouter()

	// route endpoints
	// get all cars
	router.HandleFunc("/api/v1/cars", getCars).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/cars", createCar).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/cars/{id}", updateCarById).Methods(http.MethodPut)
	router.HandleFunc("/api/v1/cars/{id}", deleteCarById).Methods(http.MethodDelete) 

	fmt.Println("Application running at localhost:8001")

	var port string
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "8001"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}