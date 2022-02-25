## benchmark

#### io with 1 core
**$ GOGC=off go test -cpu 1 -run none -bench . -benchtime 3s**
```
BenchmarkSequential               	       3	1357808278 ns/op
BenchmarkConcurrent               	      22	 151461828 ns/op
BenchmarkConcurrentGoroutine      	     168	  21234375 ns/op
BenchmarkSequentialAgain          	       3	1315390361 ns/op
BenchmarkConcurrentAgain          	      24	 150052781 ns/op
BenchmarkConcurrentGoroutineAgain 	     171	  20891343 ns/op
```
=> even with 1 core, more goroutines can gain better performance: 1 -> runtime.NumCPU() -> 1e3(# of workload)

#### io with 8 core
**$ GOGC=off go test -cpu 8 -run none -bench . -benchtime 3s**
```
goos: darwin
goarch: arm64
pkg: benchmark
BenchmarkSequential-8                 	       3	1367239111 ns/op
BenchmarkConcurrent-8                 	      21	 167059290 ns/op
BenchmarkConcurrentGoroutine-8        	     572	   5890108 ns/op
BenchmarkSequentialAgain-8            	       3	1331261264 ns/op
BenchmarkConcurrentAgain-8            	      22	 163396578 ns/op
BenchmarkConcurrentGoroutineAgain-8   	     601	   6133793 ns/op
```
=> case # of goroutines <= # of cores: bringing in the extra OS/hardware threads don’t provide any better performance <br/>
=> case # of goroutines > # of cores: bring more OS/hardware threads can gain better performance


#### cpu with 1 core: 1 | runtime.NumCPU() | 10000
```
goos: darwin
goarch: amd64
pkg: benchmark
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkSequentialAdd                       542           7029158 ns/op
BenchmarkConcurrentAdd                       228          17076309 ns/op
BenchmarkConcurrentAddGoroutines             136          25363133 ns/op

```
=> add more goroutines will loose performance for the sake of context switch

### cpu with 12 core: 1 goroutine | runtime.NumCPU() | 10000
```
goos: darwin
goarch: amd64
pkg: benchmark
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkSequentialAdd-12                    528           6659990 ns/op
BenchmarkConcurrentAdd-12                   1183           3035956 ns/op
BenchmarkConcurrentAddGoroutines-12          930           3779810 ns/op
```
=> add more threads while keep 1:1 to goroutines will improve performance

