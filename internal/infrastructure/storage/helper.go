package storage

import (
	"fmt"

	"github.com/google/uuid"
)

type ObjectName string

const (
	PlayerProfile ObjectName = "players/%s/profile.webp"
	TeamLogo      ObjectName = "teams/%s/logo.webp"
)

func BuildObjectName(name ObjectName, id uuid.UUID) string {
	return fmt.Sprintf(string(name), id)
}
