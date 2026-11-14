package models

import (
	"time"

	dbutils "git.ovh.im/yeqing/EdgeAPI/internal/db/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/types"
)

const (
	TunnelNodeStateEnabled  = 1
	TunnelNodeStateDisabled = 0
)

type TunnelNodeDAO dbs.DAO

func NewTunnelNodeDAO() *TunnelNodeDAO {
	return dbs.NewDAO(&TunnelNodeDAO{
		DAOObject: dbs.DAOObject{
			DB:     Tea.Env,
			Table:  "edgeTunnelNodes",
			Model:  new(TunnelNode),
			PkName: "id",
		},
	}).(*TunnelNodeDAO)
}

var SharedTunnelNodeDAO *TunnelNodeDAO

func init() {
	dbs.OnReady(func() {
		SharedTunnelNodeDAO = NewTunnelNodeDAO()
	})
}

func (this *TunnelNodeDAO) CreateTunnelNode(tx *dbs.Tx, adminId int64, userId int64, clusterId int64, name string, token string, serverAddr string, serverPort int, osName string, arch string) (int64, error) {
	var op = NewTunnelNodeOperator()
	op.AdminId = adminId
	op.UserId = userId
	op.ClusterId = clusterId
	op.Name = name
	op.Token = token
	op.ServerAddr = serverAddr
	op.ServerPort = serverPort
	op.OS = osName
	op.Arch = arch
	op.IsOn = true
	op.IsInstalled = false
	op.ConfigVersion = time.Now().Unix()
	op.CreatedAt = time.Now().Unix()
	op.UpdatedAt = time.Now().Unix()
	op.State = TunnelNodeStateEnabled
	err := this.Save(tx, op)
	if err != nil {
		return 0, err
	}
	return types.Int64(op.Id), nil
}

func (this *TunnelNodeDAO) FindEnabledTunnelNode(tx *dbs.Tx, id int64) (*TunnelNode, error) {
	one, err := this.Query(tx).
		Pk(id).
		State(TunnelNodeStateEnabled).
		Find()
	if one == nil || err != nil {
		return nil, err
	}
	return one.(*TunnelNode), nil
}

func (this *TunnelNodeDAO) FindTunnelNodeWithToken(tx *dbs.Tx, token string) (*TunnelNode, error) {
	one, err := this.Query(tx).
		Attr("token", token).
		State(TunnelNodeStateEnabled).
		Find()
	if one == nil || err != nil {
		return nil, err
	}
	return one.(*TunnelNode), nil
}

func (this *TunnelNodeDAO) ListTunnelNodes(tx *dbs.Tx, userId int64, clusterId int64, keyword string, offset int64, size int64) ([]*TunnelNode, error) {
	var query = this.Query(tx).State(TunnelNodeStateEnabled)
	if userId > 0 {
		query.Attr("userId", userId)
	}
	if clusterId > 0 {
		query.Attr("clusterId", clusterId)
	}
	if keyword != "" {
		query.Where("name LIKE :keyword").Param("keyword", dbutils.QuoteLike(keyword))
	}
	var result []*TunnelNode
	_, err := query.Offset(offset).Limit(size).DescPk().Slice(&result).FindAll()
	return result, err
}

func (this *TunnelNodeDAO) CountTunnelNodes(tx *dbs.Tx, userId int64, clusterId int64, keyword string) (int64, error) {
	var query = this.Query(tx).State(TunnelNodeStateEnabled)
	if userId > 0 {
		query.Attr("userId", userId)
	}
	if clusterId > 0 {
		query.Attr("clusterId", clusterId)
	}
	if keyword != "" {
		query.Where("name LIKE :keyword").Param("keyword", dbutils.QuoteLike(keyword))
	}
	return query.Count()
}

func (this *TunnelNodeDAO) DeleteTunnelNode(tx *dbs.Tx, id int64) error {
	_, err := this.Query(tx).
		Pk(id).
		Set("state", TunnelNodeStateDisabled).
		Update()
	return err
}

func (this *TunnelNodeDAO) UpdateTunnelNodeIsInstalled(tx *dbs.Tx, id int64, isInstalled bool) error {
	_, err := this.Query(tx).
		Pk(id).
		Set("isInstalled", isInstalled).
		Set("updatedAt", time.Now().Unix()).
		Update()
	return err
}

func (this *TunnelNodeDAO) UpdateTunnelNode(tx *dbs.Tx, id int64, name string, serverAddr string, serverPort int, isOn bool) error {
	var op = NewTunnelNodeOperator()
	op.Id = id
	op.Name = name
	op.ServerAddr = serverAddr
	op.ServerPort = serverPort
	op.IsOn = isOn
	op.UpdatedAt = time.Now().Unix()
	op.ConfigVersion = time.Now().Unix()
	return this.Save(tx, op)
}
