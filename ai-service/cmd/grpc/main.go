package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/demianfiot/ticketproject/ai-service/internal/config"
	"github.com/demianfiot/ticketproject/ai-service/internal/llm"
	"github.com/demianfiot/ticketproject/ai-service/internal/parser"
	"github.com/demianfiot/ticketproject/ai-service/internal/prompt"
	"github.com/demianfiot/ticketproject/ai-service/internal/service"
	grpctransport "github.com/demianfiot/ticketproject/ai-service/internal/transport/grpc"
	aipb "github.com/demianfiot/ticketproject/ai-service/proto"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	promptBuilder := prompt.NewBuilder()
	analysisParser := parser.NewParser()

	llmClient := llm.NewOpenAIClient(
		cfg.OpenAIAPIKey,
		cfg.OpenAI.Model,
		cfg.OpenAI.BaseURL,
		time.Duration(cfg.OpenAI.TimeoutSeconds)*time.Second,
		promptBuilder,
		analysisParser,
	)

	analysisService := service.NewAnalysisService(llmClient)
	handler := grpctransport.NewHandler(analysisService)

	port := cfg.Server.GRPCPort
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	aipb.RegisterAIServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Println("ai-service gRPC server started on :" + port)
	log.Fatal(grpcServer.Serve(lis))
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	return viper.ReadInConfig()
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
