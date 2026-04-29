package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/demianfiot/ticketproject/auth-service/internal/repository"
	"github.com/demianfiot/ticketproject/auth-service/internal/service"
	httptransport "github.com/demianfiot/ticketproject/auth-service/internal/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not loaded")
	}
	logrus.SetFormatter(new(logrus.JSONFormatter))
	dbConfig := DBConfigFromEnv()
	bd, err := NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer bd.Close()
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	tokenTTLHours := getEnvAsInt("TOKEN_TTL_HOURS", 12)

	repos := repository.NewRepository(bd)
	services := service.NewService(repos, jwtSecret, tokenTTLHours)
	handler := httptransport.NewHandler(services)

	router := gin.Default()
	httptransport.RegisterRoutes(router, handler)

	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("auth-service started on :%s", port)

	if err := router.Run(":" + port); err != nil {
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
