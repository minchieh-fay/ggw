package repository

import "ggw/internal/domain"

// ServiceRepository 服务注册仓储接口
type ServiceRepository interface {
	// Register 注册服务
	Register(service *domain.Service) error
	
	// Unregister 注销服务
	Unregister(serviceID string) error
	
	// GetByID 根据ID获取服务
	GetByID(serviceID string) (*domain.Service, error)
	
	// GetByName 根据名称获取服务
	GetByName(serviceName string) ([]*domain.Service, error)
	
	// GetByType 根据类型获取服务列表
	GetByType(serviceType domain.ServiceType) ([]*domain.Service, error)
	
	// UpdateLastActive 更新最后活跃时间
	UpdateLastActive(serviceID string) error
	
	// ListAll 列出所有服务
	ListAll() ([]*domain.Service, error)
}

