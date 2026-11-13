package dbutils_test

import (
	"testing"

	dbutils "git.ovh.im/yeqing/EdgeAPI/internal/db/utils"
	_ "github.com/iwind/TeaGo/bootstrap"
)

func TestHasFreeSpace(t *testing.T) {
	t.Log(dbutils.CheckHasFreeSpace())
	t.Log(dbutils.LocalDatabaseDataDir)
}
