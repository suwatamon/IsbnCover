# IsbnCover
Go 言語練習で作ったISBN入力するとその本の表紙を返すWebページ

## 動作確認環境
- Windows10 1909
- go version go1.14.1 windows/amd64
- Python 3.8.2
- Google Chrome 80.0.3987.163

## 事前準備
go 言語と Python をインストール
### Python ライブラリ
バーコード解釈用に pyzbar と pillow をインストール
```pwsh
> pip install pyzbar pillow
```
数字認識用tensorflow とかインストール
```pwsh
> pip install tensorflow keras matplotlib Flask flask-cors
```

### 数字認識用データ学習
```pwsh
> cd script
> python learning.py
```
cnn.h5 が作成される。

## 作成・実行
```pwsh
> go build
> .\IsbnCover.exe
```
として実行しておいて
http://localhost:8888
にwebブラウザでアクセスするとrootページが表示される

以下のいずれかでISBNを送信すると、表紙とAmazon商品ページへのリンクを表示するページが返る

- テキストボックスにISBNを入力
- アップロードフォームからISBNバーコードの写真を送る
- Canvas 要素にISBNを手書き
