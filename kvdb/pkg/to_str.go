package pkg

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func ToStr(v any) (string, error) {
	switch v.(type) {
	case string:
		return v.(string), nil
	case int:
		return strconv.Itoa(v.(int)), nil
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v.(bool)), nil
	case []byte:
		return string(v.([]byte)), nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("failed to marshal value: %w", err)
		}

		return string(b), nil

	}
}
