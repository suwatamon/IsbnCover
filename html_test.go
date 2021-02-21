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
	cases := []struct {
		input    string
		expected tmplData
	}{
		{"9784065216989", tmplData{Isbn10: "4065216982", Isbn13: "9784065216989"}},
		{"9784621300251", tmplData{Isbn10: "4621300253", Isbn13: "9784621300251"}},
	}

	for _, tc := range cases {
		// execute
		d := setTemplateData(tc.input)
		if d.Isbn10 != tc.expected.Isbn10 {
			t.Errorf("Returnd ISBN10 of %s is not %s: %s\n", tc.input, tc.expected.Isbn10, d.Isbn10)
		}
		if d.Isbn13 != tc.expected.Isbn13 {
			t.Errorf("Returnd ISBN13 of %s is not %s: %s\n", tc.input, tc.expected.Isbn13, d.Isbn13)
		}
	}

}
