package main

import (
	"context"

	_ "github.com/mattn/go-sqlite3"

	"ent-grpc-example/ent/category"
	"ent-grpc-example/ent/enttest"
	"ent-grpc-example/ent/proto/entpb"
	"ent-grpc-example/ent/user"
	"testing"
)

func TestUserProto(t *testing.T) {
	user := entpb.User{
		Name:         "granddaifuku",
		EmailAddress: "granddaifuku@example.com",
	}
	if user.GetName() != "granddaifuku" {
		t.Fatal("expected user name to be granddaifuku")
	}
	if user.GetEmailAddress() != "granddaifuku@example.com" {
		t.Fatal("expected email address to be granddaifuku@example.com")
	}
}

func TestServiceWithEdges(t *testing.T) {
	// インメモリのsqliteインスタンスに接続されたentクライアントの初期化から始めます
	ctx := context.Background()
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// 次に、Userサービスを初期化します。 ここでは、実際にポートを開いてgRPCサーバーを作成するのではなく
	// ライブラリのコードを直接呼び出していることに注目してください。
	svc := entpb.NewUserService(client)

	// 次に、entクライアントを使って直接Categoryを作成します。
	// Userとは無関係に初期化していることに注意してください。
	cat := client.Category.Create().SetName("cat_1").SaveX(ctx)

	// 次に、User サービスの `Create` メソッドを呼び出します。
	// IDのみが設定されたentpb.Categoryインスタンスのリストを渡していることに注意してください。
	create, err := svc.Create(ctx, &entpb.CreateUserRequest{
		User: &entpb.User{
			Name:         "user",
			EmailAddress: "user@service.code",
			Administered: []*entpb.Category{
				{Id: int32(cat.ID)},
			},
		},
	})
	if err != nil {
		t.Fatal("failed creating user using UserService", err)
	}

	// すべてが正しく動作したことを確認するために, カテゴリーテーブルをクエリします。
	// 作成したユーザーが管理するカテゴリーが1つだけあることを確認します。
	count, err := client.Category.
		Query().
		Where(
			category.HasAdminWith(
				user.ID(int(create.Id)),
			),
		).
		Count(ctx)
	if err != nil {
		t.Fatal("failed counting categories admin by created user", err)
	}
	if count != 1 {
		t.Fatal("expected exactly one group to managed by the created user")
	}
}
