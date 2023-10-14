# backend

## golang-migrate

```shell
migrate create -ext sql -dir db/migrations -seq init
migrate -database postgres://postgres:postgres@db:5432/app_database?sslmode=disable -path db/migrations up
migrate -database postgres://postgres:postgres@db:5432/app_database?sslmode=disable -path db/migrations down
# 全削除
migrate -database postgres://postgres:postgres@db:5432/app_database?sslmode=disable -path db/migrations drop
```

## sqlboiler

```shell
# モデルの生成
sqlboiler psql
# 生成されたテストを実行
go test ./models
# seedデータ作成用コード生成
boilingseed psql --sqlboiler-models github.com/msasaki666/backend/models
```
