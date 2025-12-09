package usecase

import (
	"fmt"

	"ggw/internal/domain"
	"ggw/internal/repository"
)

// ForwardUsecase 转发用例
type ForwardUsecase struct {
	serviceRepo repository.ServiceRepository
	connManager ConnectionManager
}

// ConnectionManager 连接管理器接口
type ConnectionManager interface {
	SendMessage(serviceID string, msg *domain.ForwardMessage) error
}

// NewForwardUsecase 创建转发用例
func NewForwardUsecase(serviceRepo repository.ServiceRepository, connManager ConnectionManager) *ForwardUsecase {
	return &ForwardUsecase{
		serviceRepo: serviceRepo,
		connManager: connManager,
	}
}

// ForwardMessage 转发消息到目标服务
func (u *ForwardUsecase) ForwardMessage(fromServiceID string, msg *domain.ForwardMessage) error {
	// 查找目标服务
	services, err := u.serviceRepo.GetByName(msg.TargetService)
	if err != nil {
		return fmt.Errorf("failed to find target service: %w", err)
	}
	
	if len(services) == 0 {
		return fmt.Errorf("target service not found: %s", msg.TargetService)
	}
	
	// 设置发送方ID
	msg.FromServiceID = fromServiceID
	
	// 转发到第一个匹配的服务（可以根据策略选择）
	// 这里简化处理，实际可以根据负载均衡策略选择
	targetService := services[0]
	
	// 通过连接管理器发送消息
	if err := u.connManager.SendMessage(targetService.ID, msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	
	return nil
}

