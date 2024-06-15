package utils

import (
	"fmt"
	"time"
)

func TimestampID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
