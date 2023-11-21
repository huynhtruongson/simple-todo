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

gen-mock:
	mockery --dir lib --all --output mocks/lib --with-expecter
	mockery --dir services/user --all --output mocks/user --with-expecter
	mockery --dir services/task --all --output mocks/task --with-expecter
	mockery --dir services/auth --all --output mocks/auth --with-expecter
	mockery --dir token --all --output mocks/token --with-expecter

gen-proto:
	rm -f pb/*.proto
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans:
	evans --host localhost --port 3001 -r repl

unit-test:
	go test -v -cover ./...