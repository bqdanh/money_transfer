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

dependencies:
	go mod tidy
	go mod download

run-http-service:
	go run main.go http_service --config=configs/http_server/config.yaml

docker-compose-up:
	docker-compose up --build


generate_sqlc:## generate database model
	sqlc generate -f ./internal/adapters/repository/sqlc/sqlc.yaml