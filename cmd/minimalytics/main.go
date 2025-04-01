package main

import (
	"log"
	"net/http"
	"os"

	"github.com/antontuzov/minimalytics/internal/handlers"
	"github.com/antontuzov/minimalytics/internal/middleware"
	"github.com/antontuzov/minimalytics/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	store, err := storage.NewSQLiteStorage("./analytics.db")
	if err != nil {
		log.Fatal("Failed to initialize storage: ", err)
	}
	defer store.Close()

	h := handlers.NewHandler(store)

	http.HandleFunc("/track", middleware.RateLimit(h.TrackHandler))
	http.HandleFunc("/dashboard", middleware.BasicAuth(h.DashboardHandler))
	http.HandleFunc("/api/", middleware.BasicAuth(h.APIHandler))

	port := getPort()
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}
