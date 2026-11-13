package models

import (
	"encoding/json"

	"git.ovh.im/yeqing/EdgeCommon/pkg/rpc/pb"
)

// ToPB 转换成PB对象
func (this *HTTPAccessLog) ToPB() (*pb.HTTPAccessLog, error) {
	p := &pb.HTTPAccessLog{}
	err := json.Unmarshal(this.Content, p)
	if err != nil {
		return nil, err
	}
	p.RequestId = this.RequestId
	p.RequestBody = this.RequestBody
	return p, nil
}
