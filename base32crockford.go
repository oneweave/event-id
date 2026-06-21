// Copyright 2025 Nadrama Pty Ltd
// SPDX-License-Identifier: Apache-2.0

package eventid

import (
	"encoding/base32"
	"fmt"
)

// crockfordAlphabet is a Base32 alphabet as defined by Douglas Crockford.
// It removes I, L, O, U to avoid confusion.
// https://www.crockford.com/base32.html
const crockfordAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"

// crockfordEncoding is a Base32 encoding schema using the crockfordAlphabet
var crockfordEncoding = base32.NewEncoding(crockfordAlphabet).WithPadding(base32.NoPadding)

// encodeBase32CrockfordToBuf encodes the input bytes to Crockford Base32 in the destination buffer.
// The destination buffer must have at least 26 bytes.
func encodeBase32CrockfordToBuf(dst []byte, data []byte) {
	crockfordEncoding.Encode(dst, data)
}

// encodeBase32Crockford encodes the input bytes to Crockford Base32 string
func encodeBase32Crockford(data []byte) string {
	return crockfordEncoding.EncodeToString(data)
}

// decodeBase32Crockford decodes a Crockford Base32 string to bytes
func decodeBase32Crockford(s string) ([]byte, error) {
	var buf [128]byte
	var normalized []byte
	if len(s) <= len(buf) {
		normalized = buf[:0]
	} else {
		normalized = make([]byte, 0, len(s))
	}

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' || c == '-' || c == '\n' || c == '\r' || c == '\t' {
			continue
		}
		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}
		normalized = append(normalized, c)
	}

	if len(normalized) == 0 {
		return nil, fmt.Errorf("empty string")
	}

	dst := make([]byte, crockfordEncoding.DecodedLen(len(normalized)))
	n, err := crockfordEncoding.Decode(dst, normalized)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, fmt.Errorf("invalid Base32 string")
	}
	return dst[:n], nil
}


