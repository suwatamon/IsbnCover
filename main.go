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
	"sync"
)

const nIsbnMax = 13

var (
	chNumrecogIn  = make(chan string, nIsbnMax)
	chNumrecogOut = make(chan string, nIsbnMax)
	muNumrecog    sync.Mutex

	chBarcodeIn  = make(chan string, 1)
	chBarcodeOut = make(chan string, 1)
	muBarcode    sync.Mutex
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
	const tmpDir = "tmp"
	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		fmt.Fprintln(w, "Temporary directory make error")
		fmt.Println(err)
		return
	}

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

	// ファイルを tmpDir 以下に書き込む
	localFileName := tmpDir + "/" + fileHeader.Filename
	BarcodeFile, err := os.Create(localFileName)
	if err != nil {
		fmt.Fprintln(w, "Can not create temporary file on server")
		log.Fatal(err)
	}
	defer func() {
		BarcodeFile.Close()
		os.Remove(localFileName)
	}()

	size, err := io.Copy(BarcodeFile, file)
	if err != nil {
		fmt.Fprintln(w, "Failed to output uploaded file to server")
		log.Fatal(err)
	}
	fmt.Println("Written bytes", size)

	// アップロードされた画像をバーコードとして解釈
	// Python スクリプトを外部コマンドとして呼び出し
	muBarcode.Lock()
	chBarcodeIn <- localFileName
	isbnFromBarcode := <-chBarcodeOut
	muBarcode.Unlock()

	isbn := strings.TrimSpace(string(isbnFromBarcode))
	generateHTML(w, isbn)
}

func handlerPredict(w http.ResponseWriter, r *http.Request) {
	const PixelSize = 28

	r.ParseForm()
	image := r.FormValue("imageList")

	var u [][]int
	err := json.Unmarshal([]byte(image), &u)
	if err != nil {
		log.Fatal(err)
	}

	isbn := ""
	muNumrecog.Lock()
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
	muNumrecog.Unlock()

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
	defer execpy.Wait()

	go func() {
		for {
			str, ok := <-chIn
			if ok == false {
				stdin.Close()
				return
			}
			io.WriteString(stdin, str+"\n")
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			chOut <- scanner.Text()
		}
	}()
}

func main() {

	http.Handle("/style/",
		http.StripPrefix("/style/",
			http.FileServer(http.Dir("style/"))))
	http.Handle("/script/",
		http.StripPrefix("/script/",
			http.FileServer(http.Dir("script/"))))
	http.HandleFunc("/", handlerRoot)
	http.HandleFunc("/reply", handlerReply)

	// go routine で Pythonスクリプトを起動して
	// channel で やりとりさせる
	go callPyWithChan("barcode.py", chBarcodeIn, chBarcodeOut)
	http.HandleFunc("/barcode", handlerBarcode)

	go callPyWithChan("numrecog.py", chNumrecogIn, chNumrecogOut)
	http.HandleFunc("/predict", handlerPredict)

	http.ListenAndServe(":8888", nil)
}
