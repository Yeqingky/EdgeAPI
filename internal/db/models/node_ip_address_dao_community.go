//go:build !plus
// +build !plus

package models

import (
	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
	"github.com/iwind/TeaGo/dbs"
)

// FireThresholds 触发阈值
func (this *NodeIPAddressDAO) FireThresholds(tx *dbs.Tx, role nodeconfigs.NodeRole, nodeId int64) error {
	return nil
}
