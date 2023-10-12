# backend

## golang-migrate

```shell
migrate create -ext sql -dir db/migrations -seq init
migrate -database postgres://postgres:postgres@db:5432/app_database?sslmode=disable -path db/migrations up
migrate -database postgres://postgres:postgres@db:5432/app_database?sslmode=disable -path db/migrations down
```

## sqlboiler

```shell
# モデルの生成
sqlboiler psql
# 生成されたテストを実行
go test ./models
```
