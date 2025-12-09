package domain

import "time"

// ServiceType 服务类型
type ServiceType int

const (
	ServiceTypeUnknown ServiceType = iota
	ServiceTypeProducer
	ServiceTypeConsumer
)

// String 返回服务类型的字符串表示
func (st ServiceType) String() string {
	switch st {
	case ServiceTypeProducer:
		return "producer"
	case ServiceTypeConsumer:
		return "consumer"
	default:
		return "unknown"
	}
}

// Service 服务注册信息
type Service struct {
	ID           string      // 服务唯一ID
	ServiceName  string      // 服务名称
	Type         ServiceType // 服务类型
	RegisteredAt time.Time   // 注册时间
	LastActiveAt time.Time   // 最后活跃时间
}

// IsProducer 判断是否为生产者
func (s *Service) IsProducer() bool {
	return s.Type == ServiceTypeProducer
}

// IsConsumer 判断是否为消费者
func (s *Service) IsConsumer() bool {
	return s.Type == ServiceTypeConsumer
}
