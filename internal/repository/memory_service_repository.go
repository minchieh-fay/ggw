package repository

import (
	"sync"
	"time"

	"ggw/internal/domain"
)

// MemoryServiceRepository 内存实现的服务仓储
type MemoryServiceRepository struct {
	services map[string]*domain.Service
	mu       sync.RWMutex
}

// NewMemoryServiceRepository 创建内存服务仓储
func NewMemoryServiceRepository() *MemoryServiceRepository {
	return &MemoryServiceRepository{
		services: make(map[string]*domain.Service),
	}
}

// Register 注册服务
func (r *MemoryServiceRepository) Register(service *domain.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.services[service.ID] = service
	return nil
}

// Unregister 注销服务
func (r *MemoryServiceRepository) Unregister(serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	delete(r.services, serviceID)
	return nil
}

// GetByID 根据ID获取服务
func (r *MemoryServiceRepository) GetByID(serviceID string) (*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	service, exists := r.services[serviceID]
	if !exists {
		return nil, ErrServiceNotFound
	}
	return service, nil
}

// GetByName 根据名称获取服务
func (r *MemoryServiceRepository) GetByName(serviceName string) ([]*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*domain.Service
	for _, service := range r.services {
		if service.ServiceName == serviceName {
			result = append(result, service)
		}
	}
	return result, nil
}

// GetByType 根据类型获取服务列表
func (r *MemoryServiceRepository) GetByType(serviceType domain.ServiceType) ([]*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*domain.Service
	for _, service := range r.services {
		if service.Type == serviceType {
			result = append(result, service)
		}
	}
	return result, nil
}

// UpdateLastActive 更新最后活跃时间
func (r *MemoryServiceRepository) UpdateLastActive(serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	service, exists := r.services[serviceID]
	if !exists {
		return ErrServiceNotFound
	}
	
	service.LastActiveAt = time.Now()
	return nil
}

// ListAll 列出所有服务
func (r *MemoryServiceRepository) ListAll() ([]*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]*domain.Service, 0, len(r.services))
	for _, service := range r.services {
		result = append(result, service)
	}
	return result, nil
}

