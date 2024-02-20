package utils

import (
	"github.com/google/uuid"
)

func GenerateUUID() string {
	uuidObj := uuid.New()
	return uuidObj.String()
}

