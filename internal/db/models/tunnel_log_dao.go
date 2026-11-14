package models

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/dbs"
)

type TunnelLogDAO dbs.DAO

func NewTunnelLogDAO() *TunnelLogDAO {
	return dbs.NewDAO(&TunnelLogDAO{
		DAOObject: dbs.DAOObject{
			DB:     Tea.Env,
			Table:  "edgeTunnelLogs",
			Model:  new(TunnelLog),
			PkName: "id",
		},
	}).(*TunnelLogDAO)
}

var SharedTunnelLogDAO *TunnelLogDAO

func init() {
	dbs.OnReady(func() {
		SharedTunnelLogDAO = NewTunnelLogDAO()
	})
}

func (this *TunnelLogDAO) CreateTunnelLog(tx *dbs.Tx, tunnelNodeId int64, level string, message string) error {
	var op = NewTunnelLogOperator()
	op.TunnelNodeId = tunnelNodeId
	op.Level = level
	op.Message = message
	op.CreatedAt = time.Now().Unix()
	return this.Save(tx, op)
}

func (this *TunnelLogDAO) ListTunnelLogs(tx *dbs.Tx, tunnelNodeId int64, offset int64, size int64) ([]*TunnelLog, error) {
	var result []*TunnelLog
	_, err := this.Query(tx).
		Attr("tunnelNodeId", tunnelNodeId).
		Offset(offset).
		Limit(size).
		DescPk().
		Slice(&result).
		FindAll()
	return result, err
}
