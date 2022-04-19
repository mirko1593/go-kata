## http web service test, data race and benchmark


Two ways to test http web service:
1.
```go
func TestHandleRoot_Recorder(t *testing.T) {
        rw := httptest.NewRecorder()
        handleHi(rw, req(t, "GET / HTTP/1.0\r\n\r\n"))
        if !strings.Contains(rw.Body.String(), "visitor number") {
                t.Errorf("Unexpected output: %s", rw.Body)
        }
}
```
2.
```
func TestHandleHi_TestServer(t *testing.T) {
        ts := httptest.NewServer(http.HandlerFunc(handleHi))
        defer ts.Close()
        res, err := http.Get(ts.URL)
        if err != nil {
                t.Error(err)
                return
        }
	if g, w := res.Header.Get("Content-Type"), "text/html; charset=utf-8"; g != w {
		t.Errorf("Content-Type = %q; want %q", g, w)
	}
        slurp, err := ioutil.ReadAll(res.Body)
        defer res.Body.Close()
        if err != nil {
                t.Error(err)
                return
        }
        t.Logf("Got: %s", slurp)
}
```

### Race Detector
Go has concurrency built-in to the language and automatically parallelizes code as necessary over any available CPUs. Unlike Rust, in Go you can write code with a data race if you're not careful. A data race is when multiple goroutine access shared data concurrently without synchronization, when at least one of the gouroutines is doing a write.

```go
func TestHandleHi_TestServer_Parallel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleHi))
	defer ts.Close()
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := http.Get(ts.URL)
			if err != nil {
				t.Error(err)
				return
			}
			if g, w := res.Header.Get("Content-Type"), "text/html; charset=utf-8"; g != w {
				t.Errorf("Content-Type = %q; want %q", g, w)
			}
			slurp, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				t.Error(err)
				return
			}
			t.Logf("Got: %s", slurp)
		}()
	}
	wg.Wait()
}
```
**Fix
```
  var visitors struct {
    sync.Mutex
    n int
  }
...
  func foo() {
    ...
    visitors.Lock()
    visitors.n++
    yourVisitorNumber := visitors.n
    visitors.Unlock()
```

### CPU profiling
```
func BenchmarkHi(b *testing.B) {
        r := req(b, "GET / HTTP/1.0\r\n\r\n")
        rw := httptest.NewRecorder()
        for i := 0; i < b.N; i++ {
                handleHi(rw, r)
        }
}
```

*** $ go test -run none -bench -benchtime 3s -cpuprofile=cpu.prof ***

```
var colorRx = regexp.MustCompile(`\w*$`)
if !colorRx.MatchString(r.FormValue("color")) {
```
=> 10x faster!

### Memory Profiling

*** $ go test -run none -bench -benchtime 3s -benchmem -memprofile=mem.prof ***

1. Remove Content-Type header line (the net/http Server will do it for us)
2. use fmt.Fprintf(w, ... instead of concats

```
   fmt.Fprintf(w, "<h1 style='color: %s'>Welcome!</h1>You are visitor number %d!", r.FormValue("color"), num)
```


### Removing all allocations
```
var bufPool = sync.Pool{
        New: func() interface{} {
                return new(bytes.Buffer)
        },
}
```
. to make a per-processor buffer pool at global scope, and then in the handler:
```
        buf := bufPool.Get().(*bytes.Buffer)
        defer bufPool.Put(buf)
        buf.Reset()
        buf.WriteString("<h1 style='color: ")
        buf.WriteString(r.FormValue("color"))
        buf.WriteString(">Welcome!</h1>You are visitor number ")
        b := strconv.AppendInt(buf.Bytes(), int64(num), 10)
        b = append(b, '!')
        w.Write(b)
```

### Contention Profiling
```
func BenchmarkHiParallel(b *testing.B) {
        r := req(b, "GET / HTTP/1.0\r\n\r\n")
        b.RunParallel(func(pb *testing.PB) {
                rw := httptest.NewRecorder()
                for pb.Next() {
                        handleHi(rw, r)
                        reset(rw)
                }
        })
}
```
And measure:
```
$ go test -bench=Parallel -blockprofile=prof.block
```
And fix:
```
var colorRxPool = sync.Pool{
        New: func() interface{} { return regexp.MustCompile(`\w*$`) },
}
...
func handleHi(w http.ResponseWriter, r *http.Request) {
        if !colorRxPool.Get().(*regexp.Regexp).MatchString(r.FormValue("color")) {
                http.Error(w, "Optional color is invalid", http.StatusBadRequest)
                return
        }
```
And Refactor: 
```
        num := nextVisitorNum()
...
func nextVisitorNum() int {
	visitors.Lock()
	defer visitors.Unlock()
	visitors.n++
	return visitors.n
}
```
And Write some benchmark:
```
func BenchmarkVisitCount(b *testing.B) {
        b.RunParallel(func(pb *testing.PB) {
                for pb.Next() {
                        incrementVisitorNum()
                }
        })
}
```

### Converage
```
$ go test -cover -coverprofile=cover
PASS
coverage: 54.8% of statements
ok      yapc/demo       0.066s
$ go tool cover -html=cover
(opens web browser)
```
