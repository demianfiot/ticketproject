package ai

import (
	"context"
	"time"

	aipb "github.com/demianfiot/ticketproject/ai-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client aipb.AIServiceClient
}

type AnalyzeInput struct {
	TicketID    string
	Title       string
	Description string
	UserID      string
}

type AnalyzeResult struct {
	Summary        string
	Category       string
	Priority       string
	SuggestedReply string
}

func NewGRPCClient(addr string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := aipb.NewAIServiceClient(conn)

	return &GRPCClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

func (c *GRPCClient) AnalyzeTicket(ctx context.Context, in AnalyzeInput) (AnalyzeResult, error) {
	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {
		callCtx, cancel := context.WithTimeout(ctx, 2*time.Second)

		resp, err := c.client.AnalyzeTicket(callCtx, &aipb.AnalyzeTicketRequest{
			TicketId:    in.TicketID,
			Title:       in.Title,
			Description: in.Description,
			UserId:      in.UserID,
		})

		cancel()

		if err == nil {
			return AnalyzeResult{
				Summary:        resp.Summary,
				Category:       resp.Category,
				Priority:       resp.Priority,
				SuggestedReply: resp.SuggestedReply,
			}, nil
		}

		lastErr = err
		time.Sleep(200 * time.Millisecond)
	}

	return AnalyzeResult{}, lastErr
}
