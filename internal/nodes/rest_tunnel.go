package nodes

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"git.ovh.im/yeqing/EdgeAPI/internal/db/models"
	"github.com/iwind/TeaGo/Tea"
)

func (this *RestServer) handleTunnel(writer http.ResponseWriter, req *http.Request) bool {
	switch req.URL.Path {
	case "/tunnels/download":
		this.handleTunnelDownload(writer, req)
		return true
	case "/tunnels/config":
		this.handleTunnelConfig(writer, req)
		return true
	case "/tunnels/logs":
		this.handleTunnelLogs(writer, req)
		return true
	}
	return false
}

func (this *RestServer) findTunnelNode(req *http.Request) (*models.TunnelNode, bool) {
	token := req.URL.Query().Get("token")
	if token == "" {
		return nil, false
	}
	node, err := models.SharedTunnelNodeDAO.FindTunnelNodeWithToken(nil, token)
	if err != nil || node == nil || !node.IsOn {
		return nil, false
	}
	return node, true
}

func (this *RestServer) handleTunnelDownload(writer http.ResponseWriter, req *http.Request) {
	node, ok := this.findTunnelNode(req)
	if !ok {
		http.Error(writer, "invalid token", http.StatusUnauthorized)
		return
	}

	osName := req.URL.Query().Get("os")
	arch := req.URL.Query().Get("arch")
	if osName == "" {
		osName = node.OS
	}
	if arch == "" {
		arch = node.Arch
	}
	if osName != "linux" && osName != "windows" {
		http.Error(writer, "unsupported os", http.StatusBadRequest)
		return
	}
	if arch != "amd64" && arch != "arm64" {
		http.Error(writer, "unsupported arch", http.StatusBadRequest)
		return
	}

	name := "frpc-" + osName + "-" + arch
	if osName == "windows" {
		name += ".exe"
	}
	path := filepath.Join(Tea.Root, "deploy", name)
	file, err := os.Open(path)
	if err != nil {
		http.Error(writer, "binary not found", http.StatusNotFound)
		return
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		http.Error(writer, "binary not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	writer.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
	_, _ = io.Copy(writer, file)
}

func (this *RestServer) handleTunnelConfig(writer http.ResponseWriter, req *http.Request) {
	node, ok := this.findTunnelNode(req)
	if !ok {
		http.Error(writer, "invalid token", http.StatusUnauthorized)
		return
	}

	cluster, err := models.SharedNodeClusterDAO.FindEnabledNodeCluster(nil, int64(node.ClusterId))
	if err != nil || cluster == nil {
		http.Error(writer, "cluster not found", http.StatusNotFound)
		return
	}
	frpConfig, err := models.SharedNodeClusterDAO.FindClusterFRP(nil, int64(node.ClusterId))
	if err != nil {
		http.Error(writer, "read frp config failed", http.StatusInternalServerError)
		return
	}
	if !frpConfig.IsOn {
		http.Error(writer, "frp is not enabled for this cluster", http.StatusServiceUnavailable)
		return
	}

	baseURL := this.requestBaseURL(req)
	serverAddrs, err := models.FindTunnelServerAddrs(nil, int64(node.ClusterId), frpConfig.BindPort)
	if err != nil {
		http.Error(writer, "read server addresses failed", http.StatusInternalServerError)
		return
	}
	data, err := node.BuildFRPCConfig(frpConfig, baseURL, serverAddrs)
	if err != nil {
		http.Error(writer, "encode config failed", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Tunnel-Config-Version", strconv.FormatUint(node.ConfigVersion, 10))
	_, _ = writer.Write(data)
}

func (this *RestServer) handleTunnelLogs(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	node, ok := this.findTunnelNode(req)
	if !ok {
		http.Error(writer, "invalid token", http.StatusUnauthorized)
		return
	}

	scanner := bufio.NewScanner(io.LimitReader(req.Body, 4<<20))
	scanner.Buffer(make([]byte, 4096), 64*1024)
	var count int
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry struct {
			Time    string `json:"time"`
			Level   string `json:"level"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		level := entry.Level
		if level == "" {
			level = "info"
		}
		message := entry.Message
		if len(message) > 16384 {
			message = message[:16384]
		}
		err := models.SharedTunnelLogDAO.CreateTunnelLog(nil, int64(node.Id), level, message)
		if err != nil {
			http.Error(writer, "write log failed", http.StatusInternalServerError)
			return
		}
		count++
		if count >= 1000 {
			break
		}
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (this *RestServer) requestBaseURL(req *http.Request) string {
	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	if forwarded := req.Header.Get("X-Forwarded-Proto"); forwarded != "" {
		scheme = forwarded
	}
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	return scheme + "://" + host
}
