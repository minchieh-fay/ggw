package handler

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"ggw/internal/domain"
	"ggw/internal/service"
	"ggw/internal/usecase"
	pb "ggw/proto"
)

// GatewayHandler gRPC 处理器
type GatewayHandler struct {
	pb.UnimplementedGatewayServiceServer
	
	registerUsecase *usecase.RegisterUsecase
	forwardUsecase  *usecase.ForwardUsecase
	connManager     *service.ConnectionManager
}

// NewGatewayHandler 创建网关处理器
func NewGatewayHandler(
	registerUsecase *usecase.RegisterUsecase,
	forwardUsecase *usecase.ForwardUsecase,
	connManager *service.ConnectionManager,
) *GatewayHandler {
	return &GatewayHandler{
		registerUsecase: registerUsecase,
		forwardUsecase:  forwardUsecase,
		connManager:     connManager,
	}
}

// Register 处理服务注册
func (h *GatewayHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 转换服务类型
	var serviceType domain.ServiceType
	switch req.ServiceType {
	case pb.ServiceType_PRODUCER:
		serviceType = domain.ServiceTypeProducer
	case pb.ServiceType_CONSUMER:
		serviceType = domain.ServiceTypeConsumer
	default:
		return &pb.RegisterResponse{
			Success: false,
			Message: "invalid service type",
		}, status.Error(codes.InvalidArgument, "invalid service type")
	}
	
	// 注册服务
	service, err := h.registerUsecase.RegisterService(req.ServiceName, serviceType)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}
	
	return &pb.RegisterResponse{
		Success:   true,
		Message:   "registered successfully",
		ServiceId: service.ID,
	}, nil
}

// Stream 处理双向流消息转发
func (h *GatewayHandler) Stream(stream pb.GatewayService_StreamServer) error {
	// 等待第一个消息来确定服务ID（应该在注册后建立流）
	// 这里简化处理，实际应该通过上下文传递服务ID
	
	var serviceID string
	var initialized bool
	
	for {
		// 接收消息
		msg, err := stream.Recv()
		if err == io.EOF {
			// 客户端关闭连接
			if initialized {
				h.connManager.Unregister(serviceID)
				h.registerUsecase.UnregisterService(serviceID)
			}
			return nil
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			if initialized {
				h.connManager.Unregister(serviceID)
				h.registerUsecase.UnregisterService(serviceID)
			}
			return err
		}
		
		// 如果是第一个消息，需要包含服务ID（通过metadata传递）
		if !initialized {
			// 从metadata获取服务ID（需要在客户端设置）
			md, ok := metadata.FromIncomingContext(stream.Context())
			if !ok {
				return status.Error(codes.Unauthenticated, "metadata not found")
			}
			
			serviceIDs := md.Get("service-id")
			if len(serviceIDs) == 0 || serviceIDs[0] == "" {
				return status.Error(codes.Unauthenticated, "service not registered")
			}
			serviceID = serviceIDs[0]
			
			// 注册连接
			conn := &StreamConnection{
				stream:    stream,
				serviceID: serviceID,
			}
			h.connManager.Register(serviceID, conn)
			initialized = true
			continue
		}
		
		// 转发消息
		forwardMsg := &domain.ForwardMessage{
			TargetService: msg.TargetService,
			Payload:       msg.Payload,
			Metadata:      msg.Metadata,
		}
		
		if err := h.forwardUsecase.ForwardMessage(serviceID, forwardMsg); err != nil {
			log.Printf("Error forwarding message: %v", err)
			// 发送错误响应
			response := &pb.ForwardResponse{
				Success: false,
				Message: err.Error(),
			}
			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}

// StreamConnection 流连接实现
type StreamConnection struct {
	stream    pb.GatewayService_StreamServer
	serviceID string
}

// Send 发送消息
func (sc *StreamConnection) Send(msg *domain.ForwardResponse) error {
	response := &pb.ForwardResponse{
		Success: msg.Success,
		Message: msg.Message,
		Payload: msg.Payload,
	}
	return sc.stream.Send(response)
}

// GetID 获取服务ID
func (sc *StreamConnection) GetID() string {
	return sc.serviceID
}

