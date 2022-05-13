package benchmark

import (
	"runtime"
	"testing"
)

var numbers []int

func init() {
	numbers = generateNumbers(1e7)
}

func BenchmarkSequentialAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		add(numbers)
	}
}

func BenchmarkConcurrentAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addConcurrent(runtime.NumCPU(), numbers)
	}
}

func BenchmarkConcurrentAddGoroutines(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addConcurrent(10000, numbers)
	}
}
