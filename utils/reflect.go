package utils

import (
	"errors"
	"reflect"
)

var errValueDiff = errors.New("newValue is a different type from obj, couldn't be set")

// ModifyIT 修改interface的数据
func ModifyIT(obj interface{}, value interface{}) error {
	v := reflect.ValueOf(obj)
	// 确认传入的是否是指针
	if v.Kind() == reflect.Ptr {
		// 获得所指向的元素
		v = v.Elem()
	}
	// 创建新值的reflect.Value
	valueReflect := reflect.ValueOf(value).Elem()
	// 确保新值和元素值类型一致
	if valueReflect.Type() == v.Type() {
		v.Set(valueReflect)
	} else {
		return errValueDiff
	}
	return nil
}

// StructToMap 将结构体转换为map
func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	objValue := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)
	for i := 0; i < objValue.NumField(); i++ {
		fieldName := objType.Field(i).Name
		fieldValue := objValue.Field(i).Interface()
		result[fieldName] = fieldValue
	}
	return result
}
