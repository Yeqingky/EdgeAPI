package utils_test

import (
	"testing"

	"git.ovh.im/yeqing/EdgeAPI/internal/utils"
	"git.ovh.im/yeqing/EdgeCommon/pkg/dnsconfigs"
)

func TestLookupCNAME(t *testing.T) {
	t.Log(utils.LookupCNAME("www.yun4s.cn"))
}

func TestLookupNS(t *testing.T) {
	t.Log(utils.LookupNS("goedge.cn", nil))
}

func TestLookupNSExtra(t *testing.T) {
	t.Log(utils.LookupNS("goedge.cn", []*dnsconfigs.DNSResolver{
		{
			Host: "192.168.2.2",
		},
		{
			Host: "192.168.2.2",
			Port: 58,
		},
		{
			Host: "8.8.8.8",
			Port: 53,
		},
	}))
}

func TestLookupTXT(t *testing.T) {
	t.Log(utils.LookupTXT("yanzheng.goedge.cn", nil))
}
