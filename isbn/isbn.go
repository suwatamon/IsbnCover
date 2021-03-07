package isbn

import "strconv"

// Isbn13to10 Convert ISBN 13-style to 10-style
func Isbn13to10(isbn13 string) (isbn10 string) {
	isbn10 = isbn13[3:13]
	cd := getCheckDigit10(isbn10)
	isbn10 = isbn10[:9] + cd
	return
}

// Isbn10to13 Convert ISBN 10-style to 13-style
func Isbn10to13(isbn10 string) (isbn13 string) {
	// アルゴリズム: モジュラス10 ウェイト3・1
	isbn13 = "978" + isbn10[0:10]
	cd := getCheckDigit13(isbn13)
	isbn13 = isbn13[0:12] + cd
	return
}

func getCheckDigit10(isbn10 string) (digit string) {
	/// アルゴリズム：モジュラス11 ウェイト10-2
	const MaxWeight = 10
	const MinWeight = 2
	const nWegiht = MaxWeight - MinWeight + 1

	sum := 0
	for idx := 0; idx < nWegiht; idx++ {
		weight := MaxWeight - idx
		digit, _ := strconv.Atoi(isbn10[idx : idx+1])
		sum += weight * digit
	}

	c := 11 - (sum % 11)
	if c == 10 {
		digit = "X"
	} else if c == 11 {
		digit = "0"
	} else {
		digit = strconv.Itoa(c)
	}

	return
}

func getCheckDigit13(isbn13 string) (digit string) {
	const Ndigit = 12
	sum := 0
	for idx := 0; idx < Ndigit; idx += 2 {
		digitOdd, _ := strconv.Atoi(string(isbn13[idx]))
		sum += 1 * digitOdd
		digitEven, _ := strconv.Atoi(string(isbn13[idx+1]))
		sum += 3 * digitEven
	}

	modulus := sum % 10
	if modulus == 0 {
		digit = "0"
	} else {
		digit = strconv.Itoa(10 - modulus)
	}

	return
}
