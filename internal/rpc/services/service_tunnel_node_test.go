package services

import (
	"encoding/json"
	"testing"

	"git.ovh.im/yeqing/EdgeAPI/internal/db/models"
	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
)

func TestBuildTunnelNodeConfig(t *testing.T) {
	node := &models.TunnelNode{
		ServerAddr:    "192.0.2.10",
		ServerPort:    7000,
		Token:         "tunnel-token",
		ConfigVersion: 3,
	}
	frpConfig := &nodeconfigs.FRPConfig{Token: "frp-token"}

	data, err := buildTunnelNodeConfig(node, frpConfig, "https://admin.example.com/")
	if err != nil {
		t.Fatal(err)
	}
	var config struct {
		ServerAddr string `json:"serverAddr"`
		ServerPort int    `json:"serverPort"`
		Auth       struct {
			Token string `json:"token"`
		} `json:"auth"`
		Log struct {
			Report struct {
				URL   string `json:"url"`
				Token string `json:"token"`
			} `json:"report"`
		} `json:"log"`
		RemoteConfig struct {
			URL   string `json:"url"`
			Token string `json:"token"`
		} `json:"remoteConfig"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}
	if config.ServerAddr != "192.0.2.10" || config.ServerPort != 7000 {
		t.Fatalf("unexpected server: %#v", config)
	}
	if config.Auth.Token != "frp-token" {
		t.Fatalf("unexpected auth token: %s", config.Auth.Token)
	}
	if config.Log.Report.URL != "https://admin.example.com/tunnels/logs" || config.Log.Report.Token != "tunnel-token" {
		t.Fatalf("unexpected log report: %#v", config.Log.Report)
	}
	if config.RemoteConfig.URL != "https://admin.example.com/tunnels/config" || config.RemoteConfig.Token != "tunnel-token" {
		t.Fatalf("unexpected remote config: %#v", config.RemoteConfig)
	}
}
