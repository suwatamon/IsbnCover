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
	d = tmplData{
		Isbn10: "4621300253",
		Isbn13: "9784621300251",
	}
	return
}
