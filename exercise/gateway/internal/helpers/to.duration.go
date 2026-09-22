package helpers

import (
	"fmt"
	"time"
)

// anyToTime converts any type to time.Duration type
// if cannot convert result is 0
func ToDuration(v any) (time.Duration, error) {

	switch val := v.(type) {
	case time.Duration:
		return val, nil
	case float64: // JSON numbers come as float64
		return time.Duration(val * float64(time.Second)), nil
	case int:
		return time.Duration(val) * time.Second, nil
	case int64:
		return time.Duration(val) * time.Second, nil
	case string:
		return time.ParseDuration(val)
	default:
		return 0, fmt.Errorf("cannot convert %T to time.Duration", v)
	}
}
