package tasks

import (
	"git.ovh.im/yeqing/EdgeAPI/internal/db/models"
	"git.ovh.im/yeqing/EdgeAPI/internal/remotelogs"
)

type BaseTask struct {
}

func (this *BaseTask) logErr(taskType string, errString string) {
	remotelogs.Error("TASK", "run '"+taskType+"' failed: "+errString)
}

func (this *BaseTask) IsPrimaryNode() bool {
	return models.SharedAPINodeDAO.CheckAPINodeIsPrimaryWithoutErr()
}
