package main

import (
	"bytes"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

var visitors struct {
	sync.Mutex
	n int
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

var colorReg = regexp.MustCompile(`^\w*$`)

func handleHi(w http.ResponseWriter, r *http.Request) {
	if !colorReg.MatchString("color") {
		http.Error(w, "Optional color is invalid", http.StatusBadRequest)
		return
	}

	visitors.Lock()
	visitors.n++
	visitNum := visitors.n
	visitors.Unlock()
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.Write([]byte("<h1 style='color: " + r.FormValue("color") + "'>Welcome!</h1>You are visitor number " + fmt.Sprint(visitNum) + "!"))
	// fmt.Fprintf(w, "<h1 style='color: %s'>Welcome!</h1>You are visitor number %d!",
	// 	r.FormValue("color"), visitNum)
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()
	buf.WriteString("<h1 style='color: ")
	buf.WriteString(r.FormValue("color"))
	buf.WriteString(">Welcome!</h1>You are visitor number ")
	b := strconv.AppendInt(buf.Bytes(), int64(visitNum), 10)
	b = append(b, '!')
	w.Write(b)
}

func main() {
	log.Printf("Starting on port 8080")
	http.HandleFunc("/hi", handleHi)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
