package main

import (
	"eminentcodex/railway_ticket_system/protos/ticket"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := os.Args[1]
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(address, opts)
	if err != nil {
		log.Fatalln("failed to connect to GRPC server with error %s", err.Error())
	}
	defer func() {
		if conn != nil {
			err = conn.Close()
			if err != nil {
				log.Println("failed to close connection")
			}
		}
	}()

	client := ticket.NewRailwayServiceClient(conn)
	ca := NewClientAction(client)
	ca.Start()
}
