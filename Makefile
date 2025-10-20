.PHONY: postgres adminer migrate

# Start Postgres container on port 5432
postgres:
	docker run --rm -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres

# Start Adminer on port 8080 host : host.docker.internal
adminer:
	docker run --rm -p 8080:8080 adminer

# Run migrations against Postgres    #     make migrate (cmd)
migrate:
	migrate -source file://migrations -database "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable" up

migrate-down:
	migrate -source file://migrations -database "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable" down

#compile  CompileDaemon -command="go run ./cmd/goreddit"

#compile CompileDaemon -command="go run ."