package models

// TunnelNode 隧道节点 (frpc)
type TunnelNode struct {
	Id            uint32 `field:"id"`            // ID
	AdminId       uint32 `field:"adminId"`       // 管理员ID
	UserId        uint32 `field:"userId"`        // 用户ID
	ClusterId     uint32 `field:"clusterId"`     // 集群ID
	Name          string `field:"name"`          // 名称
	Token         string `field:"token"`         // 令牌
	ServerAddr    string `field:"serverAddr"`    // frps 地址
	ServerPort    uint16 `field:"serverPort"`    // frps 端口
	OS            string `field:"os"`            // 操作系统
	Arch          string `field:"arch"`          // CPU架构
	Version       string `field:"version"`       // 版本
	IsOn          bool   `field:"isOn"`          // 是否启用
	IsInstalled   bool   `field:"isInstalled"`   // 是否已安装
	ConfigVersion uint64 `field:"configVersion"` // 配置版本
	CreatedAt     uint64 `field:"createdAt"`     // 创建时间
	UpdatedAt     uint64 `field:"updatedAt"`     // 修改时间
	State         uint8  `field:"state"`         // 状态
}

type TunnelNodeOperator struct {
	Id            any // ID
	AdminId       any // 管理员ID
	UserId        any // 用户ID
	ClusterId     any // 集群ID
	Name          any // 名称
	Token         any // 令牌
	ServerAddr    any // frps 地址
	ServerPort    any // frps 端口
	OS            any // 操作系统
	Arch          any // CPU架构
	Version       any // 版本
	IsOn          any // 是否启用
	IsInstalled   any // 是否已安装
	ConfigVersion any // 配置版本
	CreatedAt     any // 创建时间
	UpdatedAt     any // 修改时间
	State         any // 状态
}

func NewTunnelNodeOperator() *TunnelNodeOperator {
	return &TunnelNodeOperator{}
}
