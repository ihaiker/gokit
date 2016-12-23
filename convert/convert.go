package convertKit

import (
	"errors"
	"strconv"
)

// args: value, precision(only for float)
func String(args ...interface{}) (string, error) {
	value := args[0]
	var precision int = 12 // default

	switch value.(type) {
	case string:
		v, _ := value.(string)
		return v, nil
	case int:
		v, _ := value.(int)
		return strconv.Itoa(v), nil
	case int32:
		v, _ := value.(int32)
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		v, _ := value.(int64)
		return strconv.FormatInt(v, 10), nil
	case float32:
		v, _ := value.(float32)
		if len(args) >= 2 {
			precision = args[1].(int)
		}
		return strconv.FormatFloat(float64(v), 'f', precision, 64), nil
	case float64:
		v, _ := value.(float64)
		if len(args) >= 2 {
			precision = args[1].(int)
		}
		return strconv.FormatFloat(v, 'f', precision, 64), nil
	default:
		return "", errors.New("unknown type")
	}
}

func SafeString(args... interface{}) string {
	s,_ := String(args...)
	return s
}
func Int(value interface{}) (int, error) {
	switch value.(type) {
	case string:
		v, _ := value.(string)
		return strconv.Atoi(v)
	case int:
		v, _ := value.(int)
		return v, nil
	case int32:
		v, _ := value.(int32)
		return int(v), nil
	case int64:
		v, _ := value.(int64)
		return int(v), nil
	case float32:
		v, _ := value.(float32)
		return int(v), nil
	case float64:
		v, _ := value.(float64)
		return int(v), nil
	default:
		return int(0), errors.New("unknown type")
	}
}

func SafeInt(value interface{}) int {
	s,_ := Int(value)
	return s
}

func Int32(value interface{}) (int32, error) {
	switch value.(type) {
	case string:
		v, _ := value.(string)
		result, err := strconv.ParseInt(v, 10, 32)
		return int32(result), err
	case int:
		v, _ := value.(int)
		return int32(v), nil
	case int32:
		v, _ := value.(int32)
		return int32(v), nil
	case int64:
		v, _ := value.(int64)
		return int32(v), nil
	case float32:
		v, _ := value.(float32)
		return int32(v), nil
	case float64:
		v, _ := value.(float64)
		return int32(v), nil
	default:
		return int32(0), errors.New("unknown type")
	}
}

func SafeInt32(value interface{}) int32 {
	s,_:=Int32(value)
	return s
}

func Int64(value interface{}) (int64, error) {
	switch value.(type) {
	case string:
		v, _ := value.(string)
		return strconv.ParseInt(v, 10, 32)
	case int:
		v, _ := value.(int)
		return int64(v), nil
	case int32:
		v, _ := value.(int32)
		return int64(v), nil
	case int64:
		v, _ := value.(int64)
		return v, nil
	case float32:
		v, _ := value.(float32)
		return int64(v), nil
	case float64:
		v, _ := value.(float64)
		return int64(v), nil
	default:
		return int64(0), errors.New("unknown type")
	}
}

func SafeInt64(value interface{}) int64 {
	s , _ := Int64(value)
	return s
}

func Float32(value interface{}) (float32, error) {
	switch value.(type) {
	case string:
		v, _ := value.(string)
		result, err := strconv.ParseFloat(v, 32)
		return float32(result), err
	case int:
		v, _ := value.(int)
		return float32(v), nil
	case int32:
		v, _ := value.(int32)
		return float32(v), nil
	case int64:
		v, _ := value.(int64)
		return float32(v), nil
	case float32:
		v, _ := value.(float32)
		return v, nil
	case float64:
		v, _ := value.(float64)
		return float32(v), nil
	default:
		return float32(0), errors.New("unknown type")
	}
}

func SafeFloat32(value interface{}) float32 {
	s ,_ := Float32(value)
	return s
}

func Float64(value interface{}) (float64, error) {
	switch value.(type) {
	case string:
		v, _ := value.(string)
		return strconv.ParseFloat(v, 64)
	case int:
		v, _ := value.(int)
		return float64(v), nil
	case int32:
		v, _ := value.(int32)
		return float64(v), nil
	case int64:
		v, _ := value.(int64)
		return float64(v), nil
	case float32:
		v, _ := value.(float32)
		return float64(v), nil
	case float64:
		v, _ := value.(float64)
		return v, nil
	default:
		return float64(0), errors.New("unknown type")
	}
}

func SafeFloat64(value interface{}) float64 {
	s ,_ := Float64(value)
	return s
}

func Bool(value interface{}) (bool, error) {
	r, ok := value.(bool)
	if ok {
		return r, nil
	}
	if s, ok := value.(string); ok {
		return s == "true", nil
	} else if i, ok := value.(int); ok {
		return i > 0, nil
	} else if i, ok := value.(int32); ok {
		return i > 0, nil
	} else if i, ok := value.(int64); ok {
		return i > 0, nil
	} else if i, ok := value.(uint32); ok {
		return i > 0, nil
	} else if i, ok := value.(uint64); ok {
		return i > 0, nil
	}
	return r, errors.New("assert type fails")
}

func SafeBool(value interface{}) bool {
	s,_ := Bool(value)
	return s
}