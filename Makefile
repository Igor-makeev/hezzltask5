#example: make migrate type=up
migrate:
	migrate -path internal/migrations/postgres -database "postgres://postgres:asdasd@localhost:5436/postgres?sslmode=disable" -verbose $(type)

clickhouse_migrate:
	migrate -path internal/migrations/clickhouse -database "clickhouse://localhost:9000?debug=true" -verbose $(type)

