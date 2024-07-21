package main

import (
	"context"
	"eminentcodex/railway_ticket_system/protos/ticket"
	"encoding/json"
	"fmt"
	"log"
)

type ClientAction struct {
	c ticket.RailwayServiceClient
}

func NewClientAction(client ticket.RailwayServiceClient) *ClientAction {
	return &ClientAction{
		c: client,
	}
}

func (ca *ClientAction) Start() {
	var (
		done   bool
		action int
	)

	for !done {
		action = 0
		// ask for action
		fmt.Println("Select actions to perform")
		fmt.Println("\t1. Purchase Ticket")
		fmt.Println("\t2. Get Receipt")
		fmt.Println("\t3. Get User By Section")
		fmt.Println("\t4. Remove User")
		fmt.Println("\t5. Update User Seat")
		fmt.Println("\t6. Exit")
		fmt.Print("Action: ")
		fmt.Scanf("%d", &action)
		if action == 6 {
			break
		}
		ca.executeAction(action)
	}
	log.Println("Exiting...")
}

func (ca *ClientAction) executeAction(action int) {
	switch action {
	case 1:
		request := &ticket.PurchaseTicketRequest{
			User: &ticket.User{},
		}
		fmt.Print("Enter source location: ")
		fmt.Scanf("%s", &request.From)
		fmt.Print("Enter destination location: ")
		fmt.Scanf("%s", &request.To)
		fmt.Print("Enter Price: ")
		fmt.Scanf("%f", &request.Price)
		fmt.Print("Enter user first name: ")
		fmt.Scanf("%s", &request.User.FirstName)
		fmt.Print("Enter user last name: ")
		fmt.Scanf("%s", &request.User.LastName)
		fmt.Print("Enter user email name: ")
		fmt.Scanf("%s", &request.User.Email)
		response, err := ca.c.PurchaseTicket(context.Background(), request)
		if err != nil {
			log.Printf("There was an error: %s\n", err.Error())
			break
		}
		log.Println("Ticket is :")
		prettyPrint(response)
	case 2:
	case 3:
	case 4:
	case 5:
	default:
		log.Fatalln("invalid option")
	}
}

func prettyPrint(s interface{}) {
	b, _ := json.MarshalIndent(s, "", "\t")
	fmt.Println(string(b))
}
