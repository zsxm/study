package httpsrv

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"gopkg.in/goyy/goyy.v0/util/times"
)

//检查字段
func (me *HttpService) check() (bool, int, string) {
	if me.Project == "" {
		return true, -1000, "项目为空"
	} else if me.Module == "" {
		return true, -1010, "模块为空"
	} else if me.ServiceUrl == "" {
		return true, -1020, "接口地址为空" //接口地址为空
	} else if me.Method == "" && (me.Method != METHOD_GET || me.Method != METHOD_POST) {
		return true, -1030, "请求方式为空 或不是 GET POST" //请求方式为空 或不是 GET POST
	} else if me.Ver == "" {
		return true, -1040, "接口版本为空" //接口版本为空
	} else if me.DataFormat != DATA_FORMAT_URL && me.DataFormat != DATA_FORMAT_XML && me.DataFormat != DATA_FORMAT_JSON {
		return true, -1050, "交互的数据格式字段不正确" //交互的数据格式字段不正确
	}
	return false, 0, ""
}

//http 请求
//参数map转换为请求body, key参数名 value参数值
func (me *HttpService) Send() {
	defer func() {
		me.Log()
		if err := recover(); err != nil {
			me.StatusCode = -1080
			me.StatusMsg = "未知错误"
			logger.Error(err)
			debug.PrintStack()
		}
	}()
	//验证请求基本字段信息
	b, s, r := me.check()
	if b { //b=true校验失败
		me.StatusCode = s
		me.StatusMsg = r
		return
	}
	me.ReqTime = times.NowUnix() //请求时间
	var client = &http.Client{
		Timeout: time.Duration(30 * time.Second), //设置超时 单位秒
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: me.HttpsTLS},
			//DisableCompression: true,//压缩
		},
	}

	var bodyReader *strings.Reader
	if me.Method == METHOD_POST {
		if me.Values != nil && len(me.Values) > 0 && me.RequestBody != "" {
			me.ServiceUrl += "?" + me.Values.Encode()
			bodyReader = strings.NewReader(me.RequestBody)
		} else if me.Values != nil && len(me.Values) > 0 && me.RequestBody == "" {
			bodyReader = strings.NewReader(me.Values.Encode())
		} else if me.RequestBody != "" {
			bodyReader = strings.NewReader(me.RequestBody)
		}
	} else if me.Method == METHOD_GET {
		if me.Values != nil && len(me.Values) > 0 {
			me.ServiceUrl += "?" + me.Values.Encode()
		} else if me.RequestBody != "" {
			me.ServiceUrl += "?" + me.RequestBody
		}
		bodyReader = strings.NewReader("")
	}

	logger.Println("send url :", me.ServiceUrl)
	logger.Println("request params :", me.Values.Encode())
	logger.Println("request body :", me.RequestBody)

	request, err := http.NewRequest(me.Method, me.ServiceUrl, bodyReader)
	if err != nil {
		logger.Errorln(err.Error())
		me.StatusCode = -1098
		me.StatusMsg = "创建请求失败"
		return
	}
	//执行请求
	if me.Method == METHOD_POST {
		var bodyType = "application/x-www-form-urlencoded"
		if me.DataFormat == DATA_FORMAT_XML {
			bodyType = "text/xml"
		} else if me.DataFormat == DATA_FORMAT_JSON {
			bodyType = "text/json"
		}
		request.Header.Set("Content-Type", bodyType)
	}

	//执行请求
	response, err := client.Do(request)
	if err != nil {
		logger.Errorln("client request err :", err.Error())
		me.StatusCode = -1097
		me.StatusMsg = err.Error() //执行请求失败或超时
		return
	}
	defer response.Body.Close()
	me.StatusCode = response.StatusCode
	logger.Println("response status code =", response.StatusCode)
	me.ResTime = times.NowUnix() //响应时间
	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body) //读取响应内容
		if err != nil {
			me.StatusCode = -1096
			me.StatusMsg = err.Error() //读取响应内容异常
		} else {
			me.ResponseBody = string(body) //响应内容
		}
	} else {
		me.StatusMsg = "服务请求出错"
	}
	logger.Println("response time :", me.ResTime, "-", me.ReqTime, "=", me.ResTime-me.ReqTime)
}

//json转换map
func (me *HttpService) JsonToMap() (map[string]interface{}, error) {
	var object interface{}
	err := json.Unmarshal([]byte(me.ResponseBody), &object)
	if err != nil {
		return make(map[string]interface{}), err
	} else if mmap, ok := object.(map[string]interface{}); ok {
		return mmap, nil
	}
	return make(map[string]interface{}), err
}

//map转换json
func (me *HttpService) MapToJson(mmap map[string]string) ([]byte, error) {
	jsn, err := json.Marshal(mmap)
	if err != nil {
		return nil, err
	}
	return jsn, err
}

//记录日志
func (me *HttpService) Log() {
	logger.Printf("http server : %+v\n", me)
}
