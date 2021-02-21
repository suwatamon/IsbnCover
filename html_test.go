package main

import "testing"

// TestSubstitutionOfISBN10StringToTemplateData は
// テンプレートにISBN文字列を代入するテスト
func TestSubstitutionOfISBN10StringToTemplateData(t *testing.T) {
	// テストパターンの定義
	cases := []struct {
		input    string
		expected tmplData
	}{
		{"4065216982", tmplData{Isbn10: "4065216982", Isbn13: "9784065216989"}},
		{"4621300253", tmplData{Isbn10: "4621300253", Isbn13: "9784621300251"}},
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

func TestSubstitutionOfISBN13StringToTemplateData(t *testing.T) {
	// テストパターンの定義
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
