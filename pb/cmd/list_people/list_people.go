package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"pb/api"

	"google.golang.org/protobuf/proto"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s ADDRESS_BOOK_FILE\n", os.Args[0])
	}
	fname := os.Args[1]

	in, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalln("Error reading file: ", err)
	}

	book := &api.AddressBook{}
	if err := proto.Unmarshal(in, book); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	listPeople(os.Stdout, book)
}

func listPeople(w io.Writer, book *api.AddressBook) {
	for _, p := range book.People {
		writePerson(w, p)
	}
}

func writePerson(w io.Writer, p *api.Person) {
	fmt.Fprintln(w, "Person ID:", p.Id)
	fmt.Fprintln(w, "  Name:", p.Name)
	if p.Email != "" {
		fmt.Fprintln(w, "  E-mail address:", p.Email)
	}

	for _, pn := range p.Phones {
		switch pn.Type {
		case api.Person_MOBILE:
			fmt.Fprint(w, "  Mobile phone #: ")
		case api.Person_HOME:
			fmt.Fprint(w, "  Home phone #: ")
		case api.Person_WORK:
			fmt.Fprint(w, "  Work phone #: ")
		}
		fmt.Fprintln(w, pn.Number)
	}
}
