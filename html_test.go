package main

import "testing"

// TestSubstitutionOfISBN10StringToTemplateData は
// テンプレートにISBN文字列を代入するテスト
func TestSubstitutionOfISBN10StringToTemplateData(t *testing.T) {
	// isbnArg の定義
	// ISBN-10 : 4621300253
	// ISBN-13 : 978-4621300251
	isbnArg := "4621300253"
	d := setTemplateData(isbnArg)
	// d のチェック
	if d.Isbn10 != "4621300253" {
		t.Errorf("Returnd ISBN10 of %s is not %s: %s\n", isbnArg, "4621300253", d.Isbn10)
	}
	if d.Isbn13 != "9784621300251" {
		t.Errorf("Returnd ISBN13 of %s is not %s: %s\n", isbnArg, "9784621300251", d.Isbn13)
	}
}

func TestSubstitutionOfISBN13StringToTemplateData(t *testing.T) {
	// isbnArg の定義
	// ISBN-10 : 4065216982
	// ISBN-13 : 978-4065216989
	isbnArg := "9784065216989"
	d := setTemplateData(isbnArg)
	// d のチェック
	if d.Isbn10 != "4065216982" {
		t.Errorf("Returnd ISBN10 of %s is not %s: %s\n", isbnArg, "4065216982", d.Isbn10)
	}
	if d.Isbn13 != "9784065216989" {
		t.Errorf("Returnd ISBN13 of %s is not %s: %s\n", isbnArg, "9784065216989", d.Isbn13)
	}
}
