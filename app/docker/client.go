package docker

import (
	"context"

	"github.com/moby/moby/client"
)

type Client struct {
	ctx       context.Context
	apiClient *client.Client
}

// Новый инстанс docker-api
func New() (*Client, error) {
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	return &Client{
		ctx:       ctx,
		apiClient: apiClient,
	}, err
}

// Close closes docker client
func (c *Client) Close() {
	if c.apiClient != nil {
		_ = c.apiClient.Close()
	}
}
