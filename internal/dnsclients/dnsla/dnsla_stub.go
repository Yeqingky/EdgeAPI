// stub: 该包在开源版中被裁剪但 provider_dns_la.go 仍有引用, 此实现仅为通过编译, 功能不可用
package dnsla

type baseResp struct {
	ok  bool
	err error
}

func (this *baseResp) Success() bool { return this.ok }
func (this *baseResp) Error() error  { return this.err }

type DomainListResponse struct {
	baseResp
	Data struct {
		Results []struct {
			Domain string `json:"domain"`
		} `json:"results"`
	} `json:"data"`
}

type RecordListResponse struct {
	baseResp
	Data struct {
		Results []struct {
			Id       string `json:"id"`
			Host     string `json:"host"`
			Type     int    `json:"type"`
			Data     string `json:"data"`
			LineCode string `json:"line_code"`
			TTL      any    `json:"ttl"`
		} `json:"results"`
	} `json:"data"`
}

type AllLineListResponse struct {
	baseResp
	Data []struct {
		Id       string                     `json:"id"`
		Name     string                     `json:"name"`
		Code     string                     `json:"code"`
		Children []AllLineListResponseChild `json:"children"`
	} `json:"data"`
}

type AllLineListResponseChild struct {
	Id       string                     `json:"id"`
	Name     string                     `json:"name"`
	Code     string                     `json:"code"`
	Children []AllLineListResponseChild `json:"children"`
}

type DomainResponse struct {
	baseResp
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
}

type RecordCreateResponse struct {
	baseResp
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
}

type RecordUpdateResponse struct{ baseResp }

type RecordDeleteResponse struct {
	baseResp
	Code int `json:"code"`
}
