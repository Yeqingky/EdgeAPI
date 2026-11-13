package services

import (
	"context"

	"git.ovh.im/yeqing/EdgeAPI/internal/db/models/acme"
	"git.ovh.im/yeqing/EdgeAPI/internal/errors"
	"git.ovh.im/yeqing/EdgeCommon/pkg/rpc/pb"
)

// ACME认证相关
type ACMEAuthenticationService struct {
	BaseService
}

// 获取Key
func (this *ACMEAuthenticationService) FindACMEAuthenticationKeyWithToken(ctx context.Context, req *pb.FindACMEAuthenticationKeyWithTokenRequest) (*pb.FindACMEAuthenticationKeyWithTokenResponse, error) {
	_, err := this.ValidateNode(ctx)
	if err != nil {
		return nil, err
	}
	if len(req.Token) == 0 {
		return nil, errors.New("'token' should not be empty")
	}

	var tx = this.NullTx()

	auth, err := acme.SharedACMEAuthenticationDAO.FindAuthWithToken(tx, req.Token)
	if err != nil {
		return nil, err
	}
	if auth == nil {
		return &pb.FindACMEAuthenticationKeyWithTokenResponse{Key: ""}, nil
	}
	return &pb.FindACMEAuthenticationKeyWithTokenResponse{Key: auth.Key}, nil
}
