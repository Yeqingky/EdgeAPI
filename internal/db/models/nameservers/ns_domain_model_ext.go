package nameservers

import (
	"encoding/json"

	"git.ovh.im/yeqing/EdgeAPI/internal/db/models"
	"git.ovh.im/yeqing/EdgeAPI/internal/remotelogs"
)

func (this *NSDomain) DecodeGroupIds() []int64 {
	if models.IsNull(this.GroupIds) {
		return nil
	}

	var result = []int64{}
	err := json.Unmarshal(this.GroupIds, &result)
	if err != nil {
		remotelogs.Error("NSDomain", "DecodeGroupIds:"+err.Error())
	}
	return result
}
