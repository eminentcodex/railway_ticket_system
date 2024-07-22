.PHONY: compile-proto
compile-proto:
	rm -rf protos/ticket
	mkdir -p protos/ticket
	protoc --go_out=./protos --go-grpc_out=./protos protos/ticket.proto


.PHONY: run-server
run-server:
	./bin/server

.PHONY: run-client
run-client:
	./bin/client localhost:32000

.PHONY: build-sever
build-server:
	rm -f bin/server
	go build -o bin/server server/* 

.PHONY: build-client
build-client:
	rm -f bin/client
	go build -o bin/client client/*

.PHONY: test-server
test-server:
	go test ./server/...