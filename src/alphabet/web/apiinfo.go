// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"reflect"
)

type ApiInfo struct {
	AppName      string
	CodePkgPath  string
	CodeFuncName string
	WebUrl       string
	WebParams    []ApiInfo_Param
}

type ApiInfo_Param struct {
	FormName     string
	Alias        string
	TypeInfo     string
	TypeInfoKind reflect.Kind
	Doc          string
	IsArray      bool
}

/*

func GenerateApiActionInfo() []ApiInfo {
	apiList := make([]ApiInfo, 0, 1)
	log4go.InfoLog(routeActionerList)
	for _, routeActioner := range routeActionerList {
		apiInfo := ApiInfo{}
		apiInfo.AppName = routeActioner.Appname
		apiInfo.WebUrl = routeActioner.Url
		if routeActioner.Action1 != nil {
			action1Type := reflect.TypeOf(routeActioner.Action1)
			apiInfo.CodePkgPath = action1Type.PkgPath()
			apiInfo.CodeFuncName = action1Type.Name()
			apiInfo.WebParams = nil
		} else if routeActioner.Action2 != nil {
			action2Type := reflect.TypeOf(routeActioner.Action2)
			log4go.InfoLog(action2Type)
			param0Type := action2Type.In(0)
			log4go.InfoLog(param0Type)
			if param0Type.Kind() == reflect.Ptr &&
				param0Type.Elem().Kind() == reflect.Struct { // 必须是结构体
				param0Value := reflect.New(param0Type.Elem())
				log4go.InfoLog(param0Value)
				log4go.InfoLog(param0Value.Elem())

				log4go.InfoLog(webutils.EncodeHttpParamsMap(param0Value))

				//apiInfo.WebParams = webutils.EncodeHttpParams(param0Value)
			} else {
				apiInfo.WebParams = nil
			}
		}
		apiList = append(apiList, apiInfo)

	}

	bytes, err := json.Marshal(apiList)
	if err == nil {
		log4go.InfoLog(string(bytes))
	} else {
		log4go.InfoLog(err)
	}

	return apiList
}

func parseApiParams(beanType reflect.Type) []ApiInfo_Param {
	if beanType.Kind() == reflect.Ptr {
		beanType = beanType.Elem()
	}
	var beanValue reflect.Value
	if beanType.Kind() == reflect.Struct {
		beanValue := reflect.New(beanType)
	} else {
		log4go.ErrorLog(message.ERR_WEB0_39079, beanType)
		return nil
	}

	returnApiInfo_Param := make([]ApiInfo_Param, 0, 1)

	parseApiParamsStruct(beanValue, "", "", returnApiInfo_Param)
	return returnApiInfo_Param

}

func parseApiParamsStruct(beanValue reflect.Value, root string, fmtString string, apiInfo_Param []ApiInfo_Param) {
	if beanValue.Kind() == reflect.Ptr {
		beanValue = beanValue.Elem()
	}

	if beanValue.Kind() == reflect.Struct {
		beanValueType := beanValue.Type()
		for i := 0; i < beanValue.NumField(); i++ {
			name := getApiParamsAlias(beanValueType.Field(i), root)
			if beanValue.Field(i).Kind() == reflect.Slice { // 支持： 基础类型数组（[]string） , struct数组（[]struct）,struct指针数组[]*struct）
				parseApiParamsSlice(beanValue.Field(i), name, fmtString, apiInfo_Param)
			} else if beanValue.Field(i).Kind() == reflect.Map {
				parseApiParamsMap(beanValue.Field(i), name, fmtString, apiInfo_Param)
			} else if beanValue.Field(i).Kind() == reflect.Struct {
				parseApiParamsStruct(beanValue.Field(i), name, fmtString, apiInfo_Param)
			} else if beanValue.Field(i).Kind() == reflect.Ptr {
				parseApiParamsStruct(beanValue.Field(i), name, fmtString, apiInfo_Param)
			} else {
				parseApiParamsBaseType(beanValue.Field(i), name, fmtString, apiInfo_Param)
			}
		}
	}

}

func parseApiParamsMap(beanValue reflect.Value, root string, fmtString string, apiInfo_Param []ApiInfo_Param) {
	valueMapType := beanValue.Type().Elem() // map[xx]yy , yy的类型
	keyMaps := beanValue.MapKeys()
	for _, k1 := range keyMaps {
		if valueMapType.Kind() == reflect.Slice {
			parseApiParamsSlice(beanValue.MapIndex(k1), fmt.Sprintf("%s.%v", root, k1), fmtString, apiInfo_Param)
		} else if valueMapType.Kind() == reflect.Map {
			parseApiParamsMap(beanValue.MapIndex(k1), fmt.Sprintf("%s.%v", root, k1), fmtString, apiInfo_Param)
		} else if valueMapType.Kind() == reflect.Struct {
			parseApiParamsStruct(beanValue.MapIndex(k1), fmt.Sprintf("%s.%v", root, k1), fmtString, apiInfo_Param)
		} else if valueMapType.Kind() == reflect.Ptr {
			parseApiParamsStruct(beanValue.MapIndex(k1), fmt.Sprintf("%s.%v", root, k1), fmtString, apiInfo_Param)
		} else {
			parseApiParamsBaseType(beanValue.MapIndex(k1), fmt.Sprintf("%s.%v", root, k1), fmtString, apiInfo_Param)
		}
	}
}

func parseApiParamsSlice(beanValue reflect.Value, root string, fmtString string, apiInfo_Param []ApiInfo_Param) {
	for j := 0; j < beanValue.Len(); j++ {
		if beanValue.Index(j).Kind() == reflect.Struct {
			parseApiParamsStruct(beanValue.Index(j), fmt.Sprintf("%s.%d", root, j), fmtString, apiInfo_Param)
		} else if beanValue.Index(j).Kind() == reflect.Ptr {
			parseApiParamsStruct(beanValue.Index(j), fmt.Sprintf("%s.%d", root, j), fmtString, apiInfo_Param)
		} else if beanValue.Index(j).Kind() == reflect.Slice {
			parseApiParamsSlice(beanValue.Index(j), fmt.Sprintf("%s.%d", root, j), fmtString, apiInfo_Param)
		} else {
			parseApiParamsBaseType(beanValue.Index(j), root, fmtString, apiInfo_Param)
		}
	}
}

func parseApiParamsBaseType(beanValue reflect.Value, root string, fmtString string, apiInfo_Param []ApiInfo_Param) {
	if beanValue.Type().Kind() == reflect.String {
		if beanValue.String() == "" {
			return
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

func getApiParamsAlias(fieldType reflect.StructField, root string) string {
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

func getApiParamsDoc(fieldType reflect.StructField) string {
	doc := fieldType.Tag.Get("doc")
	return doc
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

*/
