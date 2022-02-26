## Profiling a Long Running Web Service

***Run project:***
```
$ go build
$ GODEBUG=gctrace=1 ./profiling > /dev/null
```


***Benchmark long running program:***
```
$ hey -m POST -c 100 -n 10000 "http://localhost:5000/search?term=trump&cnn=on&bbc=on&nyt=on"
```

***Memprofile:***
```
$ go tool pprof http://0.0.0.0:5000/debug/pprof/allocs
```

***Cpuprofile:***
```
$ go tool pprof http://0.0.0.0:5000/debug/pprof/profile\?seconds\=5
```


Documentation of memory profile options.

    // Useful to see current status of heap.
	-inuse_space  : Allocations live at the time of profile  	** default
	-inuse_objects: Number of bytes allocated at the time of profile

	// Useful to see pressure on heap over time.
	-alloc_space  : All allocations happened since program start
	-alloc_objects: Number of object allocated at the time of profile

