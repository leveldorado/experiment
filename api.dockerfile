FROM jaegertracing/protobuf as protoc

COPY protobuf/ports.proto ./protobuf/ports.proto

RUN protoc --proto_path=protobuf --go_out=plugins=grpc:. protobuf/ports.proto

FROM golang as builder

WORKDIR /app

COPY . .

COPY --from=protoc /grpc ./grpc

RUN go mod download

RUN go build api/main.go

EXPOSE 8000

ENTRYPOINT ["./main"]
