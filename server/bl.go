package main

import (
	"eminentcodex/railway_ticket_system/protos/ticket"
)

func CheckAvailableSection(totalSeats map[string]map[int32]*ticket.ReceiptResponse, totalSeatsPerSection int) (string, int32) {
	for section, seats := range totalSeats {
		if len(seats) < totalSeatsPerSection {
			// check available seats in the section
			var i int32
			for i = 1; i <= int32(totalSeatsPerSection); i++ {
				if _, ok := seats[i]; !ok {
					return section, i
				}
			}
		}
	}

	return "", 0
}

func AddUser(users map[string]*ticket.User, request *ticket.PurchaseTicketRequest) {
	if _, ok := users[request.User.Email]; !ok {
		users[request.User.Email] = request.User
	}
}

func GenerateTicket(totalSeats map[string]map[int32]*ticket.ReceiptResponse, assignedSection string, seatNum int32, ticketID string, request *ticket.PurchaseTicketRequest) *ticket.ReceiptResponse {
	newTicket := &ticket.ReceiptResponse{
		TicketID: ticketID,
		From:     request.From,
		To:       request.To,
		User:     request.User,
		Price:    request.Price,
		Section:  assignedSection,
		Seat:     seatNum,
	}
	totalSeats[assignedSection][seatNum] = newTicket

	return newTicket
}

func SearchReceipt(tickets map[string]map[int32]*ticket.ReceiptResponse, ticketID string) *ticket.ReceiptResponse {
	for _, sectionTickets := range tickets {
		for _, t := range sectionTickets {
			if ticketID == t.TicketID {
				return t
			}
		}
	}

	return nil
}

func GetUserBySection(tickets map[string]map[int32]*ticket.ReceiptResponse, section string) []*ticket.User {
	var users []*ticket.User
	sectionTickets := tickets[section]
	for _, t := range sectionTickets {
		users = append(users, t.User)
	}

	return users
}

func RemoveUser(users map[string]*ticket.User, email string) bool {
	if _, ok := users[email]; ok {
		delete(users, email)
		return true
	}

	return false
}

func UpdateUserSeat(tickets map[string]map[int32]*ticket.ReceiptResponse, ticketID string, section string, seat int32) bool {
	// remove the ticket first
	found := false
	var (
		oldSection string
		oldSeat    int32
		oldTicket  *ticket.ReceiptResponse
	)

	for s1, sectionTickets := range tickets {
		for s2, t := range sectionTickets {
			oldSeat = s2
			if ticketID == t.TicketID {
				oldTicket = t
				oldSection = s1
				found = true
				break
			}
		}
	}
	if !found {
		return false
	}

	delete(tickets[oldSection], oldSeat)
	ticketRequest := &ticket.PurchaseTicketRequest{
		From:  oldTicket.From,
		To:    oldTicket.To,
		User:  oldTicket.User,
		Price: oldTicket.Price,
	}

	_ = GenerateTicket(tickets, section, seat, oldTicket.TicketID, ticketRequest)
	return true
}

func CheckSetAvailability(tickets map[string]map[int32]*ticket.ReceiptResponse, section string, seat int32) bool {
	if _, ok := tickets[section][seat]; ok {
		return false // seat exists can't book
	}

	return true
}
