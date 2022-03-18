postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:alpine

redis:
	docker run --name redis -d -p 6379:6379 redis:alpine redis-server --save 60 1

createdb:
	docker exec -it postgres createdb --username=root --owner=root valkyrie

dropdb:
	docker exec -it postgres dropdb valkyrie

recreate:
	make dropdb && make createdb

start:
	docker start postgres && docker start redis

test:
	go test -v -cover ./service/... ./handler/...

e2e:
	go test github.com/sentrionic/valkyrie

lint:
	golangci-lint run

mock:
	mockery --all

build:
	go build github.com/sentrionic/valkyrie

fmt:
	go fmt github.com/sentrionic/...

swag:
	swag init

workflow:
	make fmt && make lint && make test