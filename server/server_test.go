package main

import (
	"context"
	"eminentcodex/railway_ticket_system/protos/ticket"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var (
	lis        *bufconn.Listener
	grpcServer *grpc.Server
)

func setup() {
	buf := 1024 * 1024
	lis = bufconn.Listen(buf)
	grpcServer = grpc.NewServer()
	// register the grpc service
	RegisterTicketServer(grpcServer)
	// start the mock server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("server failed with :", err)
		}
	}()
}

func teardown() {
	if lis != nil {
		lis.Close()
	}

	if grpcServer != nil {
		grpcServer.Stop()
	}
}

func bufconnDialer(context.Context, string) (net.Conn, error) { return lis.Dial() }

func createClient() *grpc.ClientConn {
	conn, err := grpc.NewClient("passthrough://", grpc.WithContextDialer(bufconnDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to create client", err.Error())
	}

	return conn
}

func TestMain(t *testing.M) {
	setup()
	t.Run()
	teardown()
}

func TestPurchaseTicket(t *testing.T) {
	conn := createClient()
	defer conn.Close()
	client := ticket.NewRailwayServiceClient(conn)
	table := []struct {
		name          string
		request       *ticket.PurchaseTicketRequest
		wantResp      bool
		expectedError bool
	}{
		{
			name: "success",
			request: &ticket.PurchaseTicketRequest{
				From:  "Source",
				To:    "Destination",
				Price: 12.49,
				User: &ticket.User{
					FirstName: "First Name",
					LastName:  "Last Name",
					Email:     "mock@example.com",
				},
			},
			wantResp:      true,
			expectedError: false,
		},
		{
			name: "failure",
			request: &ticket.PurchaseTicketRequest{
				From:  "Source",
				To:    "Destination",
				Price: 12.49,
			},
			wantResp:      false,
			expectedError: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.PurchaseTicket(context.Background(), tt.request)
			if tt.wantResp {
				assert.NotNil(t, resp)
			}
			if tt.expectedError {
				assert.NotNil(t, err)
			}
		})
	}
}
