package main

import (
	"log"
	"order/config/db"
	"order/proto/compiled"
	"order/service"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func createGRPCClient() (compiled.ProductServiceClient, func()) {
	conn, err := grpc.Dial("localhost:50001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := compiled.NewProductServiceClient(conn)
	return client, func() { conn.Close() }
}

func main() {
	app := fiber.New()

	client, closeConn := createGRPCClient()
	defer closeConn()

	db, closeDb := db.OpenConnection()
	defer closeDb()

	app.Get("/order/:code", func(c *fiber.Ctx) error {
		return service.GetOrder(c, db, &client)
	})
	app.Post("/order", func(c *fiber.Ctx) error {
		return service.AddOrder(c, db, &client)
	})
	app.Patch("/order/:code", func(c *fiber.Ctx) error {
		return service.UpdateOrder(c, db, &client)
	})
	app.Delete("/order/:code", func(c *fiber.Ctx) error {
		return service.DeleteOrder(c, db)
	})

	log.Fatal(app.Listen(":8080"))
}
