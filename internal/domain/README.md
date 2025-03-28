# ドメインレイヤー

## 概要
ドメインレイヤーは、ドメインモデルの知識をコードとして表現したもので、ビジネスロジックの核心部分が配置されます。このレイヤーには以下のコンポーネントが含まれます:
- ドメインエンティティ
- ドメインサービス
- 永続化用のリポジトリインターフェイス

他のレイヤーはドメインレイヤーに依存可能ですが、ドメインレイヤーは他のレイヤーに依存してはいけません。使用可能なライブラリとして以下を定義していますが、使用する場合は依存関係等をよく検討するようにして下さい。
- Goの標準ライブラリ
- Goの準標準ライブラリ
- pkgディレクトリに配置されたアプリケーション独自のユーティリティ関数

*注意
プレゼンテーションレイヤーやインフラストラクチャーレイヤーに強く関連する型や処理は、このレイヤーに含めてはいけません。
例: echo.Context 型、トランザクション、SQLエラーなど

## 集約
ドメイン内部のディレクトリを集約単位で区切ります。pocgoでは以下のリストのような集約単位で分けて管理しています。
- account
- authentication
- transaction
- user

### 集約の具体例
例えば、account集約では以下のような要素が含まれます：
- Account エンティティ（集約ルート）
- IAccountRepository インターフェース
- AccountService ドメインサービス

これらの要素は、「口座」というビジネスコンセプトを中心に凝集されています。

### 集約で考慮する点
- 集約の設計は、ビジネスルールや整合性に基づいて決定します。
- 場合によっては、RESTfulリソースに合わせて分けることも有効です
    - これはAPI設計とモデル設計の整合性を保ち、理解を容易にするためです。
    - ただしRESTfulリソースはデータ指向、DDDの集約はルール指向という違いがあります。この違いを理解した上で設計方針を選択する必要があります。
- 集約の値などを他の集約で使用してはいけないルールはありません。
    - 例えば、transactionサービスでは引き落としなどの処理の際にaccountエンティティの保存ロジックが含まれています。エンティティが使いやすいように集約を設計してください。

### 集約間の参照
集約間の参照に関するベストプラクティスは以下の通りです。

1. 参照の基本ルール
    - 集約間の参照は、IDによる参照を基本とします
    - 直接的なオブジェクト参照は避け、必要な場合はリポジトリを介して取得します
2. 参照整合性の確保
    - 集約間の参照整合性は、アプリケーションレイヤーまたはドメインサービスで確保します
    - 参照先の存在チェックは、操作の実行前に必ず行います
3. 循環参照の防止
    - 集約間の循環参照は避けます
    - 必要な場合は、一方向の参照に変更するか、新しい集約の作成を検討します

## 値オブジェクト
値オブジェクトは、システム固有の値（例: Money）やIDを表現し、ドメインモデルの一部として利用します。value_object ディレクトリに配置します。

### 特性
1. 不変性: 値オブジェクトの状態は作成後に変更されません。
2. 交換可能性: 同じ値を持つ場合、異なるインスタンスであっても等価とみなされます。
3. 等価性による比較: 値オブジェクトは等価性によって比較され、IDの概念を持ちません（エンティティとの違い）。

### ID型の実装
システムで使用される各種IDは、型安全の観点から値オブジェクトとして実装します。重複を避けるためにジェネリクスを活用したID型を定義します。
```go
// 型パラメータTによって異なるID型を表現
type ID[T any] struct {
    value string
}

// 具体的なID型の定義
type UserID = ID[userIDType]
type AccountID = ID[accountIDType]
```

## 各ファイルの責務
### {domainName}.go
このファイルでは、ドメインモデルにおけるエンティティを定義します。エンティティはライフサイクルを持ちidentityにより識別されます。またエンティティには以下の要素が含まれます:

1. プロパティ
    - エンティティが持つデータを定義します。原則として直接公開せず、外部からアクセスする場合はゲッターメソッドを使用します。
    - 例: id, userID, balance など。

2. コンストラクタ関数:
    - エンティティの生成方法を定義します。ビジネスルールを考慮し、エンティティの一貫性を確保します。
    - New 関数: 新規作成時に使用される。IDや値オブジェクトもこの関数で生成します。
    - Reconstruct 関数: 主にインフラストラクチャレイヤーから復元する際に使用されます。
    - newHoge 関数: New関数とReconstruct関数の共通部分を実装した関数です。

