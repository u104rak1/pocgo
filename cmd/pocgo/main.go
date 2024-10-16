package main

import (
	_ "github.com/ucho456job/pocgo/docs"
	"github.com/ucho456job/pocgo/internal/server"
)

// @title pocgo API
// @version 1.0
// @description This is a sample server.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	server.Start()
}
