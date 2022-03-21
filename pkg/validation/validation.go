package validation

import "encoding/json"

func IsJSON[T string | []byte](str T) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
