package main

import (
	"context"
	"eminentcodex/railway_ticket_system/protos/ticket"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	ticket.UnimplementedRailwayServiceServer
	// global in-memory storage
	TotalSeatsPerSection int
	Tickets              map[string]map[int32]*ticket.ReceiptResponse //TODO: modify to allow concurrent access and race condition
	Users                map[string]*ticket.User                      //TODO: modify to allow concurrent access and race condition
}

func RegisterTicketServer(s *grpc.Server) {
	ticket.RegisterRailwayServiceServer(s, &Server{
		Tickets: map[string]map[int32]*ticket.ReceiptResponse{
			"A": make(map[int32]*ticket.ReceiptResponse, 0),
			"B": make(map[int32]*ticket.ReceiptResponse, 0),
		},
		TotalSeatsPerSection: 5,
		Users:                make(map[string]*ticket.User, 0),
	})
}

func (s *Server) PurchaseTicket(ctx context.Context, request *ticket.PurchaseTicketRequest) (*ticket.ReceiptResponse, error) {
	// check if any section is available
	var assignedSection string
	var seatNum int32
	if request.GetUser() == nil {
		return nil, errors.New("user cannot be blank")
	}
	if assignedSection, seatNum = CheckAvailableSection(s.Tickets, s.TotalSeatsPerSection); assignedSection == "" {
		return nil, errors.New("no seats available")
	}
	// add user
	AddUser(s.Users, request)
	// create a ticket using the available section
	ticketID := fmt.Sprintf("T%d", time.Now().Unix())
	newTicket := GenerateTicket(s.Tickets, assignedSection, seatNum, ticketID, request)

	return newTicket, nil
}
func (s *Server) GetReceipt(ctx context.Context, request *ticket.ReceiptRequest) (*ticket.ReceiptResponse, error) {
	ticketID := request.GetTicketID()
	if ticketID == "" {
		return nil, errors.New("please provide a ticket id")
	}
	receipt := SearchReceipt(s.Tickets, ticketID)
	if receipt == nil {
		return nil, errors.New("ticket doesn't exists")
	}

	return receipt, nil
}
func (s *Server) GetUserBySection(ctx context.Context, request *ticket.SectionUserRequest) (*ticket.SectionUserResonse, error) {
	users := GetUserBySection(s.Tickets, request.GetSection())
	if len(users) == 0 {
		return nil, errors.New("no users found")
	}
	return &ticket.SectionUserResonse{Users: users}, nil
}
func (s *Server) RemoveUser(ctx context.Context, request *ticket.RemoveUserRequest) (*ticket.RemoveUserResponse, error) {
	if request.GetEmail() == "" {
		return nil, errors.New("provide a valid email id")
	}
	removed := RemoveUser(s.Users, request.GetEmail())
	response := &ticket.RemoveUserResponse{}
	if removed {
		response.Message = "Removed successfully"
	} else {
		response.Message = "failed to remove"
	}
	return response, nil
}
func (s *Server) UpdateUserSeat(ctx context.Context, request *ticket.UpdateUserSeatRequest) (*ticket.UpdateUserSeatResponse, error) {
	if !CheckSetAvailability(s.Tickets, request.GetSection(), request.GetSeat()) {
		return nil, errors.New("seat is not available")
	}
	updated := UpdateUserSeat(s.Tickets, request.GetTicketID(), request.GetSection(), request.GetSeat())
	response := &ticket.UpdateUserSeatResponse{}
	if updated {
		response.Message = "updated successfully"
	} else {
		response.Message = "failed to update"
	}
	return response, nil
}
