package util

import (
	"strings"
)

// base64urlEncode performs Base64 URL encoding by replacing '+' with '-' and '/' with '_'
func Base64urlEncode(base64Str string) string {
	encoded := strings.ReplaceAll(base64Str, "+", "-")
	return strings.ReplaceAll(encoded, "/", "_")
}

// base64urlDecode performs Base64 URL decoding by replacing '-' with '+' and '_' with '/'
func Base64urlDecode(encoded string) string {
	encoded = strings.ReplaceAll(encoded, "-", "+")
	return strings.ReplaceAll(encoded, "_", "/")
}
