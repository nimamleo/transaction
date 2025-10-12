package tigerbeetle

import (
	"fmt"

	tb "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type Client struct {
	client tb.Client
}

func NewClient(cfg Config) (*Client, error) {
	addresses := []string{fmt.Sprintf(":%s", cfg.Port)}

	client, err := tb.NewClient(types.ToUint128(cfg.ClusterID), addresses, 32)
	if err != nil {
		return nil, fmt.Errorf("create tigerbeetle client: %w", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) GetClient() tb.Client {
	return c.client
}
