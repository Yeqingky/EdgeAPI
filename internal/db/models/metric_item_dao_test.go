package models_test

import (
	"testing"

	"git.ovh.im/yeqing/EdgeAPI/internal/db/models"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/iwind/TeaGo/bootstrap"
)

func TestMetricStatDAO_Clean(t *testing.T) {
	var dao = models.NewMetricStatDAO()
	t.Log(dao.Clean(nil))
}
