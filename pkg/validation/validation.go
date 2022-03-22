package validation

import (
	"encoding/json"

	"github.com/google/uuid"
)

func IsJSON[T string | []byte](str T) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
