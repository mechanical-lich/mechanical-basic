package functions

import "errors"

func EnsureFloat(input interface{}) (float64, error) {
	var out float64
	switch v := input.(type) {
	case float64:
		out = v
	case int:
		out = float64(v)
	case int64:
		out = float64(v)
	case string:
		return 0, errors.New("cannot convert string to float")
	default:
		return 0, errors.New("invalid argument type for pow")
	}
	return out, nil
}

func EnsureInt(input interface{}) (int, error) {
	var out int
	switch v := input.(type) {
	case float64:
		out = int(v)
	case int:
		out = v
	case int64:
		out = int(v)
	case string:
		return 0, errors.New("cannot convert string to float")
	default:
		return 0, errors.New("invalid argument type for pow")
	}
	return out, nil
}
