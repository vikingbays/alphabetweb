// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package webutils

import (
	"fmt"
	"reflect"
)

var flagOfPrintEncodeHttpParams = false

/*
把 bean 转换成 http param 结构,字符串形式，例如： s1=1&s2=2。
其中 bean 支持的类型： int, int32, int64, float32, float64, bool, string,
其中 bean 支持的类型： []int, []int32, []int64, []float32, []float64, []bool, []string
其中 bean 支持的类型： [][]yy (不建议，会额外增加坐标)
其中 bean 支持的类型： map[xx]yy, 其中 xx 是 int, int32, int64, float32, float64, bool, string
其中 bean 支持的类型： map[xx]yy, 其中 yy 是 int, int32, int64, float32, float64, bool, string ,
                            []int, []int32, []int64, []float32, []float64, []bool, []string,
                            struct, *struct
其中 bean 支持的类型： struct, *struct, []struct(不建议，会额外增加坐标), []*struct(不建议，会额外增加坐标)
*/
func EncodeHttpParams(beanValue reflect.Value) string {
	paramStrings := encodeHttpParamsStruct(beanValue, "", "", nil)
	if len(paramStrings) > 0 {
		paramStrings = paramStrings[1:]
	}

	return paramStrings

}

/*
把 bean 转换成 http param 结构, map[string][]string形式，例如： [s1:[1],s2:[2]]。
*/
func EncodeHttpParamsMap(beanValue reflect.Value) map[string][]string {
	returnParamStringsMap := make(map[string][]string)
	encodeHttpParamsStruct(beanValue, "", "", returnParamStringsMap)
	return returnParamStringsMap

}

/*
func EncodeHttpParams(bean interface{}) string {

	if bean != nil {
		beanValue := reflect.ValueOf(bean)
		paramStrings := encodeHttpParamsStruct(beanValue, "", "")
		paramStrings = paramStrings[1:]
		return paramStrings
	} else {
		return ""
	}

}
*/

func encodeHttpParamsStruct(beanValue reflect.Value, root string, fmtString string, returnParamStringsMap map[string][]string) string {
	if beanValue.Kind() == reflect.Ptr {
		beanValue = beanValue.Elem()
	}

	if beanValue.Kind() == reflect.Struct {
		beanValueType := beanValue.Type()
		for i := 0; i < beanValue.NumField(); i++ {
			name := getHttpParamAlias(beanValueType.Field(i), root)
			if beanValue.Field(i).Kind() == reflect.Slice { // 支持： 基础类型数组（[]string） , struct数组（[]struct）,struct指针数组[]*struct）
				fmtString = encodeHttpParamsSlice(beanValue.Field(i), name, fmtString, returnParamStringsMap)
			} else if beanValue.Field(i).Kind() == reflect.Map {
				fmtString = encodeHttpParamsMap(beanValue.Field(i), name, fmtString, returnParamStringsMap)
			} else if beanValue.Field(i).Kind() == reflect.Struct {
				fmtString = encodeHttpParamsStruct(beanValue.Field(i), name, fmtString, returnParamStringsMap)
			} else if beanValue.Field(i).Kind() == reflect.Ptr {
				fmtString = encodeHttpParamsStruct(beanValue.Field(i), name, fmtString, returnParamStringsMap)
			} else {
				fmtString = encodeHttpParamsBaseType(beanValue.Field(i), name, fmtString, returnParamStringsMap)
			}
		}
	}

	return fmtString
}

func encodeHttpParamsMap(beanValue reflect.Value, root string, fmtString string, returnParamStringsMap map[string][]string) string {
	valueMapType := beanValue.Type().Elem() // map[xx]yy , yy的类型
	keyMaps := beanValue.MapKeys()
	for _, k1 := range keyMaps {
		if valueMapType.Kind() == reflect.Slice {
			fmtString = encodeHttpParamsSlice(beanValue.MapIndex(k1), fmt.Sprintf("%s.%v", root, k1), fmtString, returnParamStringsMap)
		} else if valueMapType.Kind() == reflect.Map {
			fmtString = encodeHttpParamsMap(beanValue.MapIndex(k1), fmt.Sprintf("%s.%v", root, k1), fmtString, returnParamStringsMap)
		} else if valueMapType.Kind() == reflect.Struct {
			fmtString = encodeHttpParamsStruct(beanValue.MapIndex(k1), fmt.Sprintf("%s.%v", root, k1), fmtString, returnParamStringsMap)
		} else if valueMapType.Kind() == reflect.Ptr {
			fmtString = encodeHttpParamsStruct(beanValue.MapIndex(k1), fmt.Sprintf("%s.%v", root, k1), fmtString, returnParamStringsMap)
		} else {
			fmtString = encodeHttpParamsBaseType(beanValue.MapIndex(k1), fmt.Sprintf("%s.%v", root, k1), fmtString, returnParamStringsMap)
		}
	}
	return fmtString
}

func encodeHttpParamsSlice(beanValue reflect.Value, root string, fmtString string, returnParamStringsMap map[string][]string) string {
	for j := 0; j < beanValue.Len(); j++ {
		if beanValue.Index(j).Kind() == reflect.Struct {
			fmtString = encodeHttpParamsStruct(beanValue.Index(j), fmt.Sprintf("%s.%d", root, j), fmtString, returnParamStringsMap)
		} else if beanValue.Index(j).Kind() == reflect.Ptr {
			fmtString = encodeHttpParamsStruct(beanValue.Index(j), fmt.Sprintf("%s.%d", root, j), fmtString, returnParamStringsMap)
		} else if beanValue.Index(j).Kind() == reflect.Slice {
			fmtString = encodeHttpParamsSlice(beanValue.Index(j), fmt.Sprintf("%s.%d", root, j), fmtString, returnParamStringsMap)
		} else {
			fmtString = encodeHttpParamsBaseType(beanValue.Index(j), root, fmtString, returnParamStringsMap)

		}
	}
	return fmtString
}

func encodeHttpParamsBaseType(beanValue reflect.Value, root string, fmtString string, returnParamStringsMap map[string][]string) string {
	if beanValue.Type().Kind() == reflect.String {
		if beanValue.String() == "" {
			return fmtString
		}
	}
	if returnParamStringsMap != nil {
		if returnParamStringsMap[root] == nil {
			returnParamStringsMap[root] = make([]string, 0, 1)
		}
		returnParamStringsMap[root] = append(returnParamStringsMap[root], fmt.Sprintf("%v", beanValue))
	} else if flagOfPrintEncodeHttpParams {
		fmtString = fmt.Sprintf("%s\n&%s=%v", fmtString, root, beanValue)
	} else {
		fmtString = fmt.Sprintf("%s&%s=%v", fmtString, root, beanValue)
	}

	return fmtString
}

func getHttpParamAlias(fieldType reflect.StructField, root string) string {
	name := fieldType.Name

	alias := fieldType.Tag.Get("alias")
	if alias != "" {
		name = alias
	}
	if root != "" {
		name = root + "." + name
	}
	return name
}

func isBaseType(type1 reflect.Type) bool {
	if type1.Kind() != reflect.Slice &&
		type1.Kind() != reflect.Map &&
		type1.Kind() != reflect.Struct {
		return true
	} else {
		return false
	}
}
