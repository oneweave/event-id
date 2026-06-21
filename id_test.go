// Copyright 2025 Nadrama Pty Ltd
// SPDX-License-Identifier: Apache-2.0

package eventid

import (
	"regexp"
	"strings"
	"testing"
)

func TestDecodePuidv7(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		want    string
		wantErr bool
	}{
		{
			name:    "valid event-id lowercase",
			id:      "abc_06awcb4f5hzmfey7qwt7s8a6q4",
			want:    "0195c62c-8f2c-7f47-bbc7-bf347ca146b9",
			wantErr: false,
		},
		{
			name:    "valid event-id with uppercase input",
			id:      "ABC_06AWCB4F5HZMFEY7QWT7S8A6Q4",
			want:    "0195c62c-8f2c-7f47-bbc7-bf347ca146b9",
			wantErr: false,
		},
		{
			name:    "valid event-id with 1 char prefix",
			id:      "a_06awcb4f5hzmfey7qwt7s8a6q4",
			want:    "0195c62c-8f2c-7f47-bbc7-bf347ca146b9",
			wantErr: false,
		},
		{
			name:    "valid event-id with 5 char prefix",
			id:      "abcde_06awcb4f5hzmfey7qwt7s8a6q4",
			want:    "0195c62c-8f2c-7f47-bbc7-bf347ca146b9",
			wantErr: false,
		},
		{
			name:    "invalid prefix",
			id:      "123_06AWCB4F5HZMFEY7QWT7S8A6Q4",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid characters in payload",
			id:      "abc_IIIOOO789abcdefghjkmnpqrstvwx",
			want:    "",
			wantErr: true,
		},
		{
			name:    "too short",
			id:      "abc_123",
			want:    "",
			wantErr: true,
		},
		{
			name:    "too long",
			id:      "abc_123456789abcdefghjkmnpqrstvwxyz123",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string",
			id:      "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "missing separator",
			id:      "abc06awcb4f5hzmfey7qwt7s8a6q4",
			want:    "",
			wantErr: true,
		},
		{
			name:    "special characters",
			id:      "abc_!@#$%^&*()_+{}[]|\\:;<>,.?/~`",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.id, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name    string
		uuid    string
		prefix  string
		want    string
		wantErr bool
	}{
		{
			name:    "valid uuid and prefix",
			uuid:    "0195c62c-8f2c-7f47-bbc7-bf347ca146b9",
			prefix:  "abc",
			want:    "abc_06awcb4f5hzmfey7qwt7s8a6q4",
			wantErr: false,
		},
		{
			name:    "valid uuid with spaces",
			uuid:    "  0195c62c-8f2c-7f47-bbc7-bf347ca146b9  ",
			prefix:  "xyz",
			want:    "xyz_06awcb4f5hzmfey7qwt7s8a6q4",
			wantErr: false,
		},
		{
			name:    "valid 1-character prefix",
			uuid:    "0195c62c-8f2c-7f47-bbc7-bf347ca146b9",
			prefix:  "a",
			want:    "a_06awcb4f5hzmfey7qwt7s8a6q4",
			wantErr: false,
		},
		{
			name:    "valid 5-character prefix",
			uuid:    "0195c62c-8f2c-7f47-bbc7-bf347ca146b9",
			prefix:  "abcde",
			want:    "abcde_06awcb4f5hzmfey7qwt7s8a6q4",
			wantErr: false,
		},
		{
			name:    "invalid prefix with numbers",
			uuid:    "0195c62c-8f2c-7f47-bbc7-bf347ca146b9",
			prefix:  "123",
			wantErr: true,
		},
		{
			name:    "invalid prefix with special chars",
			uuid:    "0195c62c-8f2c-7f47-bbc7-bf347ca146b9",
			prefix:  "a#c",
			wantErr: true,
		},
		{
			name:    "invalid uuid format",
			uuid:    "not-a-uuid",
			prefix:  "abc",
			wantErr: true,
		},
		{
			name:    "uuid with uppercase",
			uuid:    "0195C62C-8F2C-7F47-BBC7-BF347CA146B9",
			prefix:  "def",
			want:    "def_06awcb4f5hzmfey7qwt7s8a6q4",
			wantErr: false,
		},
		{
			name:    "empty uuid",
			uuid:    "",
			prefix:  "abc",
			wantErr: true,
		},
		{
			name:    "empty prefix",
			uuid:    "0195c62c-8f2c-7f47-bbc7-bf347ca146b9",
			prefix:  "",
			wantErr: true,
		},
		{
			name:    "prefix too long (6 chars)",
			uuid:    "0195c62c-8f2c-7f47-bbc7-bf347ca146b9",
			prefix:  "abcdef",
			wantErr: true,
		},
		{
			name:    "malformed uuid with correct length",
			uuid:    "0195c62c8f2c7f47bbc7bf347ca146b9aaaaa",
			prefix:  "abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.uuid, tt.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		wantErr bool
	}{
		{
			name:    "valid 3-char prefix",
			prefix:  "abc",
			wantErr: false,
		},
		{
			name:    "valid 1-char prefix",
			prefix:  "x",
			wantErr: false,
		},
		{
			name:    "valid 5-char prefix",
			prefix:  "abcde",
			wantErr: false,
		},
		{
			name:    "invalid prefix with numbers",
			prefix:  "123",
			wantErr: true,
		},
		{
			name:    "invalid prefix with uppercase",
			prefix:  "ABC",
			wantErr: true,
		},
		{
			name:    "invalid prefix with special chars",
			prefix:  "a#c",
			wantErr: true,
		},
		{
			name:    "empty prefix",
			prefix:  "",
			wantErr: true,
		},
		{
			name:    "prefix too long (6 chars)",
			prefix:  "abcdef",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Check that result matches expected pattern
				expectedPattern := `^` + tt.prefix + `_[0-9a-hj-km-np-tv-z]{26}$`
				if matched, _ := regexp.MatchString(expectedPattern, got); !matched {
					t.Errorf("New() = %v, doesn't match expected pattern %v", got, expectedPattern)
				}
				// Check that it has the correct prefix
				if !strings.HasPrefix(got, tt.prefix+"_") {
					t.Errorf("New() = %v, doesn't have prefix %v_", got, tt.prefix)
				}
				// Check that we can decode it back to a UUID
				uuid, err := Decode(got, tt.prefix)
				if err != nil {
					t.Errorf("New() generated invalid ID that can't be decoded: %v", err)
				}
				// Check that the UUID has the correct format
				uuidPattern := `^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
				if matched, _ := regexp.MatchString(uuidPattern, uuid); !matched {
					t.Errorf("New() generated UUID %v doesn't match UUIDv7 pattern", uuid)
				}
			}
		})
	}
}

func TestNewUniqueness(t *testing.T) {
	prefix := "tst"
	generated := make(map[string]bool)
	// Generate 1000 IDs and check they're all unique
	for i := 0; i < 1000; i++ {
		id, err := New(prefix)
		if err != nil {
			t.Fatalf("New() failed: %v", err)
		}
		if generated[id] {
			t.Errorf("New() generated duplicate ID: %v", id)
		}
		generated[id] = true
	}
}

