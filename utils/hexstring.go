package utils

import (
	"fmt"
	"strings"
)

func ToHexString(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	sb := strings.Builder{}
	for _, b := range data {
		sb.WriteString(fmt.Sprintf(" %02x", b))
	}
	return sb.String()[1:]
}
