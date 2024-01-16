# l4go/buildfs ライブラリ

embed.FSの更新時刻を変更したfs.FSを提供するライブラリーです。  
embed.FSは、更新時刻がtime.Timeのゼロ値に固定されていますが、
この時刻を別の時間に替えることができます。

## embed.FSのゼロ値の更新時刻が起こす問題

embed.FSは、更新時刻(FileInfoのModTime())として、常に`time.Time`型のゼロ値を返します。
そして、ゼロ値を特別扱いするプログラムで、この仕様が問題になります。

例えば、httpモジュールの`http.FileServer()`や`http.ServeContent()`は、更新時刻がゼロ値の時、
`Last-Modified`ヘッダの出力を止めます。  
ということは、embed.FSをそのまま利用すると`Last-Modified`ヘッダ無くなるので、HTTPキャッシュ制御に影響が出てしまいます。

## 利用方法

使い方は単純で、`BuildInFS()`というラッパー関数呼び出すだけです。

### `BuildInFS()`仕様

ラッパー関数の`BuildInFS()`の定義は以下の通りです。

``` go
func BuildInFS(fsys embed.FS, build_time time.Time) fs.FS
```
引数にembed.FSと更新時刻を渡すと、更新時刻が替えられたfs.FSが生成されます。

### サンプルコード

以下のようなコードで、embed.FSの更新時刻を変えることが出来ます。

``` go
//go:embed testfs
var raw_testFS embed.FS

var BuildTime = time.UnixMicro(1702377180976629)
var testFS = buildfs.BuildInFS(raw_testFS, BuildTime)
```

上記のサンプルコードでは`BuildTime`変数として更新時刻を決め打ちにしていますが、
実際の利用時には、Goのgenerateの機能などで適切な時刻を生成してください。  
