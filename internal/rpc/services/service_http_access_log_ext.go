//go:build !plus

package services

import "git.ovh.im/yeqing/EdgeCommon/pkg/rpc/pb"

func (this *HTTPAccessLogService) canWriteAccessLogsToDB() bool {
	return true
}

func (this *HTTPAccessLogService) writeAccessLogsToPolicy(pbAccessLogs []*pb.HTTPAccessLog) error {
	return nil
}
