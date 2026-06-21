// Copyright 2025 Nadrama Pty Ltd
// SPDX-License-Identifier: Apache-2.0

package eventid

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// isValidPrefix checks if a prefix is exactly 3 lowercase letters [a-z].
func isValidPrefix(prefix string) bool {
	if len(prefix) != 3 {
		return false
	}
	return (prefix[0] >= 'a' && prefix[0] <= 'z') &&
		(prefix[1] >= 'a' && prefix[1] <= 'z') &&
		(prefix[2] >= 'a' && prefix[2] <= 'z')
}

// normalizeID normalizes a puidv7 ID by converting to lowercase and removing spaces and hyphens.
func normalizeID(id string) ([29]byte, error) {
	var dst [29]byte
	dstIdx := 0

	for i := 0; i < len(id); i++ {
		c := id[i]
		if c == ' ' || c == '-' || c == '\n' || c == '\r' || c == '\t' {
			continue
		}

		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}

		if dstIdx >= 29 {
			return dst, fmt.Errorf("invalid puidv7 format: too long")
		}
		dst[dstIdx] = c
		dstIdx++
	}

	if dstIdx < 29 {
		return dst, fmt.Errorf("invalid puidv7 format: too short")
	}

	// Validate that the prefix part contains only lowercase letters [a-z]
	for idx := 0; idx < 3; idx++ {
		if dst[idx] < 'a' || dst[idx] > 'z' {
			return dst, fmt.Errorf("invalid prefix in puidv7 ID")
		}
	}

	return dst, nil
}

// Encode encodes a UUID into a prefixed, crockford base32-encoded string
func Encode(uuidStr string, prefix string) (string, error) {
	if !isValidPrefix(prefix) {
		return "", fmt.Errorf("invalid prefix %s", prefix)
	}

	// Clean UUID string by trimming spaces and removing internal whitespace/hyphens
	uuidStr = strings.TrimSpace(uuidStr)
	if strings.Contains(uuidStr, " ") {
		uuidStr = strings.ReplaceAll(uuidStr, " ", "")
	}

	u, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", fmt.Errorf("invalid UUID format %s: %w", uuidStr, err)
	}

	var buf [29]byte
	copy(buf[0:3], prefix)
	encodeBase32CrockfordToBuf(buf[3:], u[:])
	return string(buf[:]), nil
}

// Decode decodes a prefixed, crockford base32-encoded string
// into a UUID string, ensuring the prefix matches the one supplied
func Decode(id string, prefix string) (string, error) {
	norm, err := normalizeID(id)
	if err != nil {
		return "", err
	}

	if prefix != "" {
		if string(norm[0:3]) != prefix {
			return "", fmt.Errorf("prefix %s does not match %s", prefix, string(norm[0:3]))
		}
	}

	var decoded [16]byte
	n, err := crockfordEncoding.Decode(decoded[:], norm[3:])
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

	var buf [29]byte
	copy(buf[0:3], prefix)
	encodeBase32CrockfordToBuf(buf[3:], uuidv7[:])
	return string(buf[:]), nil
}


