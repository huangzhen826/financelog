package common

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func ParamsSetData(objValuePtr reflect.Value, dataObj interface{}) error {
	value := objValuePtr.Elem() //获取value对象
	theType := value.Type()     //获取value的类型
	theKind := theType.Kind()   //获取类型的分类
	if theKind == reflect.Struct {
		for i := 0; i < theType.NumField(); i++ {
			varName := theType.Field(i).Name            //结构体中成员变量的名称
			tagName := theType.Field(i).Tag.Get("json") //json里面的key
			opt := theType.Field(i).Tag.Get("opt")      //必填获取
			checkreq := true
			if opt == "true" {
				checkreq = false
			}
			if !value.FieldByName(varName).CanSet() { //成员变量必须是public类型
				//panic(varName + " must be public")
				continue
			}

			if mapObj, ok := dataObj.(map[string]interface{}); ok { //结构体必须跟map对应
				if _, ok := mapObj[tagName]; ok { //
					if err := ParamsSetData(value.FieldByName(varName).Addr(), mapObj[tagName]); err != nil {
						return err
					} //加载成员变量的值
					continue
				} else if checkreq {
					return errors.New("缺少必要参数:" + tagName)
				}
			}
		}
	} else if theKind == reflect.Slice {
		if arrVal, ok := dataObj.([]interface{}); ok {
			itemSize := len(arrVal)
			value.Set(reflect.MakeSlice(theType, itemSize, itemSize)) //构造slice
			for i := 0; i < itemSize; i++ {
				if err := ParamsSetData(value.Index(i).Addr(), arrVal[i]); err != nil {
					return err
				}
			}
		}
	} else if theKind == reflect.String {
		value.SetString(fmt.Sprint(dataObj))
	} else if theKind == reflect.Int {
		intV, err := strconv.ParseInt(fmt.Sprint(dataObj), 10, 64) //字符串转int64
		if err != nil {
			intV = 0
		}
		value.SetInt(intV)
	} else if theKind == reflect.Float64 {
		floatV, err := strconv.ParseFloat(fmt.Sprint(dataObj), 64) //字符串转float64
		if err != nil {
			floatV = 0.0
		}
		value.SetFloat(floatV)
	} else if theKind == reflect.Map {
		if mapVal, ok := dataObj.(map[string]interface{}); ok {
			value.Set(reflect.MakeMap(reflect.TypeOf(dataObj))) //构造map
			for k, v := range mapVal {
				value.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
			}
		}
	}

	return nil
}
