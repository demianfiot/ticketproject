package main

import (
	"log"
	"net"

	transport "github.com/demianfiot/ticketproject/ai-service/internal/transport/grpc"
	ai "github.com/demianfiot/ticketproject/ai-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	handler := transport.NewHandler()
	ai.RegisterAIServiceServer(grpcServer, handler)

	reflection.Register(grpcServer)

	log.Println("gRPC server started on :50051")
	log.Fatal(grpcServer.Serve(lis))
}
