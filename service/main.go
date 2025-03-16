package main

import (
	"context"
	"log"
	"user-service/server"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	service := server.GetService(ctx)
	log.Println("INFO: Starting user service")
	service.StartServer(ctx, cancel)

}
