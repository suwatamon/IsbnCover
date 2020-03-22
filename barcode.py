from pyzbar.pyzbar import decode
from PIL import Image
import sys

# 画像ファイルの指定
image = input()
while image:
    # バーコードの読取り
    data = decode(Image.open(image))
    # コード内容を出力
    print(data[0][0].decode('utf-8', 'ignore'))
    image = input()

