package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Item struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type Sim struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Slogan string `json:"slogan"`
}

var items []Item
var (
    sims = make(map[string]Sim)
    mutex = &sync.RWMutex{}
)

func main() {
    router := mux.NewRouter()

    // Define routes
    router.HandleFunc("/items", getItems).Methods("GET")
    router.HandleFunc("/items", createItem).Methods("POST")
    router.HandleFunc("/items/{id}", getItem).Methods("GET")
    router.HandleFunc("/items/{id}", updateItem).Methods("PUT")
    router.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

    // Instantiate a simulator, operate it, and see the results
    router.HandleFunc("/sims", createSim).Methods("POST")  // create an instance using given configuration (or defaults)
    router.HandleFunc("/sims/{id}", getSimStatus).Methods("GET")  // return configuration and operational status
    router.HandleFunc("/sims/{id}", updateSimConfig).Methods("PUT")  // modify configuration
    router.HandleFunc("/sims/{id}/runs", progressSim).Methods("POST")  // run given (or default) iterations
    router.HandleFunc("/sims/{id}/runs", getSimData).Methods("GET")  // return run stats ?? by run ID? given search parameters?

    // Start server
    log.Println("Server is starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func createSim(w http.ResponseWriter, r *http.Request) {
    mutex.RLock()
    defer mutex.RUnlock()

    var simBud Sim
    err := json.NewDecoder(r.Body).Decode(&simBud)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
        return
    }

    simBud.ID = uuid.New().String()  // could check for duplicates, especially if shorter unique key is used

    // if _, exists := sims[simBud.ID]; exists {
    //     w.WriteHeader(http.StatusConflict)
    //     json.NewEncoder(w).Encode(map[string]string{"error": "Item with this ID already exists"})
    //     return
    // }

    sims[simBud.ID] = simBud
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(simBud)
}

func getSimStatus(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    if sim, found := sims[params["id"]]; found {
        json.NewEncoder(w).Encode(sim)
    } else {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
    }
}

func updateSimConfig(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    log.Println("Update configuration of " + vars["id"])
    w.WriteHeader(http.StatusNotImplemented)
}
func progressSim(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    log.Println("Run " + vars["id"])
    w.WriteHeader(http.StatusNotImplemented)
}
func getSimData(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    log.Println("Get run data of " + vars["id"])
    w.WriteHeader(http.StatusNotImplemented)
}

func getItems(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range items {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    json.NewEncoder(w).Encode(&Item{})
}

func createItem(w http.ResponseWriter, r *http.Request) {
    var item Item
    _ = json.NewDecoder(r.Body).Decode(&item)
    items = append(items, item)
    json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, item := range items {
        if item.ID == params["id"] {
            var updatedItem Item
            _ = json.NewDecoder(r.Body).Decode(&updatedItem)
            updatedItem.ID = params["id"]
            items[index] = updatedItem
            json.NewEncoder(w).Encode(updatedItem)
            return
        }
    }
    json.NewEncoder(w).Encode(&Item{})
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    items = slices.DeleteFunc(items, func(item Item) bool {
        return item.ID == params["id"]
    })
}