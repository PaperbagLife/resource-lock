package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

const RESOURCES_FILE = "resources.json"

// ResourceType enumeration
type ResourceType int

const (
	NIC    ResourceType = 0
	SERVER ResourceType = 1
)

// ResourceStatus structure
type ResourceStatus struct {
	Locked  bool   `json:"locked"`
	Name    string `json:"name,omitempty"`    // who locked
	EndTime int64  `json:"endTime,omitempty"` // Unix timestamp in milliseconds, or -1 for infinite length
}

// Nic structure
type Nic struct {
	Type   ResourceType   `json:"type"`
	Name   string         `json:"name"`
	Status ResourceStatus `json:"status"`
}

// Server structure
type Server struct {
	Type   ResourceType   `json:"type"`
	Name   string         `json:"name"`
	Status ResourceStatus `json:"status"`
	Nics   []Nic          `json:"nics,omitempty"`
}

// EncodeJSON encodes a Server struct into JSON
func EncodeJSON(server Server) ([]byte, error) {
	return json.Marshal(server)
}

// DecodeJSON decodes JSON data into a Server struct
func DecodeJSON(data []byte) (Server, error) {
	var server Server
	err := json.Unmarshal(data, &server)
	return server, err
}

func main() {
	mux := http.NewServeMux()

	// CORS middleware setup
	c := cors.Default() // Allow all origins
	handler := c.Handler(mux)

	// Registering handlers for different endpoints
	mux.HandleFunc("/status", handleStatus)
	mux.HandleFunc("/update", handleUpdate)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read data from file
	data, err := readFile(RESOURCES_FILE)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Respond with the file content as JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Write the request body to a file
	err = writeFile(RESOURCES_FILE, string(body))
	if err != nil {
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.Write([]byte("Update request received and data written to file"))
}

func writeFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func readFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
