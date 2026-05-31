package registry

// Исходник: https://github.com/crazy-max/diun/blob/master/pkg/registry/image.go

import (
	"fmt"
	"path/filepath"
	"strings"

	"go.podman.io/image/v5/docker/reference"
)

// Информация об образе
type Image struct {
	Name    string
	Domain  string
	Path    string
	Tag     string
	HubLink string

	named reference.Named
}

// ParseImage возвращает общую информацию о контейнере, на основе введённых данных
func ParseImage(name string) (Image, error) {
	// Извлечение названия и тэга
	named, err := reference.ParseNormalizedNamed(name)
	if err != nil {
		return Image{}, fmt.Errorf("%w parsing image %s failed", err, name)
	}
	named = reference.TagNameOnly(named)

	i := Image{
		named:  named,
		Name:   named.Name(),
		Domain: reference.Domain(named),
		Path:   reference.Path(named),
	}

	// Hub-link
	i.HubLink, err = i.hubLink()
	if err != nil {
		return Image{}, fmt.Errorf("%w resolving hub link for image %s failed", err, name)
	}

	// Тэг
	if tagged, ok := named.(reference.Tagged); ok {
		i.Tag = tagged.Tag()
	}

	return i, nil
}

func (i Image) hubLink() (string, error) {
	switch i.Domain {
	case "docker.io":
		prefix := "r"
		path := i.Path
		if strings.HasPrefix(i.Path, "library/") {
			prefix = "_"
			path = strings.Replace(i.Path, "library/", "", 1)
		}
		return fmt.Sprintf("https://hub.docker.com/%s/%s", prefix, path), nil
	case "docker.bintray.io", "jfrog-docker-reg2.bintray.io":
		return fmt.Sprintf("https://bintray.com/jfrog/reg2/%s", strings.ReplaceAll(i.Path, "/", "%3A")), nil
	case "docker.pkg.github.com":
		return fmt.Sprintf("https://github.com/%s/packages", filepath.ToSlash(filepath.Dir(i.Path))), nil
	case "gcr.io":
		return fmt.Sprintf("https://%s/%s", i.Domain, i.Path), nil
	case "ghcr.io":
		ref := strings.Split(i.Path, "/")
		ghUser, ghPackage := ref[0], ref[1]
		return fmt.Sprintf("https://github.com/users/%s/packages/container/package/%s", ghUser, ghPackage), nil
	case "quay.io":
		return fmt.Sprintf("https://quay.io/repository/%s", i.Path), nil
	case "registry.access.redhat.com":
		return fmt.Sprintf("https://access.redhat.com/containers/#/registry.access.redhat.com/%s", i.Path), nil
	case "registry.gitlab.com":
		return fmt.Sprintf("https://gitlab.com/%s/container_registry", i.Path), nil
	default:
		return fmt.Sprintf("https://%s", i.Name), nil
	}
}

func (i Image) GetVersion(labels map[string]string) string {
	version := labels["org.opencontainers.image.version"]
	if version == "" {
		version = labels["version"]
	}
	if version == "" {
		version = labels["dev.containers.variant"]
	}

	if version == "" {
		version = i.Tag
	}

	return version
}
