package service

import (
	"fmt"
	"sync"

	"ggw/internal/domain"
)

// Connection 连接接口
type Connection interface {
	Send(msg *domain.ForwardResponse) error
	GetID() string
}

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[string]Connection // serviceID -> Connection
	mu          sync.RWMutex
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]Connection),
	}
}

// Register 注册连接
func (cm *ConnectionManager) Register(serviceID string, conn Connection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.connections[serviceID] = conn
}

// Unregister 注销连接
func (cm *ConnectionManager) Unregister(serviceID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	delete(cm.connections, serviceID)
}

// SendMessage 发送消息到指定服务
func (cm *ConnectionManager) SendMessage(serviceID string, msg *domain.ForwardMessage) error {
	cm.mu.RLock()
	conn, exists := cm.connections[serviceID]
	cm.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("connection not found for service: %s", serviceID)
	}
	
	response := &domain.ForwardResponse{
		Success: true,
		Message: "forwarded",
		Payload: msg.Payload,
	}
	
	return conn.Send(response)
}

// GetConnection 获取连接
func (cm *ConnectionManager) GetConnection(serviceID string) (Connection, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	conn, exists := cm.connections[serviceID]
	return conn, exists
}

// ListConnections 列出所有连接的服务ID
func (cm *ConnectionManager) ListConnections() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	ids := make([]string, 0, len(cm.connections))
	for id := range cm.connections {
		ids = append(ids, id)
	}
	return ids
}

