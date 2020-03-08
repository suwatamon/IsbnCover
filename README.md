# IsbnCover
Go 言語練習で作ったISBN入力するとその本の表紙を返すWebページ

## 動作確認環境
- Windows10 1909
- go version go1.13.4 windows/amd64
- Python 3.8.0

## 事前準備
go 言語と Python をインストール
### Python ライブラリ
バーコード解釈用に pyzbar と pillow をインストール
```powershell:install command
pip install pyzbar pillow
```
数字認識用tensorflow とかインストール
```
pip install tensorflow==2.0.0 keras matplotlib Flask flask-cors
```

## 作成・実行
```powershell:build command
go build
./barcode.exe
```
として実行しておいて
http://localhost:8888
にwebブラウザでアクセスするとrootページが表示される

テキストボックスにISBNを入力するか、
アップロードフォームからISBNバーコードの写真を送ると
表紙とAmazon商品ページへのリンクを表示するページが返る
