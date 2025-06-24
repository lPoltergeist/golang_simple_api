package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreatePersonWithAge(age int, t *testing.T) []byte {
	personWithAge := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "Ramone",
		Age:  age,
	}

	jsonBody, err := json.Marshal(personWithAge)
	if err != nil {
		t.Error(err)
	}

	return jsonBody
}

func TestHelloHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(helloHandler))
	res, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, but got: %d", res.StatusCode)
	}

	expected := "Hello World"
	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("expected %s, but got %s", expected, string(b))
	}
}

func TestAddPersonAge(t *testing.T) {
	request := CreatePersonWithAge(28, t)

	req := httptest.NewRequest(http.MethodPost, "/add-person", bytes.NewBuffer(request))
	req.Header.Set("Content-Type", "application/json")

	//NewRecorder cria um mock de http.ResponseWriter, simulando o que o client receberia.
	rr := httptest.NewRecorder()

	handleAddAge(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Return wrong status Code. It shoud be: %v, but returned: %v", http.StatusOK, status)
	}
}

func TestNegativeAge(t *testing.T) {
	request := CreatePersonWithAge(0, t)

	req := httptest.NewRequest(http.MethodPost, "/add-person", bytes.NewBuffer(request))
	req.Header.Set("Content-Type", "application/json")

	//NewRecorder cria um mock de http.ResponseWriter, simulando o que o client receberia.
	rr := httptest.NewRecorder()

	handleAddAge(rr, req)

	expected := "Age not valid! Age must not be 0!"
	if response := rr.Body.String(); response == expected {
		t.Errorf("Expected: %v \n Received: %v", expected, response)
	}

	if status := rr.Code; status == http.StatusOK {
		t.Errorf("Return wrong status Code. It shoud be: %v, but returned: %v", http.StatusBadRequest, status)
	}
}
