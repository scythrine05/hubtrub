package main

import (
	"log"
	"net/http"

	"blohub/internal/hub"

	"github.com/gorilla/mux"
)

func main() {

	// Initialize the hubs
	testHub := hub.NewMotionHub() // Movement hub
	//subHub := hub.NewSubscriptionHub()

	// Start the hubs
	go testHub.Run()
	//go subHub.Run()

	// Create a new Gorilla mux router
	router := mux.NewRouter()

	// WebSocket handlers using Gorilla mux
	router.HandleFunc("/ws/motion", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(testHub, w, r)
	})

	// router.HandleFunc("/ws/sub", func(w http.ResponseWriter, r *http.Request) {
	// 	hub.ServeWs(subHub, w, r)
	// })

	// Use the router as the HTTP handler
	http.Handle("/", router)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
