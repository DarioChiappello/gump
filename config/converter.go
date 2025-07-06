package config

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertToString(val interface{}) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(v), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func ConvertToInt(val interface{}, key string) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		// Clean spaces
		clean := strings.TrimSpace(v)

		// Direct try to convert to int
		if result, err := strconv.Atoi(clean); err == nil {
			return result, nil
		}

		// Try to convert to float and then to int
		if result, err := strconv.ParseFloat(clean, 64); err == nil {
			return int(result), nil
		}
	}

	return 0, &TypeError{Key: key, Expected: "int", Actual: fmt.Sprintf("%T", val)}
}

func ConvertToBool(val interface{}, key string) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil

	case string:
		// strings variants
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "true", "t", "yes", "y", "on", "1":
			return true, nil
		case "false", "f", "no", "n", "off", "0":
			return false, nil
		default:
			return false, &TypeError{Key: key, Expected: "bool", Actual: "string"}
		}

	case int:
		return v != 0, nil
	case int8:
		return v != 0, nil
	case int16:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case uint:
		return v != 0, nil
	case uint8:
		return v != 0, nil
	case uint16:
		return v != 0, nil
	case uint32:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	case float32:
		return v != 0.0, nil
	case float64:
		return v != 0.0, nil

	default:
		return false, &TypeError{Key: key, Expected: "bool", Actual: fmt.Sprintf("%T", val)}
	}
}
