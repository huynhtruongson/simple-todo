postgres-container:
	docker run --name postgresC -p 5432:5432 -e POSTGRES_PASSWORD=admin -d postgres:16-alpine

create-db:
	docker exec -it postgresC createdb --username=postgres --owner=postgres simple_todo

drop-db:
	docker exec -it postgresC dropdb simple_todo

migrate-up:
	migrate -path db/migration -database "postgresql://postgres:admin@localhost:5432/simple_todo?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://postgres:admin@localhost:5432/simple_todo?sslmode=disable" -verbose down

unit-test:
	go test -v -cover ./...