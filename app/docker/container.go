package docker

import (
	"github.com/moby/moby/client"
)

// Docker-api: получение списка контейнеров
func (c *Client) ContainerList() client.ContainerListResult {
	containers, err := c.apiClient.ContainerList(c.ctx, client.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}
	return containers
}
