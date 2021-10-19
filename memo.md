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

## .protocファイル
`service/gateway/proto`からコード生成する例

普通にコード生成。-Iはインポートパス。
```
$ protoc -I. \                      
  -I./services/catalog/proto \           
  -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
  --proto_path=services/gateway/proto \               
  --go_out . \                   
  --go_opt module=github.com/dai65527/microservice-handson  \
  --go-grpc_out . \
  --go-grpc_opt module=github.com/dai65527/microservice-handson \
  services/gateway/proto/*.proto
```

getewayのコード（`services/gateway/proto/gateway.pb.gw.go`）の生成
```
$ protoc -I.
    --grpc-gateway_out . \
    -I./services/catalog/proto \
    -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \       
    services/gateway/proto/*.proto
```

## bufによるprotocのコンパイル
buf.gen.yamlを用意する。使用するプラグインとオプションなどを書く。
(spec：https://docs.buf.build/configuration/v1/buf-gen-yaml)
```yaml
version: v1beta1

plugins:
  - name: go
    # path: ./bin/protoc-gen-go # binaryがGOPATHにない場合は指定する
    out: .
    opt: paths=source_relative
  - name: go-grpc
    # path: ./bin/protoc-gen-go-grpc
    out: .
    opt: paths=source_relative
  - name: grpc-gateway
  #   path: ./bin/protoc-gen-grpc-gateway # binaryがGOPATHにない場合は指定する
    out: .
    opt: paths=source_relative
```

buf.yamlを用意する。依存するパッケージなどを書く。
`buf.build/beta/googleapis`の実態はhttps://buf.build/beta/googleapisにある。(gRPC Gatewayの設定を書くのに必要)
(spec: https://docs.buf.build/configuration/v1/buf-yaml)
```yaml
version: v1beta1
name: buf.build/dnakano/microservice-handson
deps:
  - buf.build/beta/googleapis
```

プロジェクトのルートディレクトリで、以下を実行。
依存パッケージをインストール、更新して、buf.lockを生成してくれる。
```
$ buf mod update
```

コード生成。
```
$ buf generate
```

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
$ grpcurl -plaintext localhost:5000 list
dnakano.microservice_handson.db.DBServicegrpc.channelz.v1.Channelz
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

## Clientの生成
- grpc.WithInsecure(): セキュリティの設定をしない。Note that transport security is required unless WithInsecure is set.
- grpc.DialContextはデフォルトでは、Non-Blokingで接続する。
- grpc.WithBlock()でBlockingにする。
- grpc.WithDefaultCallOptions(grpc.WaitForReady(true)): 接続が確立するのをブロックして待つ。（デフォルトでは、設定されていない）
- https://github.com/grpc/grpc/blob/master/doc/wait-for-ready.md

```go
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	conn, err := grpc.DialContext(ctx, "db.db.svc.cluster.local:5000", opts...)
	if err != nil {
		return fmt.Errorf("failed to dial db server: %w", err)
	}
```

## Docker Imageの作成
### Dockerfile
- `COPY go.mod`、`go mod download`してから、`COPY . .`する理由。
	- ソース書き換え時に改めて`go mod download`しないようにするため。
	- キャッシュを使ってbuild時間短縮。
- `gcr.io/distroless/base`
	- runtime用の軽量コンテナ
	- alpineよりよい？（https://blog.unasuke.com/2021/practical-distroless/）

## Authについて
### JWT/JWS/JWK
- 分かりやすい：https://techinfoofmicrosofttech.osscons.jp/index.php?JWS
