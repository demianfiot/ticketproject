package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	go func() {
		log.Println("ai-service gRPC started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve grpc: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down ai-service...")
	grpcServer.GracefulStop()
	log.Println("ai-service stopped")
}
