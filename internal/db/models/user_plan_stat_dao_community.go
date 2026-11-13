//go:build !plus

package models

import "github.com/iwind/TeaGo/dbs"

func (this *UserPlanStatDAO) IncreaseUserPlanStat(tx *dbs.Tx, userPlanId int64, trafficBytes int64, countRequests int64, countWebsocketConnections int64) error {
	// stub
	return nil
}

func (this *UserPlanStatDAO) ResetUserPlanStatsWithUserPlanId(tx *dbs.Tx, userPlanId int64) error {
	// stub
	return nil
}
