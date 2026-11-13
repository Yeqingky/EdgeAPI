package services

import (
	"testing"

	rpcutils "git.ovh.im/yeqing/EdgeAPI/internal/rpc/utils"
	"git.ovh.im/yeqing/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/logs"
)

func TestServerHTTPFirewallDailyStatService_ComposeServerHTTPFirewallDashboard(t *testing.T) {
	dbs.NotifyReady()

	service := new(ServerHTTPFirewallDailyStatService)
	resp, err := service.ComposeServerHTTPFirewallDashboard(rpcutils.NewMockAdminNodeContext(1), &pb.ComposeServerHTTPFirewallDashboardRequest{})
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(resp, t)
}
