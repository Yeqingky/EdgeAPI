package models

import (
	"errors"

	"github.com/iwind/TeaGo/dbs"
)

// EnsureLinkedUserForAdmin 确保管理员拥有关联的普通用户 (用户名与管理员同步), 返回关联用户ID
// 已存在有效的关联用户时直接返回; 存在同名的普通用户时直接采纳; 否则创建新的关联用户
func EnsureLinkedUserForAdmin(tx *dbs.Tx, adminId int64) (int64, error) {
	if adminId <= 0 {
		return 0, errors.New("invalid adminId")
	}

	// 已关联的用户
	linkedUserId, err := SharedAdminDAO.FindAdminUserId(tx, adminId)
	if err != nil {
		return 0, err
	}
	if linkedUserId > 0 {
		exists, err := SharedUserDAO.Exist(tx, linkedUserId)
		if err != nil {
			return 0, err
		}
		if exists {
			return linkedUserId, nil
		}
	}

	// 读取管理员
	admin, err := SharedAdminDAO.FindEnabledAdmin(tx, adminId)
	if err != nil {
		return 0, err
	}
	if admin == nil {
		return 0, errors.New("admin not found")
	}

	// 同名用户直接采纳
	existingUserId, err := SharedUserDAO.FindAnyUserIdWithUsername(tx, admin.Username)
	if err != nil {
		return 0, err
	}
	if existingUserId > 0 {
		err = SharedAdminDAO.UpdateAdminUserId(tx, adminId, existingUserId)
		if err != nil {
			return 0, err
		}
		return existingUserId, nil
	}

	// 创建关联用户, 密码与管理员同步 (均为MD5)
	clusterId, err := SharedNodeClusterDAO.FindFirstEnabledClusterId(tx)
	if err != nil {
		return 0, err
	}
	userId, err := SharedUserDAO.CreateLinkedUser(tx, admin.Username, admin.Password, admin.Fullname, clusterId)
	if err != nil {
		return 0, err
	}

	err = SharedAdminDAO.UpdateAdminUserId(tx, adminId, userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

// FindAdminIdWithUserId 根据关联用户ID查找管理员ID
func (this *AdminDAO) FindAdminIdWithUserId(tx *dbs.Tx, userId int64) (int64, error) {
	if userId <= 0 {
		return 0, nil
	}
	return this.Query(tx).
		ResultPk().
		State(AdminStateEnabled).
		Attr("userId", userId).
		AscPk().
		FindInt64Col(0)
}
