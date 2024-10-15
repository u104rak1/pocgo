# pocgo

アトラス・mockgen

POST /signup
POST /signin
POST /signout

ユーザーリソース
List User
Get User
PATCH User
DELETE User

口座リソース
GET
POST
PATCH /api/v1/me/accounts/{account_id}
DELETE

取引リソース
POST /api/v1/me/accounts/{account_id}/transactions
{
  "operation": "deposit",
  "amount": 1000
}
or
{
  "operation": "withdraw",
  "amount": 500
}
or
{
  "operation": "transfer",
  "amount": 1000
  "target_account_id": "123456789",
}
取引履歴取得
GET /api/v1/users/{user_id}/accounts/{account_id}/transactions

エンティティ
User

境界値テスト
型違いのテスト <= インテグレーションテストではいるかも


やりたいこと
slog
swagger
apidog

// @Security BearerAuth
をつけるとアクセストークンを必要になるみたい。
