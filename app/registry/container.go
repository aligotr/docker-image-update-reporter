package registry

import (
	"strings"
)

func GetContainerName(names []string) string {
	if len(names) == 0 {
		return "-"
	}
	// Docker обычно хранит имя как "/my-container"
	return strings.TrimPrefix(names[0], "/")
}
