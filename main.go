package main

import (
	"bufio"
	"encoding/json"
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

const nIsbnMax = 13

var (
	chNumrecogIn  chan string
	chNumrecogOut chan string
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

func handlerPredictCh(w http.ResponseWriter, r *http.Request) {
	const PixelSize = 28

	r.ParseForm()
	image := r.Form.Get("imageList")

	var u [][]int
	err := json.Unmarshal([]byte(image), &u)
	if err != nil {
		log.Fatal(err)
	}

	isbn := ""
	for _, ii := range u {
		str := fmt.Sprintf("%v", ii)
		// 先頭と最後の1文字ずつ([])を取り除く
		str = str[1 : len(str)-1]
		chNumrecogIn <- str
	}
	for i := 0; i < len(u); i++ {
		numPredicted := <-chNumrecogOut
		isbn += numPredicted
	}
	fmt.Printf("結果: %s\n", isbn)

	generateHTML(w, isbn)
}

func callPyWithChan(pyScript string, chIn <-chan string, chOut chan<- string) {
	execpy := exec.Command("py", pyScript)
	stdin, err := execpy.StdinPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	stdout, err := execpy.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	stderr, err := execpy.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	scannerErr := bufio.NewScanner(stderr)
	go func() {
		for scannerErr.Scan() {
			fmt.Fprintln(os.Stderr, scannerErr.Text())
		}
	}()

	execpy.Start()
	defer func() {
		stdin.Close()
		execpy.Wait()
	}()

	go func() {
		for {
			str, ok := <-chIn
			if ok == false {
				return
			}
			io.WriteString(stdin, str+"\n")
		}
	}()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		chOut <- scanner.Text()
	}
}

func main() {

	// go routine で Pythonスクリプトを起動して
	// channel で やりとりさせる
	chNumrecogIn = make(chan string, nIsbnMax)
	chNumrecogOut = make(chan string, nIsbnMax)
	go callPyWithChan("numrecog.py", chNumrecogIn, chNumrecogOut)

	http.Handle("/style/",
		http.StripPrefix("/style/",
			http.FileServer(http.Dir("style/"))))
	http.HandleFunc("/", handlerRoot)
	http.HandleFunc("/reply", handlerReply)
	http.HandleFunc("/barcode", handlerBarcode)
	http.HandleFunc("/predict", handlerPredictCh)

	http.ListenAndServe(":8888", nil)
}
