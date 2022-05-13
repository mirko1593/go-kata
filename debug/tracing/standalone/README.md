## Tracing

cpu profile
```
    pprof.StartCPUProfile(os.Stdout)
	defer pprof.StopCPUProfile()
```

trace
```
	trace.Start(os.Stdout)
	defer trace.Stop()
```

generate trace data
```
    $ time ./trace > t.out
    $ go tool trace t.out
```

notes:
three implementation: sequential | fan-out pattern | bounded fan-out
generate difference trace data and head allocation.
