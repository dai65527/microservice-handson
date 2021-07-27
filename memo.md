# Memo

## protobuf + Goプラグインのインストール
Macです。

```sh
$ brew install protobuf
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
$ export PATH="$PATH:$(go env GOPATH)/bin"
```

## .protoファイルからGoのコードの生成
`service/item/proto`からコード生成する例。

```
$ protoc --proto_path=services/item/proto services/item/proto/*.proto --go_out=. --go_opt=module=github.com/dai65527/microservice-handson --go-grpc_out=. --go-grpc_opt=module=github.com/dai65527/microservice-handson
```

`--go_opt`と`--go-grpc_opt`で出力先を指定する。無いと、モジュール名も出力パスに含まれてしまう。例えば、

```
option go_package = "github.com/dai65527/microservice-handson/services/item/proto";
```

としていた場合に、`github.com/dai65527/microservice-handson/services/item/proto`というディレクトリに出力されてしまう。なので、`--go_opt=module=github.com/dai65527/microservice-handson`としておく。

## go modと必要なpackageの追加

```sh
$ go mod init microservice-handson  # go.modの生成
$ go mod tidy                       # 使用されているパッケージを全て追加する
```

## パッケージをバージョン指定してインストールし直す
`github.com/go-logr/logr`のバージョンが異なるので、@v0.4.0を入れ直す。

```
$ go get github.com/go-logr/logr@v0.4.0
```
