package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUIDv7(noHyphens bool) string {
	// UUIDのv7を生成するロジックを実装
	uuid, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	if noHyphens {
		return strings.ReplaceAll(uuid.String(), "-", "")
	}
	return uuid.String()
}
