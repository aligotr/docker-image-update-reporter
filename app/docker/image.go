package docker

import (
	"github.com/moby/moby/client"
)

// Docker-api: получение информации об образе
func (c *Client) ImageInspectResult(imageID string) (client.ImageInspectResult, error) {
	return c.apiClient.ImageInspect(c.ctx, imageID)
}

func (c *Client) IsLocalImage(image client.ImageInspectResult) bool {
	return len(image.RepoDigests) == 0
}
