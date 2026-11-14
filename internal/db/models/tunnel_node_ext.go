package models

import (
	"encoding/json"
	"strings"

	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
)

// BuildFRPCConfig builds the frpc configuration delivered to a tunnel node.
func (this *TunnelNode) BuildFRPCConfig(frpConfig *nodeconfigs.FRPConfig, baseURL string) ([]byte, error) {
	baseURL = strings.TrimRight(baseURL, "/")
	config := map[string]any{
		"serverAddr": this.ServerAddr,
		"serverPort": this.ServerPort,
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
	}
	return json.Marshal(config)
}
