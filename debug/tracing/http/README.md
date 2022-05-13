### tracing a long running program

```
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5

go tool trace trace.out
```
