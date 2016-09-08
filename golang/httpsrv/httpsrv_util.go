package httpsrv

import (
	"encoding/json"
)

//json转换map
func JsonToMap(data string) (map[string]interface{}, error) {
	var object interface{}
	err := json.Unmarshal([]byte(data), &object)
	if err != nil {
		return make(map[string]interface{}), err
	} else if mmap, ok := object.(map[string]interface{}); ok {
		return mmap, nil
	}
	return make(map[string]interface{}), err
}

//map转换json
func MapToJson(mmap map[string]interface{}) ([]byte, error) {
	jsn, err := json.Marshal(mmap)
	if err != nil {
		return nil, err
	}
	return jsn, err
}
