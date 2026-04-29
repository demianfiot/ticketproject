package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"

	"github.com/demianfiot/ticketproject/ticket-service/internal/metrics"
	"github.com/demianfiot/ticketproject/ticket-service/internal/repository"
	"github.com/demianfiot/ticketproject/ticket-service/internal/service/events"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	prometheus.MustRegister(metrics.KafkaMessages)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: viper.GetStringSlice("kafka.bootstrap_servers"),
		Topic:   viper.GetString("kafka.topic_tickets"),
		GroupID: viper.GetString("kafka.group_id"),
	})
	defer reader.Close()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("metrics server started on :2112")
		http.ListenAndServe(":2112", nil)
	}()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{
			viper.GetString("clickhouse.host") + ":" + viper.GetString("clickhouse.port"),
		},
		Auth: clickhouse.Auth{
			Database: viper.GetString("clickhouse.database"),
			Username: viper.GetString("clickhouse.user"),
			Password: viper.GetString("clickhouse.password"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("failed to ping clickhouse: %v", err)
	}

	analyticsRepo := repository.NewAnalyticsRepository(conn)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		<-sigCh

		log.Println("shutting down consumer...")
		cancel()
	}()

	log.Println("ticket events consumer started")
	//loop
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Println("reader stopped:", err)
			return
		}

		var event events.TicketCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("invalid message:", err)

			if commitErr := reader.CommitMessages(context.Background(), msg); commitErr != nil {
				log.Println("failed to commit invalid message:", commitErr)
			}

			continue
		}

		if err := analyticsRepo.InsertTicketCreatedEvent(ctx, event); err != nil {
			log.Println("failed to insert ticket event into clickhouse:", err)
			continue
		}
		metrics.KafkaMessages.Inc()

		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Println("commit failed:", err)
			continue
		}

		log.Printf(
			"ticket.created stored: event_id=%s ticket_id=%d user_id=%s status=%s created_at=%s",
			event.EventID,
			event.TicketID,
			event.UserID,
			event.Status,
			event.CreatedAt.Format(time.RFC3339),
		)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	return viper.ReadInConfig()
}
