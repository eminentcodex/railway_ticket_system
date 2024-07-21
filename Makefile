.PHONY: compile-proto
compile-proto:
	rm -rf protos/ticket
	mkdir -p protos/ticket
	protoc --go_out=./protos --go-grpc_out=./protos protos/ticket.proto