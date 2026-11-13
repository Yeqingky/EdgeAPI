package models

import (
	"encoding/json"

	"git.ovh.im/yeqing/EdgeAPI/internal/remotelogs"
	"git.ovh.im/yeqing/EdgeCommon/pkg/serverconfigs/sslconfigs"
)

func (this *SSLPolicy) DecodeCerts() []*sslconfigs.SSLCertRef {
	if len(this.Certs) == 0 {
		return nil
	}

	var refs = []*sslconfigs.SSLCertRef{}
	err := json.Unmarshal(this.Certs, &refs)
	if err != nil {
		remotelogs.Error("SSLPolicy_DecodeCerts", err.Error())
	}
	return refs
}
