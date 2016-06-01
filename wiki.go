// Copyright 2010 The Go Authors. All rights reserved.

//following the tutorial on https://golang.org/doc/articles/wiki/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Page struct { //a wiki page has two data types
	Title string //the title
	Body  []byte //the body consisting of the page information
	//body is a byte type cause that is what the io libraries will use
}

func (p *Page) save() error { //this method saves the page title to a txt file
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600) //returns any error values , which is the return type of WriteFile
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename) //readfile returns both a byte[] and error by default
	if err != nil {                        //catch any errors from ReadFile
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):] //drops the leading "view" compenet
	p, _ := loadPage(title)             //then call loadPage function
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
