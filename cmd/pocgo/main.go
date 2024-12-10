package main

import (
	_ "github.com/u104rak1/pocgo/docs"
	"github.com/u104rak1/pocgo/internal/server"
)

// @title pocgo API
// @version 1.0
// @description This is a sample server. <br />Please enter your token in the format: "Bearer <token>" in the Authorization header.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	server.Start()
}
