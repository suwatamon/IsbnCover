package main

import (
	"IsbnCover/isbn"
	"html/template"
	"log"
	"net/http"
)

type tmplData struct {
	Isbn10 string
	Isbn13 string
}

func generateHTML(w http.ResponseWriter, isbnArg string) {
	if len(isbnArg) == 13 {
		isbnArg = isbn.Isbn13to10(isbnArg)
	}

	t := template.Must(template.ParseFiles("reply.html"))

	// setTemplateData
	d := setTemplateData(isbnArg)

	// テンプレートを描画
	if err := t.ExecuteTemplate(w, "reply.html", d); err != nil {
		log.Fatal(err)
	}
}

func setTemplateData(isbnArg string) (d tmplData) {
	var isbn10, isbn13 string
	if len(isbnArg) == 10 {
		isbn10 = isbnArg
		isbn13 = "9784621300251"
	} else {
		isbn13 = isbnArg
		isbn10 = isbn.Isbn13to10(isbnArg)
	}
	d = tmplData{
		Isbn10: isbn10,
		Isbn13: isbn13,
	}

	return
}
