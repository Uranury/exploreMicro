package main

import (
	grpcsrv "github.com/Uranury/exploreMicro/service1/internal/grpc"
	"github.com/Uranury/exploreMicro/service1/internal/handlers"
	"github.com/Uranury/exploreMicro/service1/internal/storage"
	"github.com/Uranury/exploreMicro/service1/proto/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	store := storage.NewStore()

	userHandler := handlers.NewUser(store)
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, grpcsrv.NewUserService(store))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

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
