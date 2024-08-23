package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

type Item struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

var items []Item

func main() {
    router := mux.NewRouter()

    // Define routes
    router.HandleFunc("/items", getItems).Methods("GET")
    router.HandleFunc("/items", createItem).Methods("POST")
    router.HandleFunc("/items/{id}", getItem).Methods("GET")
    router.HandleFunc("/items/{id}", updateItem).Methods("PUT")
    router.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

    // Start server
    log.Println("Server is starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
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