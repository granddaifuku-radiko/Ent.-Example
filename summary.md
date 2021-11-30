# 基本

## セットアップ
1. `go get -d entgo.io/ent/cmd/ent`
2. `go run entgo.io/ent/cmd/ent init ${スキーマ名}` でプロジェクトのスキーマ生成
3. `go generete ./ent`で`ent配下`にDBと通信するためのコードが生成

## スキーマ

### フィールド
`./ent/schema/${スキーマ名}.go`に追記
`func (${スキーマ名}) Fields() []ent.Field {}`に `[]ent.Field`の形で追記していく
- Text
- String
- Time
- Enum
- Int
- Default 
- Value
- NotEmpty
- Immutable
等々サポートされている

- デフォルトで`ID`フィールドがある

- バリデーション（独自・組み込み）

- `go generete ./...`でコード生成

```
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			Annotations(
				entproto.Field(2),
			),
		field.String("email_address").
			Unique().
			Annotations(
				entproto.Field(3),
			),
		field.String("alias").
			Optional().
			Annotations(
				entproto.Field(4),
			),
	}
}
```

### インデックス
`./ent/schema/${スキーマ名}.go`に新たに`[]ent.Index`を返す関数を作成し追記

### アノテーション
- カスタムテーブル名
- 外部キー制約

## クエリ
- `Query`にメソッドチェーンしていくことで`Where`や`Only`, `Group By`などによる条件指定ができる
`client.Pet.Query().Where(pet.Not(pet.NameHasPrefix("Ari"))).All(ctx)`

## マイグレーション
- 初期化時に自動マイグレーション

## エッジ（エンティティ同士の関係）
`./ent/scheema/${スキーマ名}.go`に追記
`func (${スキーマ名}) Edges() []ent.Edge {}`に `[]ent.Edge`の形で追記していく

```
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("administered", Category.Type).
			Ref("admin").
			Annotations(entproto.Field(5)),
	}
}
```

## `sql.DB`
- `ent client`に`sql.DB`ドライバを渡すことができる

## テスト
- デフォルトで`enttest`パッケージが生成
- クライアント生成時のオプションとして渡すことができる
```Go
func TestXXX(t *testing.T) {
    opts := []enttest.Option{
        enttest.WithOptions(ent.Log(t.Log)),
        enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
    }
    client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", opts...)
    defer client.Close()
    // ...
}
```



# gRPC
- `entproto`コマンドでスキーマからProtocol Buffer定義とgRPCサービス定義を生成
- `protoc-gen-entgrpc`
- gRPCサーバは自身で実装する必要がある

## セットアップ
1. 基本的なセットアップ
2. `go get -u entgo.io/contrib/entproto`で`entproto`パッケージを追加

## Protobufの生成
- entスキーマからProtobufスキーマを生成するためにアノテーションを付与する必要がある
1. スキーマにAnnotationを返す関数を追記
2. スキーマのフィールドに一意の番号を割り当てる（entは自動でIDフィールドを生成するため、番号は2から始めるべき）
3. `./ent/generate.go`に`entproto`を呼び出すための`//go:generate go run -mod=mod entgo.io/contrib/entproto/cmd/entproto -path ./schema`を追記
4. `go generate ./...`で`./ent/proto/`配下にprotobuf関連のディレクトリが生成
- ここで生成される`./ent/proto/entpb/generate.go`は`.proto`からGoを生成するコードジェネレータ`protoc`を呼び出す
- protoc
- protoc-gen-go
- protoc-gen-go-grpc
- protoc-gen-entgrpc
上記4つをインストールする必要あり

5. `go generate ./...`で`./ent/proto/entpb/entpb.pb.go`が生成

## gRPCサービスの生成
- スキーマの`Annotation`に`entproto.Service()`と追記&`go generate ./...`でCRUDサービスが生成
