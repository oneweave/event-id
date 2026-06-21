package eventid

import (
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	uuidStr := "0195c62c-8f2c-7f47-bbc7-bf347ca146b9"
	prefix := "abc"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Encode(uuidStr, prefix)
	}
}

func BenchmarkDecode(b *testing.B) {
	puidStr := "abc06awcb4f5hzmfey7qwt7s8a6q4"
	prefix := "abc"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decode(puidStr, prefix)
	}
}

func BenchmarkNew(b *testing.B) {
	prefix := "abc"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = New(prefix)
	}
}
