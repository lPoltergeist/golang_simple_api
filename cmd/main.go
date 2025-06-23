package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/google/uuid"
)

type deleteRequest struct {
	ID string `json:"id"`
}

type ageName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var slice []ageName

func removeIndex(arr []ageName, index int) []ageName {
	return append(arr[:index], arr[index+1:]...)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func handleAddAge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method, only POST is permitted.", http.StatusMethodNotAllowed)
	}

	var person ageName

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, "Error to save this person's age", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newPersonAge := ageName{uuid.NewString(), person.Name, person.Age}

	slice = append(slice, newPersonAge)

	fmt.Printf("People Registered: %+v\n", slice)

	fmt.Fprintf(w, "Success on receiving JSON %+v\n", len(slice))
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	var idToDelete deleteRequest

	err := json.NewDecoder(r.Body).Decode(&idToDelete)
	if err != nil {
		http.Error(w, "Failed to delete!", http.StatusBadRequest)
	}
	defer r.Body.Close()

	fmt.Printf("%v", idToDelete)

	idx := slices.IndexFunc(slice, func(a ageName) bool {
		return a.ID == idToDelete.ID
	})

	if idx != -1 {
		slice = removeIndex(slice, idx)
		fmt.Printf("People removed: %+v\n", slice)
		fmt.Fprintf(w, "Success on delete! %v\n", slice)
	} else {
		http.Error(w, "Failed to delete!", http.StatusNotFound)
	}

}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/add-person", handleAddAge)
	http.HandleFunc("/remove-person", deletePerson)

	http.ListenAndServe(":8080", nil)
}
