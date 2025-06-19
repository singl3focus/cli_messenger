LOCAL_BIN:=$(CURDIR)/bin

# for windows

install-deps:
	@set "GOBIN=$(LOCAL_BIN)" && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.1
	@set "GOBIN=$(LOCAL_BIN)" && go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	@set "GOBIN=$(LOCAL_BIN)" && go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

get-deps:
	go get -u google.golang.org/grpc
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-auth-api

generate-auth-api:
	if not exist "pkg\auth_v1" md "pkg\auth_v1"
	protoc --proto_path api\auth_v1 \
		--go_out=pkg\auth_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg\auth_v1 --go-grpc_opt=paths=source_relative \
	api\auth_v1\auth.proto