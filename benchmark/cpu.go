package benchmark

import (
	"math/rand"
	"sync"
	"sync/atomic"
)

func generateNumbers(totalNumbers int) []int {
	numbers := make([]int, totalNumbers)

	for i := 0; i < totalNumbers; i++ {
		numbers[i] = rand.Intn(totalNumbers)
	}

	return numbers
}

func add(numbers []int) int {
	var v int
	for _, n := range numbers {
		v += n
	}

	return v
}

func addConcurrent(goroutines int, numbers []int) int {
	var v int64
	totalNumbers := len(numbers)
	lastGoroutines := goroutines - 1
	stride := totalNumbers / goroutines

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(g int) {
			var lv int
			defer func() {
				atomic.AddInt64(&v, int64(lv))
				wg.Done()
			}()
			start := g * stride
			end := start + stride
			if g == lastGoroutines {
				end = totalNumbers
			}

			for _, n := range numbers[start:end] {
				lv += n
			}
		}(g)
	}

	wg.Wait()

	return int(v)
}
