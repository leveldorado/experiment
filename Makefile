# generate - generate code from proto file
# requirements protoc libprotoc 3.13.0 and installed protoc-gen-go
generate:
	protoc --proto_path=protobuf --go_out=plugins=grpc:. protobuf/ports.proto

test:
	go test ./...