3. 振る舞い (メソッド):
    - エンティティが持つビジネスロジックや振る舞いを定義します。エンティティの状態を変更する操作や、データを計算・比較するロジックを含みます。
    - 例: Withdrawal, ChangeName, ComparePassword など。

### {domainName}_specification.go
このファイルではその名の通り仕様を定義します。エンティティの生成や操作時に必要なドメイン固有のバリデーションルールや定数、エラー定義を管理します。主な役割は次の通りです。

1. 定数の管理:
    - ドメインルールに基づくシステム固有の値やパラメータ値（例: 長さの上限/下限、最大許容数）を定数として定義します。

2. エラーの定義:
    - ドメインで発生しうる例外的な状況を表すエラーを定義します。エラーは一貫性を保つため、適切なメッセージを付与します。

3. バリデーション関数:
    - ドメイン固有のルールに基づき、値の妥当性を検証する関数を提供します。エンティティのコンストラクタやメソッドから再利用されます。

### {domainName}_service.go
このファイルは、ドメインサービスを実装します。ドメインサービスは、エンティティや値オブジェクトが単独で処理できないビジネスロジックを提供します。特に、集約間やリポジトリを利用する操作を実現する際に重要な役割を担います。メソッドの第一引数には`context.Context型`をとります。主な役割は以下の通りです。

1. ビジネスルールの実装:
    - エンティティだけでは完結しない、より高レベルなビジネスロジックを実装します。
2. リポジトリ操作の統合:
    - リポジトリを利用してデータを取得・操作し、ビジネスルールに従った処理を実現します。

#### Tips
ドメインサービスに定義すべきか、ユースケースに定義すべきかは議論の余地がある議題です。pocgoでは適切なビジネスロジックに区切ってドメインサービスに定義することを推奨しています。
```go
user, err := u.userRepo.FindByID(ctx, cmd.ID)
if err != nil {
    return nil, err
}
if user == nil {
    return nil, userDomain.ErrNotFound
}
```
例えば上記のようなユーザー取得ロジックがユースケースに書かれていてもルール違反ではありません。しかしこのロジックをドメインサービスのメソッドにする事で以下のメリットを享受できます。
- ユーザー取得ロジックの再利用性が高まる
- ドメインレイヤーのエラー定義がユースケースに漏れず、ユースケースのコードがシンプルになる

### {domainName}_repository.go
このファイルでは、ドメインレイヤーからデータ永続化の詳細を隠蔽するためのリポジトリのインターフェースを定義します。リポジトリは、データベースやその他のストレージ操作を抽象化し、ドメインモデルを扱う形で統一的にデータを管理します。引数が多いかつ今後の拡張が見込まれる場合、以下の例のようにparamsとして定義することを推奨します。
```go
type ListTransactionsParams struct {
    AccountID      string
    From           *time.Time
    To             *time.Time
    OperationTypes []string
    Sort           *string
    Limit          *int
    Page           *int
}
type ITransactionRepository interface {
	ListWithTotalByAccountID(ctx context.Context, params ListTransactionsParams) (transactions []*Transaction, total int, err error)
}
```

## コード規約
これまでの説明で記述できなかったコード規約を以下に記述します。

1. importしたドメインレイヤーのパッケージ名などが変数名と重複しないようにします。pocgoではドメインや値オブジェクトのimport時には`hogeDomain`, `hogeVO`という命名規則を使用しています。pkgディレクトリのパッケージ名は変数として使いそうな場合のみ`hogeUtil`という命名規則を使用しています。
    ```go
    import (
        accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
        moneyVO "github.com/u104rak1/pocgo/internal/domain/value_object/money"
        passwordUtil "github.com/u104rak1/pocgo/pkg/password"
    )
    ```
2. Exportされた関数やメソッドに外から名前を見ただけでは分かりづらいルールがある場合は、docを記述する事を推奨します。これにより使用側でルールを理解しやすくなります。
    ```go
    // ユーザーの口座を取得する。ユーザーIDとパスワードの確認はオプションであり、必要ない場合はnilを渡す。
	GetAndAuthorize(ctx context.Context, accountID AccountID, userID *userDomain.UserID, password *string) (*Account, error)
    ```
3. コンストラクタ関数でエンティティを生成する際、IDや値オブジェクトはコンストラクタ関数内で生成します。
    ```go
    func New(name, email string) (*User, error) {
	    id := UserID(ulid.New())
        return newUser(id, name, email)
    }
    ```
