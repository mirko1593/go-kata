package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()
var encoder = schema.NewEncoder()

// Person ...
type Person struct {
	Name  string
	Phone string
}

func h(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {

	}

	var p Person

	err = decoder.Decode(&p, req.PostForm)
	if err != nil {

	}

	fmt.Println(p)
}

func makeRequest() {
	person := Person{"Jane Doe", "555-5555"}
	form := url.Values{}

	err := encoder.Encode(person, form)
	if err != nil {

	}

	client := &http.Client{}
	_, err = client.PostForm("http://my-api.test", form)
}
