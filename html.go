package main

import (
	"IsbnCover/isbn"
	"html/template"
	"log"
	"net/http"
)

func generateHTML(w http.ResponseWriter, isbnArg string) {
	if len(isbnArg) == 13 {
		isbnArg = isbn.Isbn13to10(isbnArg)
	}

	type tmplData struct {
		Isbn string
	}

	t := template.Must(template.ParseFiles("reply.html"))

	// setTemplateData
	d := tmplData{Isbn: isbnArg}

	// テンプレートを描画
	if err := t.ExecuteTemplate(w, "reply.html", d); err != nil {
		log.Fatal(err)
	}
}

func setTemplateData() {}
