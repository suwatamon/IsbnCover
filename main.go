package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("root.html"))

	// テンプレートを描画
	if err := t.ExecuteTemplate(w, "root.html", nil); err != nil {
		log.Fatal(err)
	}
}
func handlerReply(w http.ResponseWriter, r *http.Request) {
	isbn := r.FormValue("isbn")
	generateHTML(w, isbn)
}

func generateHTML(w http.ResponseWriter, isbn string) {
	type tmplData struct {
		Isbn string
	}

	t := template.Must(template.ParseFiles("reply.html"))

	if len(isbn) == 13 {
		isbn = isbn13to10(isbn)
	}

	d := tmplData{Isbn: isbn}
	// 動作確認用に解釈したISBN文字列を出力
	fmt.Println(d.Isbn)

	// テンプレートを描画
	if err := t.ExecuteTemplate(w, "reply.html", d); err != nil {
		log.Fatal(err)
	}
}

func isbn13to10(isbn13 string) (isbn10 string) {
	isbn10 = isbn13[3:13]
	cd := getCheckDigit10(isbn10)
	isbn10 = isbn10[:9] + cd
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

func handlerBarcode(w http.ResponseWriter, r *http.Request) {
	// POSTメソッドのみ受け付ける
	if r.Method != "POST" {
		fmt.Fprintln(w, "The method should be POST")
		return
	}
	// アップロードされたファイルを取得
	file, fileHeader, err := r.FormFile("barcode")
	if err != nil {
		fmt.Fprintln(w, "No file detected to be upleaded")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Println("uploaded filename is ", fileHeader.Filename)

	// とりあえず書き込むファイル名は固定かつ jpg 決め打ち
	localFileName := "tmp_Barcode.jpg"

	BarcodeFile, err := os.Create("./" + localFileName)
	if err != nil {
		fmt.Fprintln(w, "Can not create temporary file on server")
		log.Fatal(err)
	}
	defer BarcodeFile.Close()

	size, err := io.Copy(BarcodeFile, file)
	if err != nil {
		fmt.Fprintln(w, "Failed to output uploaded file to server")
		log.Fatal(err)
	}
	fmt.Println("Written bytes", size)

	// アップロードされた画像をバーコードとして解釈
	// Python スクリプトを外部コマンドとして呼び出し
	// 結果は標準出力で返されるバイト列を取得
	isbnFromBarcode, err := exec.Command("py", "barcode.py").Output()
	if err != nil {
		fmt.Fprintln(w, "Barcode image can not be interpreted as ISBN")
		fmt.Println(err)
		return
	}

	// ISBN バーコードの解釈結果を確認出力
	fmt.Println(isbnFromBarcode, string(isbnFromBarcode), strings.TrimSpace(string(isbnFromBarcode)))

	isbn := strings.TrimSpace(string(isbnFromBarcode))
	generateHTML(w, isbn)
}

func handlerPredict(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	image := r.Form.Get("image")
	bitStrAry := strings.Split(image, ",")
	for i := 0; i < 28; i++ {
		for j := 0; j < 28; j++ {
			fmt.Print(bitStrAry[i*28+j])
		}
		fmt.Println()
	}
}

func main() {
	http.HandleFunc("/", handlerRoot)
	http.HandleFunc("/reply", handlerReply)
	http.HandleFunc("/barcode", handlerBarcode)
	http.HandleFunc("/predict", handlerPredict)

	http.ListenAndServe(":8888", nil)
}
