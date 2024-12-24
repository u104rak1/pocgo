package main

import (
	_ "github.com/u104rak1/pocgo/docs"
	"github.com/u104rak1/pocgo/internal/server"
)

// @title pocgo
// @version 1.0
// @description pocgoはGo * Clean Architectureで実装した簡易的な銀行操作を模したAPI Serverです。<br />詳細は<a href="https://github.com/u104rak1/pocgo">リポジトリ</a>をご覧ください。<br />*アクセストークンは`Bearer <token>`形式で入力してください。
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	server.Start()
}
