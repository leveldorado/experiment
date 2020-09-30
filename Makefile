# generate - generate code from proto file
# requirements protoc libprotoc 3.13.0 and installed protoc-gen-go
generate:
	protoc --proto_path=protobuf --go_out=plugins=grpc:. protobuf/ports.proto

test: export TEST_MONGODB_URL=mongodb://localhost:27017

test:
	docker rm -f mongo-test  2>/dev/null; true
	docker run --name mongo-test -p 27017:27017 -d  mongo
	go test ./...
	docker rm -f mongo-test

build:
	cd api && go build
	cd ports && go build

docker-build:
	docker build -t experiment-api -f api.dockerfile .
	docker build -t experiment-ports -f ports.dockerfile .

up:
	docker-compose up -d

down:
	docker-compose down