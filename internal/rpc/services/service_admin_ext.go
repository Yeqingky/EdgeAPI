//go:build !plus

package services

import (
	"context"

	"git.ovh.im/yeqing/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/dbs"
)

// ComposeAdminDashboard方法扩展
func (this *AdminService) composeAdminDashboardExt(tx *dbs.Tx, ctx context.Context, result *pb.ComposeAdminDashboardResponse) error {
	return nil
}
