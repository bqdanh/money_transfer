GO_TOOLS = 	google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/gogo/protobuf/protoc-gen-gofast \
			google.golang.org/protobuf/cmd/protoc-gen-go \
			github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
			github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
			github.com/envoyproxy/protoc-gen-validate \
			github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking \
			github.com/bufbuild/buf/cmd/protoc-gen-buf-lint \
			github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc \
			github.com/bold-commerce/protoc-gen-struct-transformer \
			github.com/googleapis/api-linter/cmd/api-linter \
			github.com/sqlc-dev/sqlc/cmd/sqlc \
			github.com/go-swagger/go-swagger/cmd/swagger \
            github.com/bufbuild/buf/cmd/buf \
            github.com/golang/mock \

install:
	go get $(GO_TOOLS)
	brew install sqlc

dependencies:
	go mod tidy
	go mod download

run-server:
	go run main.go server -c ./cmd/start_server/config/local.yaml

docker-compose-up:
	docker-compose up --build


generate_sqlc:## generate database model
	sqlc generate -f ./internal/adapters/repository/sqlc/sqlc.yaml

make proto:## generate proto
	buf generate

make sqlc:## generate sqlc querier
	sqlc generate -f ./internal/adapters/repository/sqlc/sqlc.yaml
