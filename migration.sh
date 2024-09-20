goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/jwt?sslmode=disable" status

goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/jwt?sslmode=disable" up