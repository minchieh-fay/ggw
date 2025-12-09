# GGW - gRPC Gateway

一个基于清洁架构的 gRPC 网关，用于生产者和消费者之间的消息转发。

## 功能特性

- 🏗️ **清洁架构**：采用分层架构设计，职责清晰，易于维护和扩展
- 🔌 **服务注册**：支持生产者和消费者服务注册
- 📨 **消息转发**：自动转发消息到目标服务
- 🔄 **长连接管理**：基于 gRPC 双向流的 TCP 长连接
- ⚡ **轻量级**：无需配置文件，仅需一个端口参数

## 架构设计

项目采用清洁架构（Clean Architecture），主要分为以下层次：

```
ggw/
├── cmd/                    # 应用入口
│   └── server/            # 服务器启动入口
├── internal/              # 内部代码
│   ├── domain/           # 领域模型层
│   │   ├── service.go    # 服务实体
│   │   └── message.go    # 消息实体
│   ├── repository/       # 仓储层（数据访问接口）
│   │   ├── service_repository.go          # 服务仓储接口
│   │   └── memory_service_repository.go   # 内存实现
│   ├── usecase/          # 用例层（业务逻辑）
│   │   ├── register_usecase.go   # 注册用例
│   │   └── forward_usecase.go    # 转发用例
│   ├── service/          # 服务层（基础设施）
│   │   └── connection_manager.go # 连接管理器
│   └── handler/          # 处理器层（gRPC 接口）
│       └── gateway_handler.go    # 网关处理器
└── proto/                # Protobuf 定义
    └── gateway.proto     # gRPC 服务定义
```

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- Protocol Buffers 编译器（protoc）
- Go 插件：
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

### 安装依赖

```bash
go mod download
```

### 生成 Protobuf 代码

```bash
make proto
```

或者手动执行：

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/gateway.proto
```

### 运行服务器

使用默认端口（5566）：

```bash
make run
# 或
go run cmd/server/main.go
```

指定端口：

```bash
go run cmd/server/main.go -port 8080
```

### 构建二进制文件

```bash
make build
```

生成的二进制文件在 `bin/ggw`。

## 使用方式

### 1. 服务注册

客户端首先需要调用 `Register` RPC 方法注册服务：

```go
// 注册生产者服务
registerReq := &pb.RegisterRequest{
    ServiceName: "my-producer",
    ServiceType: pb.ServiceType_PRODUCER,
}

registerResp, err := client.Register(ctx, registerReq)
serviceID := registerResp.ServiceId
```

### 2. 建立流连接

注册成功后，客户端需要建立双向流连接用于消息转发。在建立连接时，需要在 metadata 中传递服务ID：

```go
// 设置 metadata
md := metadata.New(map[string]string{
    "service-id": serviceID,
})
ctx := metadata.NewOutgoingContext(context.Background(), md)

// 建立流
stream, err := client.Stream(ctx)
```

### 3. 发送消息

通过流发送消息到目标服务：

```go
msg := &pb.ForwardMessage{
    TargetService: "target-service-name",
    Payload:       []byte("your message data"),
    Metadata:      map[string]string{"key": "value"},
}

err := stream.Send(msg)
```

### 4. 接收响应

接收来自目标服务的响应：

```go
resp, err := stream.Recv()
if resp.Success {
    // 处理响应
    data := resp.Payload
}
```

## 工作流程

1. **服务注册**：客户端连接后，首先调用 `Register` 方法注册服务，提供服务名称和类型（生产者/消费者）
2. **建立连接**：注册成功后，客户端建立双向流连接，在 metadata 中传递服务ID
3. **消息转发**：客户端通过流发送消息，网关根据目标服务名称查找对应的服务连接并转发消息
4. **自动清理**：当连接断开时，网关自动注销服务注册和连接

## 项目结构说明

### Domain 层（领域模型）

定义核心业务实体：
- `Service`：服务注册信息
- `ForwardMessage`：转发消息
- `ForwardResponse`：转发响应

### Repository 层（仓储层）

定义数据访问接口，当前提供内存实现：
- `ServiceRepository`：服务注册信息存储接口
- `MemoryServiceRepository`：内存实现（生产环境可替换为数据库实现）

### Usecase 层（用例层）

实现核心业务逻辑：
- `RegisterUsecase`：处理服务注册和注销
- `ForwardUsecase`：处理消息转发逻辑

### Service 层（服务层）

提供基础设施服务：
- `ConnectionManager`：管理所有客户端连接

### Handler 层（处理器层）

实现 gRPC 接口：
- `GatewayHandler`：处理 gRPC 请求和流

## 开发计划

- [ ] 支持负载均衡策略（轮询、随机等）
- [ ] 支持消息路由规则
- [ ] 支持服务健康检查
- [ ] 支持配置文件和更多配置选项
- [ ] 添加监控和日志
- [ ] 支持持久化存储（数据库）

## License

详见 [LICENSE](LICENSE) 文件。
