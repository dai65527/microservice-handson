# Memo

## protobuf + Go プラグインのインストール

Mac です。

```sh
$ brew install protobuf
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
$ export PATH="$PATH:$(go env GOPATH)/bin"
```

## .proto ファイルから Go のコードの生成

`service/item/proto`からコード生成する例。

```
$ protoc --proto_path=services/item/proto services/item/proto/*.proto --go_out=. --go_opt=module=github.com/dai65527/microservice-handson --go-grpc_out=. --go-grpc_opt=module=github.com/dai65527/microservice-handson
```

`--go_opt`と`--go-grpc_opt`で出力先を指定する。無いと、モジュール名も出力パスに含まれてしまう。例えば、

```
option go_package = "github.com/dai65527/microservice-handson/services/item/proto";
```

としていた場合に、`github.com/dai65527/microservice-handson/services/item/proto`というディレクトリに出力されてしまう。なので、`--go_opt=module=github.com/dai65527/microservice-handson`としておく。

## .protoc ファイル

`service/gateway/proto`からコード生成する例

普通にコード生成。-I はインポートパス。

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

geteway のコード（`services/gateway/proto/gateway.pb.gw.go`）の生成

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

## buf による protoc のコンパイル

buf.gen.yaml を用意する。使用するプラグインとオプションなどを書く。
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

