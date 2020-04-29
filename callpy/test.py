import sys
data = input()
while data:
    # コード内容を出力
    print(data)
    print(data, file=sys.stderr)
    data = input()