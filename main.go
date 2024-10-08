package main

import (
	"fmt"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the ERP system!")
}

func main() {
	http.HandleFunc("/", homeHandler)
	//http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/dashboard", dashboardHandler)

	http.HandleFunc("/start", handleControl)
	http.HandleFunc("/stop", handleControl)

	fmt.Println("Starting the web server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
