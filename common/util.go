package common

import "strings"

func FormatHex(raw string) string {
	raw = strings.ToLower(raw)
	raw = strings.TrimSpace(raw)
	return strings.TrimPrefix(raw, "0x")
}
