package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	aisvc "github.com/demianfiot/ticketproject/ticket-service/internal/client/ai"
	"github.com/demianfiot/ticketproject/ticket-service/internal/metrics"
	"github.com/demianfiot/ticketproject/ticket-service/internal/repository"
	"github.com/demianfiot/ticketproject/ticket-service/internal/service"
	"github.com/demianfiot/ticketproject/ticket-service/internal/service/events"
	httptransport "github.com/demianfiot/ticketproject/ticket-service/internal/transport/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not loaded")
	}
	logrus.SetFormatter(new(logrus.JSONFormatter))

	aiClient, err := aisvc.NewGRPCClient("localhost:50051")
	if err != nil {
		log.Fatal(err)
	}
	defer aiClient.Close()

	dbConfig := DBConfigFromEnv()
	bd, err := NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer bd.Close()
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err)
	}

	log.Println("kafka brokers:", viper.GetStringSlice("kafka.bootstrap_servers"))
	log.Println("kafka topic:", viper.GetString("kafka.topic_tickets"))
	producer := events.NewKafkaProducer(
		viper.GetStringSlice("kafka.bootstrap_servers"),
		viper.GetString("kafka.topic_tickets"),
	)
	defer producer.Close()

	repository := repository.NewRepository(bd)
	service := service.NewService(aiClient, repository, producer)
	handler := httptransport.NewHandler(service)

	router := gin.Default()
	httptransport.RegisterRoutes(router, handler)
	prometheus.MustRegister(
		metrics.HTTPRequests,
		metrics.HTTPDuration,
	)
	log.Println("HTTP server started on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLmode  string
}

func NewPostgresDB(cfg DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLmode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func DBConfigFromEnv() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLmode:  os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}
