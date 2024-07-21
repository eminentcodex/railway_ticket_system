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
		request := &ticket.ReceiptRequest{}
		fmt.Print("Enter the ticket/receipt id: ")
		fmt.Scanf("%s", &request.TicketID)
		response, err := ca.c.GetReceipt(context.Background(), request)
		if err != nil {
			log.Printf("There was an error: %s\n", err.Error())
			break
		}
		log.Println("Ticket is :")
		prettyPrint(response)
	case 3:
		request := &ticket.SectionUserRequest{}
		fmt.Print("Enter the section: ")
		fmt.Scanf("%s", &request.Section)
		response, err := ca.c.GetUserBySection(context.Background(), request)
		if err != nil {
			log.Printf("There was an error: %s\n", err.Error())
			break
		}
		log.Println("Users :")
		prettyPrint(response)
	case 4:
		request := &ticket.RemoveUserRequest{}
		fmt.Print("Enter the email id: ")
		fmt.Scanf("%s", &request.Email)
		response, err := ca.c.RemoveUser(context.Background(), request)
		if err != nil {
			log.Printf("There was an error: %s\n", err.Error())
			break
		}
		log.Println("Status :")
		prettyPrint(response)
	case 5:
		request := &ticket.UpdateUserSeatRequest{}
		fmt.Print("Enter the ticket id: ")
		fmt.Scanf("%s", &request.TicketID)

		fmt.Print("Enter the new section: ")
		fmt.Scanf("%s", &request.Section)

		fmt.Print("Enter the new seat: ")
		fmt.Scanf("%d", &request.Seat)

		response, err := ca.c.UpdateUserSeat(context.Background(), request)
		if err != nil {
			log.Printf("There was an error: %s\n", err.Error())
			break
		}
		log.Println("Status :")
		prettyPrint(response)
	default:
		log.Fatalln("invalid option")
	}
}

func prettyPrint(s interface{}) {
	b, _ := json.MarshalIndent(s, "", "\t")
	fmt.Println(string(b))
}