4. interfaceの命名は接頭辞に`I`をつけます。
    ```go
    type IAccountService interface {}
    ```
5. interfaceのモックはmockgenを使って自動生成します。makeコマンドに定義をしているので詳しくはそちらを参照して下さい。
    ```shell
    make mockgen path=internal/domain/user/user_repository.go
    ```

## ユニットテストのルール
### {domainName}_test.go
コンストラクタ関数とメソッドをメインにテストします。

#### コンストラクタ関数のテスト
1. 正常系としてビジネスルールに基づいたエンティティが正しく作成されることをテストします。生成されたエンティティのプロパティをゲッター関数を使用して期待通りであることを確認します。ゲッター関数で値を取得できる事を確認できるのでゲッター関数個別のテストは不要とします。IDの確認はカスタム型のIDを検証します。
2. 異常系としてバリデーションの包括的なテストを行います。各プロパティ毎に無効な値や境界値をセットして異常系のテスト（境界値の正常系のテスト）を行います。ここでバリデーションのテストを行う為、{domainName}_specification.go に定義されているバリデーション関数を個別にテストする必要はありません。
3. バリデーションテストにおける型の違いをチェックするテストは不要とします。言語の特性によりドメインレイヤーに異常な型が渡されることはあり得ない為、型に関しては常に正確と判断して構いません。
4. 並列実行の為の`*testing.T.Parallel()`を使用することを推奨します。テスト速度の向上とテストの独立性を確保します。これを使うことでエラーが発生するという事は、作成されたエンティティが独立していないなど、テストの書き方に問題があると判断できます。
5. Reconstruct関数のテストでは、正常系のみをテストします。なぜなら異常系においてはNew関数のテストと重複する為です。テストの際は`IDString()`ゲッターを使用してString型のIDを検証します。
6. カバレッジは可能な限りすべての分岐に対するテストを作成し、100%を目指します。ただし、必ず100%にする必要はありません。エラーハンドリングでは殆ど発生しないパターンであっても明示することが推奨される為、再現が困難な分岐が発生します。そのような発生頻度が極めて低い分岐については無理にテストを書かずに省略して構いません。（スルーする理由をコメントで明示することが推奨されます）
7. 定数やエラーの扱いについて、ドメインレイヤーのテストでは、{domainName}_specification.go に定義されている定数やエラーを直接使用せず、`期待される値を明示的に記述`します。これはビジネスルールが意図せず変更された場合にテストが失敗することで検知可能にする為です。意図的な変更の場合はテストの修正も含める事で安全性を確保します。
    - 補足
    - `他の集約で定義されているエラーや定数`は使用するべきです。なぜならその値を管理するのは定義された集約の責務だからです。値のテストを重複管理するのは極力避けるのが望ましいです。
    - しかし、そうは言っても人の目で管理する必要があるので、正確に掌握するのは難しいです。（他のレイヤーのテストであれば定数が使われていない事に違和感を持てますが、ドメインレイヤーのテストでは文脈を読まないと定数がその集約で定義されているものかどうかを瞬時に判断することが難しいです。）
    - その為、他の集約で定義されたエラーや定数であっても`期待される値を明示的に記述`する事を許容します。
    - 最も重要な事は、ビジネスルールが意図せず変更された場合にテストが失敗することで検知可能にする事です。

#### メソッドのテスト
1. メソッドがビジネスロジックに従い、正しく動作することを確認します。特に、値の変更後のエンティティの状態確認ややエラーが発生した際に値が変更されていない事・エラーメッセージの検証に重点を置きます。
2. バリデーションテストの網羅的な確認（無効な値や境界値のテスト）はコンストラクタ関数のテストで行うため、メソッドのテストでは正常系と分岐による異常系のテストに集中します。

### {domainName}_service_test.go
1. 正常系・異常系を網羅的に検証します。リポジトリや外部依存をモック化し、サービス内のすべての分岐を確認します。ただし、コンストラクタ関数のテストのルールで言及したようにカバレッジ100%は努力目標とします。
2. コンストラクタ関数のテストで言及した通り、定数やエラーに関しては`期待される値を明示的に記述`します。
3. モックに渡す引数は議論の余地がありますが、pocgoでは検証しません。理由は以下の通りです。pocgoでは`gomock.Any()`を使用して引数の検証を実質的にスキップしています。
    - 引数の検証はモックされた関数自体のテストで行うべきな為
    - テストは対象のInputによってOutputが期待通りかを検証するべきであり、実装の中身は考慮しない事が望ましい為。
