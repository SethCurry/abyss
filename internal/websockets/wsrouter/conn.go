package wsrouter

import (
	"github.com/google/uuid"
)

func newID() (string, error) {
	gotUUID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return gotUUID.String(), nil
}
