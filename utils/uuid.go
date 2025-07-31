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

// 複数のUUIDを生成するための関数
func GenerateMultipleUUIDs(count int, noHyphens bool) []string {
	if count <= 0 {
		return []string{}
	}
	uuidList := make([]string, count)
	for i := 0; i < count; i++ {
		uuidList[i] = GenerateUUIDv7(noHyphens)
	}
	return uuidList
}
