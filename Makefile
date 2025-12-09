.PHONY: proto build run clean

# 生成 protobuf 代码
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/gateway.proto

# 构建
build:
	go build -o bin/ggw cmd/server/main.go

# 运行
run:
	go run cmd/server/main.go

# 清理
clean:
	rm -rf bin/
	rm -f proto/*.pb.go proto/*_grpc.pb.go

