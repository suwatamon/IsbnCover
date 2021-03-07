package isbn

import "testing"

func TestGetCheckDigit10(t *testing.T) {
	isbnExample1 := "4103096101"
	cd1 := getCheckDigit10(isbnExample1)
	if cd1 != "1" {
		t.Errorf("%s is not 1\n", cd1)
	}

	isbnExample2 := "401077830X"
	cd2 := getCheckDigit10(isbnExample2)
	if cd2 != "X" {
		t.Errorf("%s is not X\n", cd2)
	}

	isbnExample3 := "4121006240"
	cd3 := getCheckDigit10(isbnExample3)
	if cd3 != "0" {
		t.Errorf("%s is not 0\n", cd3)
	}
}

func TestIsbn13to10(t *testing.T) {
	isbnExample1 := "9784873117522"
	result13 := Isbn13to10(isbnExample1)
	if result13 != "4873117526" {
		t.Errorf("%s is not 1\n", result13)
	}

}

func TestIsbn10to13(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"4873117526", "9784873117522"},
		{"4065216982", "9784065216989"},
		{"4621300253", "9784621300251"},
	}

	for _, tc := range cases {
		// execute
		result13 := Isbn10to13(tc.input)
		if result13 != tc.expected {
			t.Errorf("Returnd ISBN13 of %s is not %s: %s\n", tc.input, tc.expected, result13)
		}
	}
}

func TestGetCheckDigit13(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"9784873117522", "2"},
		{"9784065216989", "9"},
		{"9784621300251", "1"},
	}
	for _, tc := range cases {
		// execute
		cd := getCheckDigit13(tc.input)
		if cd != tc.expected {
			t.Errorf("Returnd checkdigit of ISBN13 for %s is not %s: %s\n", tc.input, tc.expected, cd)
		}
	}
}
