package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", ":32000")
	if err != nil {
		log.Fatalf("failed to create listener with error : %s", err.Error())
	}

	s := grpc.NewServer()

	// register all services with the grpc server
	RegisterTicketServer(s)

	log.Printf("starting grpc server at: %v\n", l.Addr())

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		if err = s.Serve(l); err != nil {
			log.Fatalf("failed to start server due to : %s ", err.Error())
		}
	}()

	<-done

	log.Println("shutting the server down")
	// gracefull shutdown
	s.GracefulStop()
	log.Println("server stopped")
}
