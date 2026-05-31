package main

import (
	"strings"
	"sync"

	"docker-check/docker"
	"docker-check/registry"

	"github.com/moby/moby/api/types/container"
	digest "github.com/opencontainers/go-digest"
)

type tableDataType struct {
	name      string
	version   string
	hasUpdate string
	link      string
	err       error
}

func main() {
	// Инициализация docker-api
	cli, err := docker.New()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// Инициализация клиента для работы с удалённым репозиторием
	reg, err := registry.New()
	if err != nil {
		panic(err)
	}

	// Получение списка контейнеров
	containers := cli.ContainerList().Items

	tableData := make(chan tableDataType, len(containers))
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 10)

	for _, ctr := range containers {
		wg.Add(1)
		go func(ctr container.Summary) {
			defer wg.Done()

			semaphore <- struct{}{}        // Занять слот
			defer func() { <-semaphore }() // Освободить слот

			// Получение информации об образе
			imageInspectResult, err := cli.ImageInspectResult(ctr.ImageID)
			if err != nil {
				return
			}

			// Если образ локальный, то пропуск
			if cli.IsLocalImage(imageInspectResult) {
				return
			}

			// Если нет дайджеста, то пропуск
			if len(imageInspectResult.RepoDigests) == 0 {
				return
			}

			// Получение дайджеста
			repoDigestParts := strings.Split(imageInspectResult.RepoDigests[0], "@")
			localDigest := digest.Digest(repoDigestParts[1])

			// Заполнение данных таблицы
			containerName := registry.GetContainerName(ctr.Names)
			image, _ := registry.ParseImage(ctr.Image)
			version := image.GetVersion(ctr.Labels)
			updateResult := reg.CheckForUpdate(ctr.Image, localDigest)

			tableData <- tableDataType{
				name:      containerName,
				version:   version,
				hasUpdate: updateResult.HasUpdate,
				link:      image.HubLink,
			}
		}(ctr)
	}

	// Завершение горутин
	go func() {
		wg.Wait()
		close(tableData)
	}()

	// Вывод данных в консоль
	tableRender(tableData)
}
