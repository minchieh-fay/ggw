package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"ggw/internal/handler"
	"ggw/internal/repository"
	"ggw/internal/service"
	"ggw/internal/usecase"
	pb "ggw/proto"
)

const (
	defaultPort = "5566"
)

func main() {
	// 解析命令行参数
	port := flag.String("port", defaultPort, "gRPC server port")
	flag.Parse()

	// 创建监听器
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 初始化依赖
	serviceRepo := repository.NewMemoryServiceRepository()
	connManager := service.NewConnectionManager()

	registerUsecase := usecase.NewRegisterUsecase(serviceRepo)
	forwardUsecase := usecase.NewForwardUsecase(serviceRepo, connManager)

	gatewayHandler := handler.NewGatewayHandler(registerUsecase, forwardUsecase, connManager)

	// 注册服务
	pb.RegisterGatewayServiceServer(grpcServer, gatewayHandler)

	log.Printf("gRPC Gateway Server starting on port %s", *port)

	// 启动服务器
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
