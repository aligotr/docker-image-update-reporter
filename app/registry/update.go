package registry

import (
	"fmt"

	digest "github.com/opencontainers/go-digest"
	"go.podman.io/image/v5/docker"
)

type UpdateStatus struct {
	RemoteDigest digest.Digest
	LocalDigest  digest.Digest
	HasUpdate    string
}

// Сверка локального и удалённого дайджестов
func (c *Client) CheckForUpdate(imageWithTag string, localDigest digest.Digest) *UpdateStatus {
	remoteDigest, err := c.getRemoteDigest(imageWithTag)

	hasUpdate := "false"
	if err != nil {
		hasUpdate = "error"
	} else if localDigest != remoteDigest {
		hasUpdate = "true"
	}

	return &UpdateStatus{
		RemoteDigest: remoteDigest,
		LocalDigest:  localDigest,
		HasUpdate:    hasUpdate,
	}
}

// Получение дайджеста из удалённого репозитория
func (c *Client) getRemoteDigest(imageWithTag string) (digest.Digest, error) {
	ctx, cancel := c.timeoutContext()
	defer cancel()

	ref, err := docker.ParseReference("//" + imageWithTag)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга имени: %w", err)
	}

	remoteDigest, err := docker.GetDigest(ctx, c.sysCtx, ref)
	if err != nil {
		return "", fmt.Errorf("не удалось получить digest из репозитория: %w", err)
	}

	return remoteDigest, nil
}
