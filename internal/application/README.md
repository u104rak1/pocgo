# アプリケーションレイヤー

## 概要
アプリケーションレイヤーは、ドメインレイヤーのビジネスロジックを実行するためのインターフェースを提供します。このレイヤーには以下のコンポーネントが含まれます:
- ユースケース
- unitOfWork関数（トランザクション管理）
- JWTサービスのインターフェース

アプリケーションレイヤーはドメインレイヤーに依存することができます。

## コマンド
コマンドはユースケースの入力パラメーターを表現します。
型はプリミティブ型を使用します。
```go
type ReadUserCommand struct {
	ID string
}
```

## DTO (Data Transfer Object)
DTOはユースケースの出力パラメーターを表現します。
コマンドと同様にプリミティブ型を使用します。
```go
type ReadUserDTO struct {
	ID    string
	Name  string
	Email string
}
```

## unitOfWork
unitOfWorkはトランザクション管理を行うためのインターフェースです。トランザクションで管理したい動作を`RunInTx`メソッドでラップして実行することで、インフラストラクチャレイヤーでトランザクションを管理することができます。`IUnitOfWork`は戻り値が`error`のみのメソッドで、結果が欲しい場合は`IUnitOfWorkWithResult[T any]`を使用します。
```go
type IUnitOfWork interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error
}
type IUnitOfWorkWithResult[T any] interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) (*T, error)) (*T, error)
}
```

## JWTサービス
JWTサービスはJWTの生成と検証を行うためのインターフェースです。アクセストークンの生成と検証の為に使用します。このロジック自体はドメインレイヤーの責務では無い為、アプリケーションレイヤーにインターフェイスを定義し、実際のロジックはインフラストラクチャレイヤーに実装します。
```go
type IJwtService interface {
	GenerateToken(ctx context.Context, userID string) (string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
}
```

## コード規約
1. パッケージは基本的にドメインレイヤーの集約と連動して分けますが絶対ではありません。
2. interfaceの命名は接頭辞に`I`をつけます。
    ```go
    type IAccountService interface {}
    ```
3. 一つのファイルに一つのユースケース及びコマンド、DTOを書きます。ユースケースの実行メソッドは`Run`という名前をつけます。またRunメソッドの引数は`ctx context.Context`と`cmd`です。返り値は`*DTO, error`です。
    ```go
    func (u *readUserUsecase) Run(ctx context.Context, cmd ReadUserCommand) (*ReadUserDTO, error) {
    }
    ```
4. コマンドの値がプリミティブ型である為、カスタム型に変換するロジックはユースケースレイヤーの責務となります。

## ユニットテストのルール
1. 正常系としてユースケースのテストを行います。
2. DIされたモックをテスト用に利用します。
3. 異常系として各エラーハンドリングの分岐をテストし、カバレッジ100%を目指します。エラーの内容については責務外とし、一貫して `assert.AnError`を使用します。失敗した場合、エラーが発生することのみ確認します。
4. モックに渡す引数は検証しません。`gomock.Any()`を使用して検証をスキップします。理由の詳細は[ドメインレイヤーのユニットテストのルール](../domain/README.md#ユニットテストのルール)を参照してください。
5. 並列実行の為の`*testing.T.Parallel()`を使用することを推奨します。テスト速度の向上とテストの独立性を確保します。これを使うことでエラーが発生するという事は、作成されたエンティティが独立していないなど、テストの書き方に問題があると判断できます。
