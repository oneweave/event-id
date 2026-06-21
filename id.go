// Copyright 2025 Nadrama Pty Ltd
// SPDX-License-Identifier: Apache-2.0

package eventid

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// isValidPrefix checks if a prefix is 1 to 5 lowercase letters [a-z].
func isValidPrefix(prefix string) bool {
	l := len(prefix)
	if l < 1 || l > 5 {
		return false
	}
	for i := 0; i < l; i++ {
		c := prefix[i]
		if c < 'a' || c > 'z' {
			return false
		}
	}
	return true
}

// normalizeID normalizes an ID by converting to lowercase, removing spaces/hyphens,
// and identifying the prefix and base32 payload indices within the normalized buffer.
func normalizeID(id string, normalized *[128]byte) (prefixLen int, normLen int, err error) {
	dstIdx := 0

	for i := 0; i < len(id); i++ {
		c := id[i]
		if c == ' ' || c == '-' || c == '\n' || c == '\r' || c == '\t' {
			continue
		}

		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}

		if dstIdx >= len(normalized) {
			return 0, 0, fmt.Errorf("invalid ID length")
		}
		normalized[dstIdx] = c
		dstIdx++
	}

	underscoreIdx := -1
	for i := 0; i < dstIdx; i++ {
		if normalized[i] == '_' {
			underscoreIdx = i
			break
		}
	}

	if underscoreIdx == -1 {
		return 0, 0, fmt.Errorf("invalid format: missing '_' separator")
	}

	prefixLen = underscoreIdx
	payloadLen := dstIdx - (underscoreIdx + 1)

	if prefixLen < 1 || prefixLen > 5 {
		return 0, 0, fmt.Errorf("invalid prefix length: %d", prefixLen)
	}
	for i := 0; i < prefixLen; i++ {
		c := normalized[i]
		if c < 'a' || c > 'z' {
			return 0, 0, fmt.Errorf("invalid prefix character %q", c)
		}
	}

	if payloadLen != 26 {
		return 0, 0, fmt.Errorf("invalid payload length: %d (expected 26)", payloadLen)
	}

	return prefixLen, dstIdx, nil
}

// Encode encodes a UUID into a prefixed, crockford base32-encoded string
func Encode(uuidStr string, prefix string) (string, error) {
	if !isValidPrefix(prefix) {
		return "", fmt.Errorf("invalid prefix %s", prefix)
	}

	uuidStr = strings.TrimSpace(uuidStr)
	if strings.Contains(uuidStr, " ") {
		uuidStr = strings.ReplaceAll(uuidStr, " ", "")
	}

	u, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID format %s: %w", uuidStr, err)
	}

	var buf [32]byte
	prefixLen := len(prefix)
	copy(buf[0:prefixLen], prefix)
	buf[prefixLen] = '_'
	encodeBase32CrockfordToBuf(buf[prefixLen+1:prefixLen+1+26], u[:])
	return string(buf[:prefixLen+1+26]), nil
}

// Decode decodes a prefixed, crockford base32-encoded string
// into a UUID string, ensuring the prefix matches the one supplied
func Decode(id string, prefix string) (string, error) {
	var normalized [128]byte
	prefixLen, normLen, err := normalizeID(id, &normalized)
	if err != nil {
		return "", err
	}

	if prefix != "" {
		if len(prefix) != prefixLen || string(normalized[:prefixLen]) != prefix {
			return "", fmt.Errorf("prefix %s does not match %s", prefix, string(normalized[:prefixLen]))
		}
	}

	var decoded [16]byte
	n, err := crockfordEncoding.Decode(decoded[:], normalized[prefixLen+1:normLen])
	if err != nil {
		return "", fmt.Errorf("invalid Base32 string: %w", err)
	}
	if n != 16 {
		return "", fmt.Errorf("invalid decoded length: expected 16, got %d", n)
	}

	u, err := uuid.FromBytes(decoded[:])
	if err != nil {
		return "", fmt.Errorf("invalid UUID bytes: %w", err)
	}

	return u.String(), nil
}

// New generates a new UUIDv7 and encodes it with the given prefix
func New(prefix string) (string, error) {
	if !isValidPrefix(prefix) {
		return "", fmt.Errorf("invalid prefix %s", prefix)
	}

	uuidv7, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	var buf [32]byte
	prefixLen := len(prefix)
	copy(buf[0:prefixLen], prefix)
	buf[prefixLen] = '_'
	encodeBase32CrockfordToBuf(buf[prefixLen+1:prefixLen+1+26], uuidv7[:])
	return string(buf[:prefixLen+1+26]), nil
}