buf.yaml を用意する。依存するパッケージなどを書く。
`buf.build/beta/googleapis`の実態はhttps://buf.build/beta/googleapisにある。(gRPC Gateway の設定を書くのに必要)
(spec: https://docs.buf.build/configuration/v1/buf-yaml)

```yaml
version: v1beta1
name: buf.build/dnakano/microservice-handson
deps:
  - buf.build/beta/googleapis
```

プロジェクトのルートディレクトリで、以下を実行。
依存パッケージをインストール、更新して、buf.lock を生成してくれる。

```
$ buf mod update
```

コード生成。

```
$ buf generate
```

## go mod と必要な package の追加

```sh
$ go mod init microservice-handson  # go.modの生成
$ go mod tidy                       # 使用されているパッケージを全て追加する
```

## パッケージをバージョン指定してインストールし直す

`github.com/go-logr/logr`のバージョンが異なるので、@v0.4.0 を入れ直す。

```
$ go get github.com/go-logr/logr@v0.4.0
```

## gRPC

### ステータスコード

予め決まったステータスコードがある。
https://grpc.github.io/grpc/core/md_doc_statuscodes.html

使う時は、こんな感じで Go のエラーに変換する。

```go
if errors.Is(err, db.ErrAlreadyExists) {
    return nil, status.Error(codes.AlreadyExists, "error msg")
}
```

## ログ出力

### logger パッケージの使い方

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

## GRPC サーバーを listen するために書く処理

```go
lis, err := net.Listen("tcp", ":50051")
// 略
s := grpc.NewServer()
pb.RegisterServiceServer(s, &server{})
s.Serve(lis)
```

`server`は grpc のインターフェイスが実装された構造体。

## Client の作成

こんな感じ。下はテストヘルパーの例。

```go
func newClient(t *testing.T, port int) proto.DBServiceClient {
	target := fmt.Sprintf("localhost: %d", port)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	require.NoError(t, err)
	return proto.NewDBServiceClient(conn)
}
```

## reflection.Register は何をしているか

> gRPC を使う上でリフレクションを有効にすると、gRPCurl や Evans といったツールを使う際に Protocol Buffers などの IDL を直接的に読み込まずにメソッドを呼び出すことができてとても便利。
> https://syfm.hatenablog.com/entry/2020/06/23/235952
> → サーバが持つメソッドなどを参照できるようにする機能？

https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md

## Channelz とは？

> gRPC による通信の場合は、Channelz という通信状況をデバッグできるツールが用意されています。
> https://zenn.dev/imamura_sh/articles/channelz-introduction

## grpcurl で動作確認

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

## Client の生成

- grpc.WithInsecure(): セキュリティの設定をしない。Note that transport security is required unless WithInsecure is set.
- grpc.DialContext はデフォルトでは、Non-Bloking で接続する。
- grpc.WithBlock()で Blocking にする。
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

## Docker Image の作成

### Dockerfile

- `COPY go.mod`、`go mod download`してから、`COPY . .`する理由。
  - ソース書き換え時に改めて`go mod download`しないようにするため。
  - キャッシュを使って build 時間短縮。
- `gcr.io/distroless/base`
  - runtime 用の軽量コンテナ
  - alpine よりよい？（https://blog.unasuke.com/2021/practical-distroless/）

## Auth について

### JWT/JWS/JWK

- 分かりやすい：https://techinfoofmicrosofttech.osscons.jp/index.php?JWS

## gRPC Gateway 経由のリクエスト（REST）

SignUp

```sh
$ curl -s -XPOST -d '{"name":"gopher"}' localhost:4000/auth/signup | jq .
{
  "customer": {
    "id": "46cafbac-fb9d-4bdf-bdb4-d7010dc82803",
    "name": "gopher"
  }
}
$ curl -s -XPOST -d '{"name":"gopher"}' localhost:4000/auth/signin | jq .
{
  "access_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImFhN2M2Mjg3LWM0NWQtNDk2Ni04NGI0LWExNjMzZTRlM2E2NCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRob3JpdHkiLCJzdWIiOiI0NmNhZmJhYy1mYjlkLTRiZGYtYmRiNC1kNzAxMGRjODI4MDMifQ.GKqJXCgsTr_yqrVcxg6zBDRXtZYuxxvRO8dV3GbgljVzO2JnQuhg2jwThCFsuY0m5RlU0DKS5X9VI7ZTL8mOlowUaoypgschJMj3juX2JiP9Daj-6w6k__zCAdaJ4eiw8QdoljYhrfRi8Q4jei2JFqTODqpC8cgmjvTACDdl-MQ4_v6Qrf925ZlW6LFuY6QIzhfmlg9a6486RjpxDBCv-BaR3k8IHl5GcE-3ACTDVq_N25ePH7zha917FROQnflTy9Ozy-_V37khiiPtu6TdkC7bf2mA4VMefR_crh5NNG90Ch2IQ1NhYKMg1l400n_7U1aiSjnImv4TqtvJ5s0Www"
}
$ export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6ImFhN2M2Mjg3LWM0NWQtNDk2Ni04NGI0LWExNjMzZTRlM2E2NCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRob3JpdHkiLCJzdWIiOiI0NmNhZmJhYy1mYjlkLTRiZGYtYmRiNC1kNzAxMGRjODI4MDMifQ.GKqJXCgsTr_yqrVcxg6zBDRXtZYuxxvRO8dV3GbgljVzO2JnQuhg2jwThCFsuY0m5RlU0DKS5X9VI7ZTL8mOlowUaoypgschJMj3juX2JiP9Daj-6w6k__zCAdaJ4eiw8QdoljYhrfRi8Q4jei2JFqTODqpC8cgmjvTACDdl-MQ4_v6Qrf925ZlW6LFuY6QIzhfmlg9a6486RjpxDBCv-BaR3k8IHl5GcE-3ACTDVq_N25ePH7zha917FROQnflTy9Ozy-_V37khiiPtu6TdkC7bf2mA4VMefR_crh5NNG90Ch2IQ1NhYKMg1l400n_7U1aiSjnImv4TqtvJ5s0Www
```

```sh
$ curl -s -XGET -H "authorization: bearer $TOKEN" localhost:4000/catalog/items | jq .
```

## 21.11.3 Gateway からの認証がうまくいかない

```sh
$ docker-compose up
WARNING: Found orphan containers (handson-db-debug) for this project. If you removed or renamed this service in your compose file, you can run this command with the --remove-orphans flag to clean it up.
Starting handson-item-service      ... done
Starting handson-db-service        ... done
Starting handson-catalog-service   ... done
Starting handson-authority-service ... done
Starting handson-customer-service  ... done
Recreating handson-gateway-service ... done
Attaching to handson-db-service, handson-catalog-service, handson-customer-service, handson-item-service, handson-authority-service, handson-gateway-service
handson-gateway-service | {"level":"info","timestamp":1635893034987.4788,"logger":"gateway.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.gateway.GatewayService/CreateItem","request_id":"49ebb263-155e-4e1c-8d4a-3c234d385f2d"}
handson-gateway-service | token: eyJhbGciOiJSUzI1NiIsImtpZCI6ImFhN2M2Mjg3LWM0NWQtNDk2Ni04NGI0LWExNjMzZTRlM2E2NCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRob3JpdHkiLCJzdWIiOiIxMzM0YTFjOC02YmQzLTQ4YmYtOWNmYy05OTRiODFjOWJkNDMifQ.TSV5mSo3ifkv7ELzEYfySfV17nTsP4NACP3lwxTqmh7GvUITySJ6hOCTKMYD_VJjIISpVSWfP4niBV3CB2e9hUv8f96bWxfEAlSDHVcsG5lI9Se_shuscVIKOQJXhjSOSre_5ei4AFumLcmcuIhjJADki5e81Wh3Ic5yrPj4h9NDMv_HatXW4BGIaBU4XZCufbWyehoriWtTvbjV5bSZhZ1AaoO8PdDqzvyrHvj6agAelb5rjsgTT5jK2flOsS19c6u4usOiwwKvAkNgQieLi2HWjXezkPbb1BRAjepiLbbyC-9pTDtOQLKoKDU5NW5aUP7Mu9DzTNtikZvbJAclOw
handson-gateway-service | res.Jwks: {"keys":[{"e":"AQAB","kid":"aa7c6287-c45d-4966-84b4-a1633e4e3a64","kty":"RSA","n":"tWT6dFoWKoWI1dSO6FN6FDFMFHAgvR4j6QEuVE-3Q2_iAnjXE3jCawBqYT3vCH6azWksLhsT12WjaNp3OjLP2yuDke2aVzD9c9g274XWvcJ90YI31GHdWbuZ1MgA9_Gmkq_v-cvSrkRctZAs_ktoHT4Fnpzfyl4fphCsYptL2cfe_PyJi2g6PIyrd21uhiBaHkWd8m3MBBl2yUHMqOlUJQ9NtlomAoKmqNIeGo7sfNPUtXlQ_cjQzP1aIotByu5N_PGpzLIqzDlGonytB55fOPEuJ1BXt3RxaxU_dje4LdSVlZk7Fms3XJ0Ug1r7PQE-ds7wDRkUb0vUDP9PALWrbQ"}]}
handson-gateway-service | key.Len(): 1
handson-gateway-service | {"level":"info","timestamp":1635893034992.5544,"logger":"gateway.grpc.server","message":"failed to verify token: failed to find matching key for verification: invalid signature algorithm : invalid jwa.SignatureAlgorithm value","request_id":"49ebb263-155e-4e1c-8d4a-3c234d385f2d"}
handson-gateway-service | {"level":"info","timestamp":1635893034992.7283,"logger":"gateway.grpc.request","message":"finished","method":"/dnakano.microservice_handson.gateway.GatewayService/CreateItem","code":"Unauthenticated","request_id":"49ebb263-155e-4e1c-8d4a-3c234d385f2d"}
handson-authority-service | {"level":"info","timestamp":1635893034989.5896,"logger":"authority.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.authority.AuthorityService/ListPublicKeys","request_id":"6be367ca-1024-42c7-9829-ab283eb64789"}
handson-authority-service | {"level":"info","timestamp":1635893034990.2715,"logger":"authority.grpc.request","message":"finished","method":"/dnakano.microservice_handson.authority.AuthorityService/ListPublicKeys","code":"OK","request_id":"6be367ca-1024-42c7-9829-ab283eb64789"}
handson-db-service | {"level":"info","timestamp":1635893052464.8687,"logger":"db.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.db.DBService/GetCustomerByName","request_id":"d8498661-c7c0-4af4-8a62-a1aa7a4b0999"}
handson-db-service | {"level":"info","timestamp":1635893052465.0283,"logger":"db.grpc.request","message":"finished","method":"/dnakano.microservice_handson.db.DBService/GetCustomerByName","code":"NotFound","request_id":"d8498661-c7c0-4af4-8a62-a1aa7a4b0999"}
handson-authority-service | {"level":"info","timestamp":1635893052457.8284,"logger":"authority.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.authority.AuthorityService/Signin","request_id":"d535d49f-c233-42a1-8595-c946da29466f"}
handson-gateway-service | {"level":"info","timestamp":1635893052456.3914,"logger":"gateway.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.gateway.GatewayService/Signin","request_id":"529a37c6-4757-4562-bb1d-4d3f5555375e"}
handson-customer-service | {"level":"info","timestamp":1635893052461.0435,"logger":"customer.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.customer.CustomerService/GetCustomerByName","request_id":"ba68b354-7bfb-4364-934d-e5277678db85"}
handson-customer-service | {"level":"info","timestamp":1635893052468.113,"logger":"customer.grpc.request","message":"finished","method":"/dnakano.microservice_handson.customer.CustomerService/GetCustomerByName","code":"AlreadyExists","request_id":"ba68b354-7bfb-4364-934d-e5277678db85"}
handson-authority-service | {"level":"info","timestamp":1635893052470.805,"logger":"authority.grpc.server","message":"failed to authenticate the customer: rpc error: code = AlreadyExists desc = not found","request_id":"d535d49f-c233-42a1-8595-c946da29466f"}
handson-authority-service | {"level":"info","timestamp":1635893052470.9062,"logger":"authority.grpc.request","message":"finished","method":"/dnakano.microservice_handson.authority.AuthorityService/Signin","code":"Unauthenticated","request_id":"d535d49f-c233-42a1-8595-c946da29466f"}
handson-gateway-service | {"level":"info","timestamp":1635893052471.6965,"logger":"gateway.grpc.request","message":"finished","method":"/dnakano.microservice_handson.gateway.GatewayService/Signin","code":"Unauthenticated","request_id":"529a37c6-4757-4562-bb1d-4d3f5555375e"}
handson-gateway-service | {"level":"info","timestamp":1635893063928.7341,"logger":"gateway.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.gateway.GatewayService/Signup","request_id":"21d12419-9138-461d-b2a2-85dc3b1620aa"}
handson-authority-service | {"level":"info","timestamp":1635893063930.1057,"logger":"authority.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.authority.AuthorityService/Signup","request_id":"f4175032-8565-4215-994c-139ad3241438"}
handson-customer-service | {"level":"info","timestamp":1635893063931.3923,"logger":"customer.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.customer.CustomerService/CreateCustomer","request_id":"e6312800-b277-4fc3-8abf-b22e0fdedf68"}
handson-gateway-service | {"level":"info","timestamp":1635893063937.1272,"logger":"gateway.grpc.request","message":"finished","method":"/dnakano.microservice_handson.gateway.GatewayService/Signup","code":"OK","request_id":"21d12419-9138-461d-b2a2-85dc3b1620aa"}
handson-authority-service | {"level":"info","timestamp":1635893063935.4856,"logger":"authority.grpc.request","message":"finished","method":"/dnakano.microservice_handson.authority.AuthorityService/Signup","code":"OK","request_id":"f4175032-8565-4215-994c-139ad3241438"}
handson-db-service | {"level":"info","timestamp":1635893063932.687,"logger":"db.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.db.DBService/CreateCustomer","request_id":"880870f2-6313-4776-8fbb-3b5559cfe1b6"}
handson-db-service | {"level":"info","timestamp":1635893063932.747,"logger":"db.grpc.request","message":"finished","method":"/dnakano.microservice_handson.db.DBService/CreateCustomer","code":"OK","request_id":"880870f2-6313-4776-8fbb-3b5559cfe1b6"}
handson-customer-service | {"level":"info","timestamp":1635893063933.9233,"logger":"customer.grpc.request","message":"finished","method":"/dnakano.microservice_handson.customer.CustomerService/CreateCustomer","code":"OK","request_id":"e6312800-b277-4fc3-8abf-b22e0fdedf68"}
handson-gateway-service | {"level":"info","timestamp":1635893072381.1436,"logger":"gateway.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.gateway.GatewayService/Signin","request_id":"1fac336f-b261-4952-b2aa-a971aebba5fd"}
handson-authority-service | {"level":"info","timestamp":1635893072382.5774,"logger":"authority.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.authority.AuthorityService/Signin","request_id":"6829cd91-0321-43b8-bfe8-a68beb6b8003"}
handson-customer-service | {"level":"info","timestamp":1635893072383.5532,"logger":"customer.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.customer.CustomerService/GetCustomerByName","request_id":"7af11b35-1800-4b2a-ba50-9ea1882f1e5d"}
handson-db-service | {"level":"info","timestamp":1635893072386.2737,"logger":"db.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.db.DBService/GetCustomerByName","request_id":"031161eb-9086-499b-9d2e-37fbe681dbd6"}
handson-db-service | {"level":"info","timestamp":1635893072386.4258,"logger":"db.grpc.request","message":"finished","method":"/dnakano.microservice_handson.db.DBService/GetCustomerByName","code":"OK","request_id":"031161eb-9086-499b-9d2e-37fbe681dbd6"}
handson-customer-service | {"level":"info","timestamp":1635893072391.382,"logger":"customer.grpc.request","message":"finished","method":"/dnakano.microservice_handson.customer.CustomerService/GetCustomerByName","code":"OK","request_id":"7af11b35-1800-4b2a-ba50-9ea1882f1e5d"}
handson-authority-service | {"level":"info","timestamp":1635893072399.0325,"logger":"authority.grpc.request","message":"finished","method":"/dnakano.microservice_handson.authority.AuthorityService/Signin","code":"OK","request_id":"6829cd91-0321-43b8-bfe8-a68beb6b8003"}
handson-gateway-service | {"level":"info","timestamp":1635893072402.0679,"logger":"gateway.grpc.request","message":"finished","method":"/dnakano.microservice_handson.gateway.GatewayService/Signin","code":"OK","request_id":"1fac336f-b261-4952-b2aa-a971aebba5fd"}
handson-gateway-service | {"level":"info","timestamp":1635893113116.1245,"logger":"gateway.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.gateway.GatewayService/Signin","request_id":"a6206893-b1fb-4506-9d9c-a0929e47736a"}
handson-customer-service | {"level":"info","timestamp":1635893113118.7732,"logger":"customer.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.customer.CustomerService/GetCustomerByName","request_id":"4e78a4ec-892e-42d0-8cf3-ee99772fe327"}
handson-authority-service | {"level":"info","timestamp":1635893113117.5051,"logger":"authority.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.authority.AuthorityService/Signin","request_id":"14cafbe6-f8e1-4ce5-b267-3e07ca9492e7"}
handson-db-service | {"level":"info","timestamp":1635893113120.35,"logger":"db.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.db.DBService/GetCustomerByName","request_id":"c9a25287-56cf-4749-a03f-b0f216cb0fad"}
handson-db-service | {"level":"info","timestamp":1635893113120.8394,"logger":"db.grpc.request","message":"finished","method":"/dnakano.microservice_handson.db.DBService/GetCustomerByName","code":"OK","request_id":"c9a25287-56cf-4749-a03f-b0f216cb0fad"}
handson-customer-service | {"level":"info","timestamp":1635893113121.808,"logger":"customer.grpc.request","message":"finished","method":"/dnakano.microservice_handson.customer.CustomerService/GetCustomerByName","code":"OK","request_id":"4e78a4ec-892e-42d0-8cf3-ee99772fe327"}
handson-authority-service | {"level":"info","timestamp":1635893113128.0361,"logger":"authority.grpc.request","message":"finished","method":"/dnakano.microservice_handson.authority.AuthorityService/Signin","code":"OK","request_id":"14cafbe6-f8e1-4ce5-b267-3e07ca9492e7"}
handson-gateway-service | {"level":"info","timestamp":1635893113128.869,"logger":"gateway.grpc.request","message":"finished","method":"/dnakano.microservice_handson.gateway.GatewayService/Signin","code":"OK","request_id":"a6206893-b1fb-4506-9d9c-a0929e47736a"}
handson-authority-service | {"level":"info","timestamp":1635893125624.3455,"logger":"authority.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.authority.AuthorityService/ListPublicKeys","request_id":"10e0df09-2446-4160-b20a-19a83ed7e179"}
handson-authority-service | {"level":"info","timestamp":1635893125624.6262,"logger":"authority.grpc.request","message":"finished","method":"/dnakano.microservice_handson.authority.AuthorityService/ListPublicKeys","code":"OK","request_id":"10e0df09-2446-4160-b20a-19a83ed7e179"}
handson-gateway-service | {"level":"info","timestamp":1635893125623.3635,"logger":"gateway.grpc.request","message":"grpc request","method":"/dnakano.microservice_handson.gateway.GatewayService/CreateItem","request_id":"dd1b59e3-3426-462a-8cd0-c98d558e4034"}
handson-gateway-service | token: eyJhbGciOiJSUzI1NiIsImtpZCI6ImFhN2M2Mjg3LWM0NWQtNDk2Ni04NGI0LWExNjMzZTRlM2E2NCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRob3JpdHkiLCJzdWIiOiI3ODVlODQ0OC1mZDQ2LTQ0NmQtYTdhMC04M2Y1Yjg3MzgzMDgifQ.AP7D5BvHXLkUxrWMQjtaNmocl-B3rJ1eLAg3Xm9BEFhdGbI4S108Dbc-RIxXGAW8tfREeUi1KvcwsMTElML-yxr8jt7NrNF7vA9iv6vU8IFVS3dH6Es6cJdJ1AujjbYt_JR6Ua3_LdNgVEI7nLwvtBqzAeL9IVgGgWDwbdXZ2wXn59WndBnkuwm03sO_leyuTAd5pINaPLLwDHI_bFj9e7AgukbYhDZLti_gJcms3euPIn_JLQDWKf3T1V06s0dQXYbWv-HvX28yGEpydAjsW41fdDZP78_OE7N9nmYIQeg6uLmhXbazjUfC2k13G2RpssGJzzFyXbdiHRuu5E6enQ
handson-gateway-service | res.Jwks: {"keys":[{"e":"AQAB","kid":"aa7c6287-c45d-4966-84b4-a1633e4e3a64","kty":"RSA","n":"tWT6dFoWKoWI1dSO6FN6FDFMFHAgvR4j6QEuVE-3Q2_iAnjXE3jCawBqYT3vCH6azWksLhsT12WjaNp3OjLP2yuDke2aVzD9c9g274XWvcJ90YI31GHdWbuZ1MgA9_Gmkq_v-cvSrkRctZAs_ktoHT4Fnpzfyl4fphCsYptL2cfe_PyJi2g6PIyrd21uhiBaHkWd8m3MBBl2yUHMqOlUJQ9NtlomAoKmqNIeGo7sfNPUtXlQ_cjQzP1aIotByu5N_PGpzLIqzDlGonytB55fOPEuJ1BXt3RxaxU_dje4LdSVlZk7Fms3XJ0Ug1r7PQE-ds7wDRkUb0vUDP9PALWrbQ"}]}
handson-gateway-service | key.Len(): 1
handson-gateway-service | k.Algorithm():
handson-gateway-service | {"level":"info","timestamp":1635893125625.8662,"logger":"gateway.grpc.server","message":"failed to verify token: failed to find matching key for verification: invalid signature algorithm : invalid jwa.SignatureAlgorithm value","request_id":"dd1b59e3-3426-462a-8cd0-c98d558e4034"}
handson-gateway-service | {"level":"info","timestamp":1635893125625.9524,"logger":"gateway.grpc.request","message":"finished","method":"/dnakano.microservice_handson.gateway.GatewayService/CreateItem","code":"Unauthenticated","request_id":"dd1b59e3-3426-462a-8cd0-c98d558e4034"}
^CGracefully stopping... (press Ctrl+C again to force)
Stopping handson-gateway-service   ... done
Stopping handson-authority-service ... done
Stopping handson-item-service      ... done
Stopping handson-catalog-service   ... done
Stopping handson-customer-service  ... done
Stopping handson-db-service        ... done
```
