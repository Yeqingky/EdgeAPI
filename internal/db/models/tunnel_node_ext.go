package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
	"github.com/iwind/TeaGo/dbs"
)

type TunnelServerAddr struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

func FindTunnelServerAddrs(tx *dbs.Tx, clusterId int64, port int) ([]TunnelServerAddr, error) {
	addresses, err := SharedNodeIPAddressDAO.FindAllAccessibleIPAddressesWithClusterId(tx, nodeconfigs.NodeRoleNode, clusterId, nil)
	if err != nil {
		return nil, err
	}
	var result []TunnelServerAddr
	seen := map[string]bool{}
	for _, address := range addresses {
		if address.Ip == "" || seen[address.Ip] {
			continue
		}
		seen[address.Ip] = true
		result = append(result, TunnelServerAddr{Addr: address.Ip, Port: port})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Addr < result[j].Addr })
	return result, nil
}

// BuildFRPCConfig builds the frpc configuration delivered to a tunnel node.
func (this *TunnelNode) BuildFRPCConfig(frpConfig *nodeconfigs.FRPConfig, baseURL string, serverAddrs []TunnelServerAddr) ([]byte, error) {
	baseURL = strings.TrimRight(baseURL, "/")
	if len(serverAddrs) == 0 {
		return nil, errors.New("no accessible FRP server address in cluster")
	}
	serverAddr := serverAddrs[0]
	config := map[string]any{
		"serverAddr":  serverAddr.Addr,
		"serverPort":  serverAddr.Port,
		"serverAddrs": serverAddrs,
		"auth": map[string]any{
			"method": "token",
			"token":  frpConfig.Token,
		},
		"log": map[string]any{
			"to":                "console",
			"level":             "info",
			"disablePrintColor": true,
			"report": map[string]any{
				"url":   baseURL + "/tunnels/logs",
				"token": this.Token,
			},
		},
		"remoteConfig": map[string]any{
			"url":      baseURL + "/tunnels/config",
			"token":    this.Token,
			"interval": 30,
		},
		"proxies": []map[string]any{
			{
				"name":      fmt.Sprintf("tunnel-%d", this.Id),
				"type":      "stcp",
				"secretKey": this.Token,
				"plugin":    map[string]any{"type": "http_proxy"},
				"transport": map[string]any{"useEncryption": true},
			},
		},
	}
	return json.Marshal(config)
}
