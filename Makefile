CWD = /go/src/apptica

lint:
	@docker run --rm -t -v $(CURDIR):$(CWD) -w $(CWD) golangci/golangci-lint golangci-lint run

unit:
	@docker run --rm -v $(CURDIR):$(CWD) -w $(CWD) golang:1.16 sh -c "go list ./... | xargs go test"

build: generate_pb lint unit
	@docker build -t apptica:latest .

up:
	@docker-compose up

down:
	@docker-compose down -v --remove-orphans

use_db:
	@docker-compose exec db mysql -uroot -pqwerty apptica

generate_pb:
	@docker run --rm -v $(CURDIR)/:/src -w /src/grpcserver/schema namely/protoc-all -f schema.proto -l go -o /src