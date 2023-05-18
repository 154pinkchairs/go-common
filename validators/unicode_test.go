package validators

import (
	"strings"
	"testing"
)

func BenchmarkEnsureUTF8Tiny(b *testing.B) {
	tiny := strings.Repeat("a", 1<<5) // ~32B string. One UTF-8 char is 1-4 bytes.

	for i := 0; i < b.N; i++ {
		EnsureUTF8(tiny)
	}
}

// 1 KB string is same as anything in English that is 150-200 words long
// e.g. a review, longer email, 3-4 tweets, etc.
func BenchmarkEnsureUTF8Small(b *testing.B) {
	small := strings.Repeat("a", 1<<10) // ~1KB string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		EnsureUTF8(small)
	}
}

func BenchmarkEnsureUTF8Mid(b *testing.B) {
	mid := strings.Repeat("a", 1<<20) // ~1MB string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		EnsureUTF8(mid)
	}
}

func BenchmarkEnsureUTF8Large(b *testing.B) {
	mid := strings.Repeat("a", 1<<20)  // ~1MB string
	large := strings.Repeat(mid, 1<<2) // ~4MB string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		EnsureUTF8(large)
	}
}

func BenchmarkGEnsureUTF8Tiny(b *testing.B) {
	tiny := strings.Repeat("a", 1<<5) // ~32B string. One UTF-8 char is 1-4 bytes.
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GEnsureUTF8(tiny)
	}
}

func BenchmarkGEnsureUTF8Small(b *testing.B) {
	small := strings.Repeat("a", 1<<10) // ~1KB string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GEnsureUTF8(small)
	}
}

func BenchmarkGEnsureUTF8Mid(b *testing.B) {
	mid := strings.Repeat("a", 1<<20) // ~1MB string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GEnsureUTF8(mid)
	}
}

func BenchmarkGEnsureUTF8Large(b *testing.B) {
	mid := strings.Repeat("a", 1<<20)  // ~1MB string
	large := strings.Repeat(mid, 1<<2) // ~4MB string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GEnsureUTF8(large)
	}
}

func TestGEnsureUTF8(t *testing.T) {
	tiny := strings.Repeat("a", 1<<5) // ~32B string. One UTF-8 char is 1-4 bytes.
	small := strings.Repeat("a", 1<<10)

	tests := []struct {
		name    string
		arg     string
		wantErr bool
	}{
		{"tiny", tiny, false},
		{"small", small, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, got := GEnsureUTF8(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GEnsureUTF8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.arg) {
				t.Errorf("GEnsureUTF8() got = %v, want %v", len(got), len(tt.arg))
			}
		})
	}
}
