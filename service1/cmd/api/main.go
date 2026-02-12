package main

import (
	"github.com/Uranury/exploreMicro/service1/internal/handlers"
	"github.com/Uranury/exploreMicro/service1/internal/storage"
	"log"
	"net/http"
)

func main() {
	store := storage.NewStore()
	userHandler := handlers.NewUser(store)
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				userHandler.GetUser(w, r)
			} else {
				userHandler.GetUsers(w, r)
			}
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		case http.MethodPatch:
			userHandler.UpdateUser(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
