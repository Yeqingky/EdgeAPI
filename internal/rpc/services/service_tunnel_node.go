package services

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"

	"git.ovh.im/yeqing/EdgeAPI/internal/db/models"
	"git.ovh.im/yeqing/EdgeCommon/pkg/nodeconfigs"
	"git.ovh.im/yeqing/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/types"
)

type TunnelNodeService struct {
	BaseService
}

func (this *TunnelNodeService) CreateTunnelNode(ctx context.Context, req *pb.CreateTunnelNodeRequest) (*pb.CreateTunnelNodeResponse, error) {
	adminId, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(req.Name)) == 0 {
		return nil, errors.New("name should not be empty")
	}
	if req.ClusterId <= 0 {
		return nil, errors.New("clusterId should not be empty")
	}
	if req.ServerPort <= 0 || req.ServerPort > 65535 {
		return nil, errors.New("invalid serverPort")
	}
	if net.ParseIP(req.ServerAddr) == nil && len(strings.TrimSpace(req.ServerAddr)) == 0 {
		return nil, errors.New("serverAddr should not be empty")
	}
	if req.Os != "linux" && req.Os != "windows" {
		return nil, errors.New("unsupported os")
	}
	if req.Arch != "amd64" && req.Arch != "arm64" {
		return nil, errors.New("unsupported arch")
	}

	var tx = this.NullTx()
	cluster, err := models.SharedNodeClusterDAO.FindEnabledNodeCluster(tx, req.ClusterId)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, errors.New("cluster not found")
	}
	frpConfig, err := models.SharedNodeClusterDAO.FindClusterFRP(tx, req.ClusterId)
	if err != nil {
		return nil, err
	}
	if !frpConfig.IsOn {
		return nil, errors.New("frp is not enabled for this cluster")
	}

	userId := req.UserId
	if userId <= 0 {
		userId, err = models.SharedAdminDAO.FindAdminUserId(tx, adminId)
		if err != nil {
			return nil, err
		}
		if userId <= 0 {
			userId, err = models.EnsureLinkedUserForAdmin(tx, adminId)
			if err != nil {
				return nil, err
			}
		}
	} else {
		exist, err := models.SharedUserDAO.Exist(tx, userId)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, errors.New("user not found")
		}
		userClusterId, err := models.SharedUserDAO.FindUserClusterId(tx, userId)
		if err != nil {
			return nil, err
		}
		if userClusterId > 0 && userClusterId != req.ClusterId {
			return nil, errors.New("cluster does not belong to the selected user")
		}
	}

	var token string
	for range 5 {
		token = rands.String(64)
		exist, err := models.SharedTunnelNodeDAO.FindTunnelNodeWithToken(tx, token)
		if err != nil {
			return nil, err
		}
		if exist == nil {
			break
		}
		token = ""
	}
	if token == "" {
		return nil, errors.New("generate token failed")
	}

	tunnelNodeId, err := models.SharedTunnelNodeDAO.CreateTunnelNode(tx,
		adminId, userId, req.ClusterId, strings.TrimSpace(req.Name), token,
		strings.TrimSpace(req.ServerAddr), int(req.ServerPort), req.Os, req.Arch)
	if err != nil {
		return nil, err
	}
	return &pb.CreateTunnelNodeResponse{
		TunnelNodeId: tunnelNodeId,
		Token:        token,
	}, nil
}

func (this *TunnelNodeService) ListTunnelNodes(ctx context.Context, req *pb.ListTunnelNodesRequest) (*pb.ListTunnelNodesResponse, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var tx = this.NullTx()
	nodes, err := models.SharedTunnelNodeDAO.ListTunnelNodes(tx, req.UserId, req.ClusterId, req.Keyword, req.Offset, req.Size)
	if err != nil {
		return nil, err
	}
	result := make([]*pb.TunnelNode, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, this.toPB(node))
	}
	return &pb.ListTunnelNodesResponse{TunnelNodes: result}, nil
}

func (this *TunnelNodeService) CountTunnelNodes(ctx context.Context, req *pb.CountTunnelNodesRequest) (*pb.RPCCountResponse, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var tx = this.NullTx()
	count, err := models.SharedTunnelNodeDAO.CountTunnelNodes(tx, req.UserId, req.ClusterId, req.Keyword)
	if err != nil {
		return nil, err
	}
	return &pb.RPCCountResponse{Count: count}, nil
}

func (this *TunnelNodeService) FindTunnelNode(ctx context.Context, req *pb.FindTunnelNodeRequest) (*pb.FindTunnelNodeResponse, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var tx = this.NullTx()
	node, err := models.SharedTunnelNodeDAO.FindEnabledTunnelNode(tx, req.TunnelNodeId)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return &pb.FindTunnelNodeResponse{}, nil
	}
	return &pb.FindTunnelNodeResponse{TunnelNode: this.toPB(node)}, nil
}

func (this *TunnelNodeService) FindTunnelNodeWithToken(ctx context.Context, req *pb.FindTunnelNodeWithTokenRequest) (*pb.FindTunnelNodeWithTokenResponse, error) {
	var tx = this.NullTx()
	node, err := models.SharedTunnelNodeDAO.FindTunnelNodeWithToken(tx, req.Token)
	if err != nil {
		return nil, err
	}
	if node == nil || !node.IsOn {
		return &pb.FindTunnelNodeWithTokenResponse{}, nil
	}
	return &pb.FindTunnelNodeWithTokenResponse{TunnelNode: this.toPB(node)}, nil
}

