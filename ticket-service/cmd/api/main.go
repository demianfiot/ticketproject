package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	_ = godotenv.Load(".env")
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err)
	}

	aiAddr := viper.GetString("ai.address")
	log.Println("AI address:", aiAddr)

	aiClient, err := aisvc.NewGRPCClient(aiAddr)
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
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("ticket-service started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down ticket-service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("ticket-service stopped")
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
	db, err := connectPostgresWithRetry(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLmode))
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
func connectPostgresWithRetry(dsn string) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for attempt := 1; attempt <= 10; attempt++ {
		db, err = sqlx.Open("postgres", dsn)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			return db, nil
		}

		db.Close()
		log.Printf("postgres not ready, attempt %d/10: %v", attempt, err)
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
