## micro profiling two algorithm

***memprofile: go test -run none -bench . -benchtime 3s -benchmem -memprofile mem.prof -gcflags -m=2***

Before:
```
goos: darwin
goarch: arm64
pkg: profiling
BenchmarkAlgorithmOne-8   	 3046042	      1179 ns/op	      53 B/op	       2 allocs/op
BenchmarkAlgorithmTwo-8   	12018472	       299.4 ns/op	       0 B/op	       0 allocs/op
```

After:
```
goos: darwin
goarch: arm64
pkg: profiling
BenchmarkAlgorithmOne-8   	 6032106	       557.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkAlgorithmTwo-8   	12136112	       297.5 ns/op	       0 B/op	       0 allocs/op
```

1. allocation happened when assign &bytes.Buffer{} to io.Reader in io.ReadFull()
```
before: io.ReadFull(input, buf[:end]) 
after: input.Read(buf[:end])
```
2. allocation happened when make a un-constant sized slice: make([]byte, size) => make([]byte, 5) // this is a trick


***cpuprofile: go test -run none -bench . -benchtime 3s  -cpuprofile cpu.prof***
```
         .      1.40s     97:		if bytes.Equal(buf, find) {
```
since algTwo use INDEX instead of function call, so here dont make further optimization
