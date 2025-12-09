package usecase

import (
	"time"

	"github.com/google/uuid"
	"ggw/internal/domain"
	"ggw/internal/repository"
)

// RegisterUsecase 注册用例
type RegisterUsecase struct {
	serviceRepo repository.ServiceRepository
}

// NewRegisterUsecase 创建注册用例
func NewRegisterUsecase(serviceRepo repository.ServiceRepository) *RegisterUsecase {
	return &RegisterUsecase{
		serviceRepo: serviceRepo,
	}
}

// RegisterService 注册服务
func (u *RegisterUsecase) RegisterService(serviceName string, serviceType domain.ServiceType) (*domain.Service, error) {
	// 生成服务ID
	serviceID := uuid.New().String()
	
	now := time.Now()
	service := &domain.Service{
		ID:           serviceID,
		ServiceName:  serviceName,
		Type:         serviceType,
		RegisteredAt: now,
		LastActiveAt: now,
	}
	
	// 保存到仓储
	if err := u.serviceRepo.Register(service); err != nil {
		return nil, err
	}
	
	return service, nil
}

// UnregisterService 注销服务
func (u *RegisterUsecase) UnregisterService(serviceID string) error {
	return u.serviceRepo.Unregister(serviceID)
}

