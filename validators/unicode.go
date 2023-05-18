package validators

import (
	"runtime"
	"strings"
	"sync"
	"unicode/utf8"
)

type EnsureUTF8Func func(s string, builder *strings.Builder)

func NewRunePool(cap uint16) (*sync.Pool, EnsureUTF8Func) {
	return &sync.Pool{
			New: func() interface{} {
				builder := &strings.Builder{}
				builder.Grow(int(cap))
				return builder
			},
		}, func(s string, builder *strings.Builder) {
			builder.Reset()
			for i, r := range s {
				if r == utf8.RuneError {
					_, size := utf8.DecodeRuneInString(s[i:])
					if size == 1 {
						builder.WriteRune('0')
						continue
					}
				}
				builder.WriteRune(r)
			}
		}
}

// EnsureUTF8 returns a string with invalid UTF-8 characters replaced by the
// Unicode replacement character.
func EnsureUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	v := make([]rune, 0, len(s))
	for i, r := range s {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(s[i:])
			if size == 1 {
				// invalid character, replace with space, zero or similar (but cannot be empty)
				v = append(v, '0')
				continue
			}
		}
		v = append(v, r)
	}
	return string(v)
}

// GEnsureUTF8 has the same functionality as EnsureUTF8 but uses concurrency
// For most cases the difference is collosal and diminishes with the size of the string
/*
* See the benchmark file for more details
* pkg: github.com/italistdev/go-common/validators
cpu: Intel(R) Core(TM) i5-6600K CPU @ 3.50GHz
BenchmarkEnsureUTF8Tiny-4        3181953               399.4 ns/op
BenchmarkEnsureUTF8Small-4       2452778               492.7 ns/op
BenchmarkEnsureUTF8Mid-4           14148             94583 ns/op
BenchmarkEnsureUTF8Large-4          3144            378204 ns/op
BenchmarkGEnsureUTF8Tiny-4      210445300                5.691 ns/op
BenchmarkGEnsureUTF8Small-4     13614853                98.89 ns/op
BenchmarkGEnsureUTF8Mid-4          14260             83527 ns/op
BenchmarkGEnsureUTF8Large-4         3525            343363 ns/op
PASS
ok      github.com/italistdev/go-common/validators      14.325s
*/

func GEnsureUTF8(s string) (error, string) {
	// Create a pool of rune slices
	runesP, ensureUTF8Chunk := NewRunePool(100)

	// If the string is small enough, just use EnsureUTF8
	if len(s) < 100 {
		builder := runesP.Get().(*strings.Builder)
		ensureUTF8Chunk(s, builder)
		result := builder.String()
		runesP.Put(builder)
		return nil, result
	}

	var (
		buffer strings.Builder
		wg     sync.WaitGroup
		mutex  sync.Mutex
	)

	// Divide string into chunks
	numGoroutines := runtime.NumCPU()
	chunkSize := len(s) / numGoroutines

	// Process each chunk concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		start := i * chunkSize
		end := start + chunkSize
		if i == numGoroutines-1 { // last chunk
			end = len(s)
		}
		go func(start, end int) {
			defer wg.Done()
			builder := runesP.Get().(*strings.Builder)
			ensureUTF8Chunk(s[start:end], builder)
			mutex.Lock()
			buffer.WriteString(builder.String())
			mutex.Unlock()
			runesP.Put(builder)
		}(start, end)
	}

	wg.Wait()

	return nil, buffer.String()
}
