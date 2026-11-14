package services

import (
	"encoding/json"
	"testing"

	"git.ovh.im/yeqing/EdgeAPI/internal/db/models"
	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
)

func TestBuildTunnelNodeConfig(t *testing.T) {
	node := &models.TunnelNode{
		Id:            7,
		ServerAddr:    "192.0.2.10",
		ServerPort:    7000,
		Token:         "tunnel-token",
		ConfigVersion: 3,
	}
	frpConfig := &nodeconfigs.FRPConfig{Token: "frp-token"}

	data, err := buildTunnelNodeConfig(node, frpConfig, "https://admin.example.com/", []models.TunnelServerAddr{{Addr: "192.0.2.10", Port: 7000}, {Addr: "192.0.2.11", Port: 7000}})
	if err != nil {
		t.Fatal(err)
	}
	var config struct {
		ServerAddr  string `json:"serverAddr"`
		ServerAddrs []struct {
			Addr string `json:"addr"`
			Port int    `json:"port"`
		} `json:"serverAddrs"`
		ServerPort int `json:"serverPort"`
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
		Proxies []struct {
			Name      string `json:"name"`
			Type      string `json:"type"`
			SecretKey string `json:"secretKey"`
			Plugin    struct {
				Type string `json:"type"`
			} `json:"plugin"`
			Transport struct {
				UseEncryption bool `json:"useEncryption"`
			} `json:"transport"`
		} `json:"proxies"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}
	if len(config.ServerAddrs) != 2 || config.ServerAddrs[1].Addr != "192.0.2.11" {
		t.Fatalf("unexpected server addrs: %#v", config.ServerAddrs)
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
	if len(config.Proxies) != 1 || config.Proxies[0].Name != "tunnel-7" || config.Proxies[0].Type != "stcp" ||
		config.Proxies[0].SecretKey != node.Token || config.Proxies[0].Plugin.Type != "http_proxy" || !config.Proxies[0].Transport.UseEncryption {
		t.Fatal("tunnel proxy must use authenticated, encrypted STCP")
	}
}

func TestBuildTunnelNodeConfigRequiresServerAddress(t *testing.T) {
	_, err := buildTunnelNodeConfig(&models.TunnelNode{}, &nodeconfigs.FRPConfig{}, "https://admin.example.com", nil)
	if err == nil {
		t.Fatal("expected an error when the cluster has no accessible server")
	}
}
