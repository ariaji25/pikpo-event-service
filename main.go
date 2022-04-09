package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"com.pikpo.events/app/database"
	eventService "com.pikpo.events/app/events"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// load the env
	if os.Getenv("mode") == "native" {
		godotenv.Load(".env")
	}
	// start db connections
	database.Database.Open()

	// run server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	eventService.RegisterEventServer(s, &eventService.Events{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
