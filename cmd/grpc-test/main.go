package main

import (
	"context"
	grpc2 "github.com/onosproject/onos-test/test/grpc"
	"google.golang.org/grpc"
	"net"
)

func main() {
	service := &Service{}
	server := grpc.NewServer()
	grpc2.RegisterTestServiceServer(server, service)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	server.Serve(lis)
}

type Service struct {
}

func (s *Service) RequestReply(ctx context.Context, message *grpc2.Message) (*grpc2.Message, error) {
	return message, nil
}
