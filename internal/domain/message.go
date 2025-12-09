package domain

// ForwardMessage 转发消息
type ForwardMessage struct {
	TargetService string            // 目标服务名称
	Payload       []byte            // 消息负载
	Metadata      map[string]string // 元数据
	FromServiceID string            // 发送方服务ID
}

// ForwardResponse 转发响应
type ForwardResponse struct {
	Success bool   // 是否成功
	Message string // 响应消息
	Payload []byte // 响应负载
}
