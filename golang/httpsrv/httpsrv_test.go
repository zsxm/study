package httpsrv_test

import (
	"creditease.com/yxhn/hn-comm/api/httpsrv"

	"net/url"
	"testing"
)

func TestHttp(t *testing.T) {
	hs := httpsrv.New()
	hs.Module = "test"
	hs.Project = "test"
	hs.ServiceUrl = "http://10.100.140.23:9090/product/yxhn_all.json"
	hs.ServerName = "/product/yxhn_all.json"
	hs.Method = httpsrv.METHOD_GET
	hs.Ver = "1"
	values := url.Values{}
	values.Add("appkey", "yuwenbn")
	hs.Values = values
	hs.Send()
}
