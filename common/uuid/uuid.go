package uuid

import (
	"github.com/google/uuid"
	"strings"
)

func Uuid() string {
	return strings.ReplaceAll(RawUuid(), "-", "")
}

func RawUuid() string {
	return uuid.New().String()
}
