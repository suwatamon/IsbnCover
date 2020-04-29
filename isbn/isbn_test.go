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

}
