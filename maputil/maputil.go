package maputil

import (
	"github.com/liuchonglin/go-utils"
	"reflect"
	"encoding/json"
)

// 通过反射方式完成结构体转换为map
// 反射方式比json方法效率高两倍左右
func StructToMapReflect(o interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	if !utils.IsStruct(o) {
		return m
	}
	v := reflect.ValueOf(o)
	if utils.IsPtr(o) {
		elem := v.Elem()
		t := elem.Type()
		for i := 0; i < t.NumField(); i++ {
			m[t.Field(i).Name] = elem.Field(i).Interface()
		}
	} else {
		t := reflect.TypeOf(o)
		for i := 0; i < t.NumField(); i++ {
			m[t.Field(i).Name] = v.Field(i).Interface()
		}
	}
	return m
}

// 通过json方式完成结构体转换为map
// json方式比反射方式效率低两倍左右，推荐用StructToMapReflect()方法
func StructToMapJson(o interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	if !utils.IsStruct(o) {
		return m
	}

	jsonData, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return nil
	}
	return m
}
