package utils

import "strings"

// NormalizeAddress trims and converts an Ethereum address to lowercase
func NormalizeAddress(address string) string {
	return strings.ToLower(strings.TrimSpace(address))
}