func (this *TunnelNodeService) FindTunnelNodeConfigWithToken(ctx context.Context, req *pb.FindTunnelNodeConfigWithTokenRequest) (*pb.FindTunnelNodeConfigWithTokenResponse, error) {
	var tx = this.NullTx()
	node, err := models.SharedTunnelNodeDAO.FindTunnelNodeWithToken(tx, req.Token)
	if err != nil {
		return nil, err
	}
	if node == nil || !node.IsOn {
		return nil, errors.New("invalid tunnel token")
	}
	frpConfig, err := models.SharedNodeClusterDAO.FindClusterFRP(tx, int64(node.ClusterId))
	if err != nil {
		return nil, err
	}
	if !frpConfig.IsOn {
		return nil, errors.New("frp is not enabled for this cluster")
	}
	data, err := buildTunnelNodeConfig(node, frpConfig, req.BaseURL)
	if err != nil {
		return nil, err
	}
	return &pb.FindTunnelNodeConfigWithTokenResponse{
		ConfigJSON:    data,
		ConfigVersion: int64(node.ConfigVersion),
	}, nil
}

func buildTunnelNodeConfig(node *models.TunnelNode, frpConfig *nodeconfigs.FRPConfig, baseURL string) ([]byte, error) {
	baseURL = strings.TrimRight(baseURL, "/")
	config := map[string]any{
		"serverAddr": node.ServerAddr,
		"serverPort": node.ServerPort,
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
				"token": node.Token,
			},
		},
		"remoteConfig": map[string]any{
			"url":      baseURL + "/tunnels/config",
			"token":    node.Token,
			"interval": 30,
		},
	}
	return json.Marshal(config)
}

func (this *TunnelNodeService) DeleteTunnelNode(ctx context.Context, req *pb.DeleteTunnelNodeRequest) (*pb.RPCSuccess, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var tx = this.NullTx()
	err = models.SharedTunnelNodeDAO.DeleteTunnelNode(tx, req.TunnelNodeId)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

func (this *TunnelNodeService) UpdateTunnelNode(ctx context.Context, req *pb.UpdateTunnelNodeRequest) (*pb.RPCSuccess, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(req.Name)) == 0 {
		return nil, errors.New("name should not be empty")
	}
	if req.ServerPort <= 0 || req.ServerPort > 65535 {
		return nil, errors.New("invalid serverPort")
	}
	var tx = this.NullTx()
	err = models.SharedTunnelNodeDAO.UpdateTunnelNode(tx, req.TunnelNodeId, strings.TrimSpace(req.Name), strings.TrimSpace(req.ServerAddr), int(req.ServerPort), req.IsOn)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

func (this *TunnelNodeService) UpdateTunnelNodeIsInstalled(ctx context.Context, req *pb.UpdateTunnelNodeIsInstalledRequest) (*pb.RPCSuccess, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var tx = this.NullTx()
	err = models.SharedTunnelNodeDAO.UpdateTunnelNodeIsInstalled(tx, req.TunnelNodeId, req.IsInstalled)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

func (this *TunnelNodeService) CreateTunnelLog(ctx context.Context, req *pb.CreateTunnelLogRequest) (*pb.RPCSuccess, error) {
	var tx = this.NullTx()
	node, err := models.SharedTunnelNodeDAO.FindTunnelNodeWithToken(tx, req.Token)
	if err != nil {
		return nil, err
	}
	if node == nil || !node.IsOn {
		return nil, errors.New("invalid tunnel token")
	}
	if len(req.Logs) > 1000 {
		req.Logs = req.Logs[:1000]
	}
	for _, log := range req.Logs {
		message := log.Message
		if len(message) > 16384 {
			message = message[:16384]
		}
		err = models.SharedTunnelLogDAO.CreateTunnelLog(tx, int64(node.Id), log.Level, message)
		if err != nil {
			return nil, err
		}
	}
	return this.Success()
}

func (this *TunnelNodeService) ListTunnelLogs(ctx context.Context, req *pb.ListTunnelLogsRequest) (*pb.ListTunnelLogsResponse, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var tx = this.NullTx()
	logs, err := models.SharedTunnelLogDAO.ListTunnelLogs(tx, req.TunnelNodeId, req.Offset, req.Size)
	if err != nil {
		return nil, err
	}
	result := make([]*pb.TunnelLog, 0, len(logs))
	for _, log := range logs {
		result = append(result, &pb.TunnelLog{
			Id:           int64(log.Id),
			TunnelNodeId: int64(log.TunnelNodeId),
			Level:        log.Level,
			Message:      log.Message,
			CreatedAt:    int64(log.CreatedAt),
		})
	}
	return &pb.ListTunnelLogsResponse{Logs: result}, nil
}

func (this *TunnelNodeService) toPB(node *models.TunnelNode) *pb.TunnelNode {
	if node == nil {
		return nil
	}
	return &pb.TunnelNode{
		Id:            int64(node.Id),
		ClusterId:     int64(node.ClusterId),
		UserId:        int64(node.UserId),
		Name:          node.Name,
		Token:         node.Token,
		ServerAddr:    node.ServerAddr,
		ServerPort:    types.Int32(node.ServerPort),
		Os:            node.OS,
		Arch:          node.Arch,
		Version:       node.Version,
		IsOn:          node.IsOn,
		IsInstalled:   node.IsInstalled,
		ConfigVersion: int64(node.ConfigVersion),
		CreatedAt:     int64(node.CreatedAt),
		UpdatedAt:     int64(node.UpdatedAt),
	}
}
