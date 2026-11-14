package models

import (
	"errors"

	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/types"
)

// EnsureLinkedUserForAdminDB 在升级/初始化场景 (裸 dbs.DB) 确保管理员拥有关联用户
func EnsureLinkedUserForAdminDB(db *dbs.DB, adminId int64) (int64, error) {
	if adminId <= 0 {
		return 0, errors.New("invalid adminId")
	}

	// 已关联的用户
	linkedUserIdCol, err := db.FindCol(0, "SELECT userId FROM edgeAdmins WHERE id=?", adminId)
	if err != nil {
		return 0, err
	}
	linkedUserId := types.Int64(linkedUserIdCol)
	if linkedUserId > 0 {
		one, err := db.FindOne("SELECT id FROM edgeUsers WHERE id=? AND state=1", linkedUserId)
		if err != nil {
			return 0, err
		}
		if len(one) > 0 {
			return linkedUserId, nil
		}
	}

	// 读取管理员
	admin, err := db.FindOne("SELECT id, username, password, fullname FROM edgeAdmins WHERE id=? AND state=1", adminId)
	if err != nil {
		return 0, err
	}
	if len(admin) == 0 {
		return 0, errors.New("admin not found")
	}
	username := admin.GetString("username")
	password := admin.GetString("password")
	fullname := admin.GetString("fullname")
	if len(username) == 0 {
		return 0, errors.New("admin username is empty")
	}

	// 同名用户直接采纳
	existingUserIdCol, err := db.FindCol(0, "SELECT id FROM edgeUsers WHERE username=? ORDER BY id ASC LIMIT 1", username)
	if err != nil {
		return 0, err
	}
	existingUserId := types.Int64(existingUserIdCol)
	if existingUserId > 0 {
		_, err = db.Exec("UPDATE edgeAdmins SET userId=? WHERE id=?", existingUserId, adminId)
		if err != nil {
			return 0, err
		}
		return existingUserId, nil
	}

	// 创建关联用户, 密码直接沿用管理员的MD5
	clusterIdCol, err := db.FindCol(0, "SELECT id FROM edgeNodeClusters WHERE state=1 ORDER BY id ASC LIMIT 1")
	if err != nil {
		return 0, err
	}
	clusterId := types.Int64(clusterIdCol)

	_, err = db.Exec("INSERT INTO edgeUsers (`isOn`, `username`, `password`, `fullname`, `source`, `clusterId`, `features`, `isVerified`, `state`, `createdAt`, `day`, `updatedAt`) VALUES (1, ?, ?, ?, 'admin', ?, '[]', 1, 1, UNIX_TIMESTAMP(), DATE_FORMAT(NOW(), '%Y%m%d'), UNIX_TIMESTAMP())", username, password, fullname, clusterId)
	if err != nil {
		return 0, err
	}
	userIdCol, err := db.FindCol(0, "SELECT id FROM edgeUsers WHERE username=? ORDER BY id ASC LIMIT 1", username)
	if err != nil {
		return 0, err
	}
	userId := types.Int64(userIdCol)

	_, err = db.Exec("UPDATE edgeAdmins SET userId=? WHERE id=?", userId, adminId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
