// Copyright 2010 The Go Authors. All rights reserved.

//following the tutorial on https://golang.org/doc/articles/wiki/

package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct { //a wiki page has two data types
	Title string //the title
	Body  []byte /*
		//the body consisting of the page information
		//body is a byte type cause that is what the io libraries will use

	*/
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
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//The http.Error function sends a specified HTTP response code (in this case "Internal Server Error") and error message.
	}
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
