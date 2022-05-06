package main

import "net/http"

var build = "develop"

func main() {
	http.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.ListenAndServe(":8080", nil)
}
