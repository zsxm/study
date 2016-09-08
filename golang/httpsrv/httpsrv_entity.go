package httpsrv

import (
	"net/url"
)

func New() *HttpService {
	e := &HttpService{DataFormat: DATA_FORMAT_URL}
	return e
}

type HttpService struct {
	Project, Module                string     //项目和模块
	ServiceUrl, ServerName, Method string     //Method = GET or POST
	ReqTime, ResTime               int64      //请求时间和响应时间
	RequestBody, ResponseBody      string     //请求数据和响应数据
	ResCode, ResMsg                string     //接口返回数据状态码和状态信息
	Ver                            string     //接口版本
	StatusCode                     int        //http状态码
	StatusMsg                      string     //http状态信息
	DataFormat                     string     //交互数据格式 URL、XML、JSON
	HttpsTLS                       bool       //https认证 true TLS接受任何证书
	Values                         url.Values //请求参数
}
