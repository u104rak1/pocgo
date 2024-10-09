#!/bin/bash

# Atlasがインストールされているか確認する
if ! command -v atlas &> /dev/null
then
    echo "Atlas is not installed. Please install it by running:"
    echo "curl -sSf https://atlasgo.sh | sh"
    exit 1
fi

# schema.sqlを更新する
go run ./cmd/postgres/main.go generate
if [ $? -ne 0 ]; then
    echo "Failed to connect to database. Please start the database by running:"
    echo "make dependencies_start"
    exit 1
fi

# データベースをリセットする
go run ./cmd/postgres/main.go migrate reset

# マイグレーションファイルとschema.sqlの差分を確認し、マイグレーションファイルを生成する
atlas migrate diff migration \
  --dir 'file://internal/infrastructure/postgres/migrations?format=golang-migrate' \
  --to 'file://internal/infrastructure/postgres/schema.sql' \
  --dev-url 'postgres://local_user:password@localhost:5432/POCGO_LOCAL_DB?sslmode=disable'
