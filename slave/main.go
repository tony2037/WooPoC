package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/tony2037/WooPoC/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedAddServiceServer
}

func main() {
	var protocol string = "tcp"
	var port string = ":3000"
	lis, err := net.Listen(protocol, port)
	if err != nil {
		log.Fatalf("Failed to listen:  %v", err)
	}

	slave_server := grpc.NewServer()
	pb.RegisterAddServiceServer(slave_server, &server{})
	reflection.Register(slave_server)
	if err := slave_server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Reference: https://github.com/shuza/kubernetes-go-grpc/blob/master/api/main.go
func (s *server) Compute(cxt context.Context, r *pb.AddRequest) (*pb.AddResponse, error) {
	result := &pb.AddResponse{}
	result.Result = r.A + r.B

	logMessage := fmt.Sprintf("A: %d   B: %d     sum: %d", r.A, r.B, result.Result)
	log.Println(logMessage)

	return result, nil
}
