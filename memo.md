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

## gRPC

### ステータスコード
予め決まったステータスコードがある。
https://grpc.github.io/grpc/core/md_doc_statuscodes.html

使う時は、こんな感じでGoのエラーに変換する。

```go
if errors.Is(err, db.ErrAlreadyExists) {
    return nil, status.Error(codes.AlreadyExists, "error msg")
}
```

## ログ出力
### loggerパッケージの使い方

```go
package main

import (
	"github.com/dai65527/microservice-handson/pkg/logger"
)

func main() {
	l, err := logger.New()
	if err != nil {
		panic(err)
	}

	// ログを出力
	l.Info("Hello, World!")
    // {"level":"info","timestamp":1628382925315.32,"message":"Hello, World!"}

	// 名前を付けてログを出力
	clogger := l.WithName("db")
	clogger.Info("Hello, World!")
    // {"level":"info","timestamp":1628382925315.366,"logger":"db","message":"Hello, World!"}
	clogger.WithName("grpc").Info("Hello, World!")
    // {"level":"info","timestamp":1628382925315.372,"logger":"db.grpc","message":"Hello, World!"}

	// 値を追加してログを出力
	clogger.WithValues(
		"key1", "value 1",
		"key2", "value 2",
	).Info("grpc request")
    // {"level":"info","timestamp":1628382925315.416,"logger":"db","message":"grpc request","key1":"value 1","key2":"value 2"}
}
```

## GRPCサーバーをlistenするために書く処理

```go
lis, err := net.Listen("tcp", ":50051")
// 略
s := grpc.NewServer()
pb.RegisterServiceServer(s, &server{})
s.Serve(lis)
```

`server`はgrpcのインターフェイスが実装された構造体。

## Clientの作成
こんな感じ。下はテストヘルパーの例。

```go
func newClient(t *testing.T, port int) proto.DBServiceClient {
	target := fmt.Sprintf("localhost: %d", port)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	require.NoError(t, err)
	return proto.NewDBServiceClient(conn)
}
```

## reflection.Registerは何をしているか
> gRPC を使う上でリフレクションを有効にすると、gRPCurl や Evans といったツールを使う際に Protocol Buffers などの IDL を直接的に読み込まずにメソッドを呼び出すことができてとても便利。
https://syfm.hatenablog.com/entry/2020/06/23/235952
→サーバが持つメソッドなどを参照できるようにする機能？

https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md

## Channelzとは？
> gRPCによる通信の場合は、Channelzという通信状況をデバッグできるツールが用意されています。
https://zenn.dev/imamura_sh/articles/channelz-introduction

## grpcurlで動作確認
- リストアップ
```
$ grpcurl -plaintext localhost:5000 list dnakano.microservice_handson.db.DBServicegrpc.channelz.v1.Channelz
grpc.reflection.v1alpha.ServerReflection
```

- メソッド一覧
```
$ grpcurl -plaintext localhost:5000 list grpc.reflection.v1alpha.ServerReflection
grpc.reflection.v1alpha.ServerReflection.ServerReflectionInfo
```

- メソッド叩いてみる
```
$ grpcurl -plaintext -d '{"name": "bunjiro"}' localhost:5000  dnakano.microservice_handson.db.DBService/CreateCustomer
{
  "customer": {
    "id": "e37e21c6-ffab-4424-b126-4464b5a6ec6e",
    "name": "bunjiro"
  }
}
$ grpcurl -plaintext -d '{"name": "bunjiro"}' localhost:5000  dnakano.microservice_handson.db.DBService/CreateCustomer
ERROR:
  Code: AlreadyExists
  Message: already exists
```

## FormError

```go
func (s *server) GetItem(ctx context.Context, req *proto.GetItemRequest) (*proto.GetItemResponse, error) {
	res, err := s.dbClient.GetItem(ctx, &db.GetItemRequest{Id: req.Id})
	if err != nil {
		// st != nil かつ ok == true の場合？
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}
	[...]
}
```
