package main

import (
	"log"
	"net"
	"product/config/db"
	"product/proto/compiled"
	srv "product/service"

	"google.golang.org/grpc"
)

func main() {
	service := grpc.NewServer()

	db, closeDb := db.OpenConnection()
	defer closeDb()

	compiled.RegisterProductServiceServer(service, srv.ProductService{
		DB: db,
	})

	l, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatalf("could not listen to %s: %v", ":50001", err)
	}
	log.Fatal(service.Serve(l))
}
