package models

// TunnelLog 隧道节点日志
type TunnelLog struct {
	Id           uint64 `field:"id"`           // ID
	TunnelNodeId uint32 `field:"tunnelNodeId"` // 隧道节点ID
	Level        string `field:"level"`        // 级别
	Message      string `field:"message"`      // 内容
	CreatedAt    uint64 `field:"createdAt"`    // 创建时间
}

type TunnelLogOperator struct {
	Id           any // ID
	TunnelNodeId any // 隧道节点ID
	Level        any // 级别
	Message      any // 内容
	CreatedAt    any // 创建时间
}

func NewTunnelLogOperator() *TunnelLogOperator {
	return &TunnelLogOperator{}
}
