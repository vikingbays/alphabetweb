// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package webutils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func DecodeHttpParams(paramStringsMap map[string][]string, beanValue reflect.Value) {
	decodeHttpParamsStruct(paramStringsMap, "", beanValue)
}

/*
func DecodeHttpParams(paramStringsMap map[string][]string, bean interface{}) {
	beanValue := reflect.ValueOf(bean)
	decodeHttpParamsStruct(paramStringsMap, "", beanValue)
}
*/

func decodeHttpParamsStruct(paramStringsMap map[string][]string, root string, beanValue reflect.Value) {
	if beanValue.Kind() == reflect.Ptr {
		beanValue = beanValue.Elem()
	}
	if beanValue.Kind() == reflect.Struct {
		beanValueType := beanValue.Type()
		for i := 0; i < beanValue.NumField(); i++ {
			name := getHttpParamAlias(beanValueType.Field(i), root)
			if beanValue.Field(i).Kind() == reflect.Slice { // 支持： 基础类型数组（[]string） , struct数组（[]struct）,struct指针数组[]*struct）
				decodeHttpParamsSlice(paramStringsMap, name, beanValue.Field(i))
			} else if beanValue.Field(i).Kind() == reflect.Map {
				decodeHttpParamsMap(paramStringsMap, name, beanValue.Field(i))
			} else if beanValue.Field(i).Kind() == reflect.Struct {
				obj := reflect.New(beanValue.Field(i).Type())
				obj = obj.Elem()
				decodeHttpParamsStruct(paramStringsMap, name, obj)
				beanValue.Field(i).Set(obj)
			} else if beanValue.Field(i).Kind() == reflect.Ptr {
				subTypeElem := beanValue.Field(i).Type().Elem()
				if subTypeElem.Kind() == reflect.Struct {
					obj := reflect.New(beanValue.Field(i).Type().Elem())
					decodeHttpParamsStruct(paramStringsMap, name, obj)
					beanValue.Field(i).Set(obj)
				}
			} else if beanValue.Field(i).Kind() == reflect.Int ||
				beanValue.Field(i).Kind() == reflect.Int32 ||
				beanValue.Field(i).Kind() == reflect.Int64 {
				decodeHttpParamsBaseType_Int(paramStringsMap, name, beanValue.Field(i))
			} else if beanValue.Field(i).Kind() == reflect.Float32 ||
				beanValue.Field(i).Kind() == reflect.Float64 {
				decodeHttpParamsBaseType_Float(paramStringsMap, name, beanValue.Field(i))
			} else if beanValue.Field(i).Kind() == reflect.Bool {
				decodeHttpParamsBaseType_Bool(paramStringsMap, name, beanValue.Field(i))
			} else if beanValue.Field(i).Kind() == reflect.String {
				decodeHttpParamsBaseType_String(paramStringsMap, name, beanValue.Field(i))
			}

		}
	}

}

func decodeHttpParamsMap(paramStringsMap map[string][]string, root string, beanValue reflect.Value) reflect.Value {
	valueMapType := beanValue.Type().Elem() // map[xx]yy , yy的类型
	keyMapType := beanValue.Type().Key()

	if valueMapType.Kind() == reflect.Struct || valueMapType.Kind() == reflect.Ptr {
		beanMap1 := reflect.MakeMap(beanValue.Type())
		keyPrefix := fmt.Sprintf("%s.", root)
		structNameMap := make(map[string]int)
		for k, _ := range paramStringsMap {

			if strings.HasPrefix(k, strings.ToLower(keyPrefix)) {
				strArr := strings.Split(k[len(keyPrefix):], ".")
				structNameMap[strArr[0]] = 1
			}
		}

		for k, _ := range structNameMap {
			k1 := reflect.New(keyMapType)
			k1 = k1.Elem()
			flag1 := decodeHttpParamsBaseType_nest(k, keyMapType.Kind(), k1)

			var obj reflect.Value
			if valueMapType.Kind() == reflect.Struct {
				obj = reflect.New(valueMapType)
				obj = obj.Elem()
				decodeHttpParamsStruct(paramStringsMap, keyPrefix+k, obj)
			} else {
				obj = reflect.New(valueMapType.Elem())
				decodeHttpParamsStruct(paramStringsMap, keyPrefix+k, obj.Elem())
			}

			if flag1 {
				beanMap1.SetMapIndex(k1, obj)
			}

		}

		if beanValue.CanSet() {
			beanValue.Set(beanMap1)
		} else {
			return beanMap1
		}

	} else if valueMapType.Kind() == reflect.Slice {
		beanMap1 := reflect.MakeMap(beanValue.Type())
		keyPrefix := fmt.Sprintf("%s.", root)
		for k, _ := range paramStringsMap {
			if strings.HasPrefix(k, strings.ToLower(keyPrefix)) {
				k1 := reflect.New(keyMapType)
				k1 = k1.Elem()
				strArr := strings.Split(k[len(keyPrefix):], ".")
				flag1 := decodeHttpParamsBaseType_nest(strArr[0], keyMapType.Kind(), k1)
				beanValue1 := reflect.MakeSlice(valueMapType, 0, 1)
				beanValue1 = decodeHttpParamsSlice(paramStringsMap, k, beanValue1)
				if flag1 {
					beanMap1.SetMapIndex(k1, beanValue1)
				}
			}
		}
		if beanValue.CanSet() {
			beanValue.Set(beanMap1)
		} else {
			return beanMap1
		}

	} else if valueMapType.Kind() == reflect.Map {
		beanMap1 := reflect.MakeMap(beanValue.Type())
		keyPrefix := fmt.Sprintf("%s.", root)
		for k, _ := range paramStringsMap {
			if strings.HasPrefix(k, strings.ToLower(keyPrefix)) {
				k1 := reflect.New(keyMapType)
				k1 = k1.Elem()
				strArr := strings.Split(k[len(keyPrefix):], ".")
				flag1 := decodeHttpParamsBaseType_nest(strArr[0], keyMapType.Kind(), k1)
				beanValue1 := reflect.MakeMap(valueMapType)
				beanValue1 = decodeHttpParamsMap(paramStringsMap, keyPrefix+strArr[0], beanValue1)
				if flag1 {
					beanMap1.SetMapIndex(k1, beanValue1)
				}
			}
		}
		if beanValue.CanSet() {
			beanValue.Set(beanMap1)
		} else {
			return beanMap1
		}

	} else {
		beanMap1 := reflect.MakeMap(beanValue.Type())
		keyPrefix := fmt.Sprintf("%s.", root)
		for k, v := range paramStringsMap {
			if strings.HasPrefix((k), strings.ToLower(keyPrefix)) {
				strArr := strings.Split(k[len(keyPrefix):], ".")
				//strArr := strings.Split(k, ".")
				if len(strArr) == 1 && len(v) == 1 {
					k1 := reflect.New(keyMapType)
					k1 = k1.Elem()
					flag1 := decodeHttpParamsBaseType_nest(strArr[0], keyMapType.Kind(), k1)
					v1 := reflect.New(valueMapType)
					v1 = v1.Elem()
					flag2 := decodeHttpParamsBaseType_nest(v[0], valueMapType.Kind(), v1)
					if flag1 == true && flag2 == true {
						beanMap1.SetMapIndex(k1, v1)
					}
				}
			}
		}
		if beanValue.CanSet() {
			beanValue.Set(beanMap1)
		} else {
			return beanMap1
		}
	}
	return beanValue
}

func decodeHttpParamsSlice(paramStringsMap map[string][]string, root string, beanValue reflect.Value) reflect.Value {

	if beanValue.Type().Elem().Kind() == reflect.Slice {
		beanValue1 := reflect.MakeSlice(beanValue.Type(), 0, 1)
		for i := 0; ; i++ {
			key := fmt.Sprintf("%s.%d", root, i)
			value := paramStringsMap[strings.ToLower(key)]
			if value != nil && len(value) > 0 {
				beanValue1_sub := reflect.MakeSlice(beanValue.Type().Elem(), 0, 1)
				beanValue1_sub = decodeHttpParamsSlice(paramStringsMap, key, beanValue1_sub)
				beanValue1 = reflect.Append(beanValue1, beanValue1_sub)
			} else {
				break
			}
		}
		if beanValue.CanSet() {
			beanValue.Set(beanValue1)
		} else {
			return beanValue1
		}
	} else if beanValue.Type().Elem().Kind() == reflect.Struct ||
		beanValue.Type().Elem().Kind() == reflect.Ptr {
		beanValue1 := reflect.MakeSlice(beanValue.Type(), 0, 1)
		keyPrefix := fmt.Sprintf("%s.", root)
		count := 0
		structArrayToParamStringsMap := make(map[int]map[string][]string)
		for k, v := range paramStringsMap {
			if strings.HasPrefix((k), strings.ToLower(keyPrefix)) {
				strArr := strings.Split(k, ".")
				if len(strArr) > 2 {
					n, errInt := strconv.Atoi(strArr[1])
					if errInt == nil {
						if structArrayToParamStringsMap[n] == nil {
							structArrayToParamStringsMap[n] = make(map[string][]string)
						}

						structArrayToParamStringsMap[n][strings.Join(strArr[2:], ".")] = v
						if count < (n + 1) {
							count = n + 1
						}
					}
				}
			}
		}

		for i := 0; i < count; i++ {
			if structArrayToParamStringsMap[i] != nil {
				if beanValue.Type().Elem().Kind() == reflect.Ptr {
					subTypeElem := beanValue.Type().Elem().Elem()
					if subTypeElem.Kind() == reflect.Struct {
						obj := reflect.New(beanValue.Type().Elem().Elem())
						decodeHttpParamsStruct(structArrayToParamStringsMap[i], "", obj)
						beanValue1 = reflect.Append(beanValue1, obj)
					}
				} else {
					obj := reflect.New(beanValue.Type().Elem())
					obj = obj.Elem()
					decodeHttpParamsStruct(structArrayToParamStringsMap[i], "", obj)
					beanValue1 = reflect.Append(beanValue1, obj)
				}
			}
		}

		if beanValue.CanSet() {
			beanValue.Set(beanValue1)
		} else {
			return beanValue1
		}

	} else {
		value := paramStringsMap[strings.ToLower(root)]
		if value != nil {
			beanValue1 := reflect.MakeSlice(beanValue.Type(), 0, 1)

			for _, v := range value {
				if beanValue.Type().Elem().Kind() == reflect.Int {
					n, err := strconv.Atoi(v)
					if err == nil {
						beanValue1 = reflect.Append(beanValue1, reflect.ValueOf(n))
					}
				} else if beanValue.Type().Elem().Kind() == reflect.Int32 {
					n, err := strconv.ParseInt(v, 10, 32)
					if err == nil {
						beanValue1 = reflect.Append(beanValue1, reflect.ValueOf(int32(n)))
					}
				} else if beanValue.Type().Elem().Kind() == reflect.Int64 {
					n, err := strconv.ParseInt(v, 10, 64)
					if err == nil {
						beanValue1 = reflect.Append(beanValue1, reflect.ValueOf((n)))
					}
				} else if beanValue.Type().Elem().Kind() == reflect.Float32 {
					n, err := strconv.ParseFloat(v, 32)
					if err == nil {
						beanValue1 = reflect.Append(beanValue1, reflect.ValueOf(float32(n)))
					}
				} else if beanValue.Type().Elem().Kind() == reflect.Float64 {
					n, err := strconv.ParseFloat(v, 64)
					if err == nil {
						beanValue1 = reflect.Append(beanValue1, reflect.ValueOf((n)))
					}
				} else if beanValue.Type().Elem().Kind() == reflect.Bool {
					n, err := strconv.ParseBool(v)
					if err == nil {
						beanValue1 = reflect.Append(beanValue1, reflect.ValueOf((n)))
					}
				} else if beanValue.Type().Elem().Kind() == reflect.String {
					beanValue1 = reflect.Append(beanValue1, reflect.ValueOf((v)))
				}
			}
			if beanValue1.Len() > 0 {
				if beanValue.CanSet() {
					beanValue.Set(beanValue1)
				} else {
					return beanValue1
				}

			}
		}
	}
	return beanValue

}

func decodeHttpParamsBaseType_Int(paramStringsMap map[string][]string, root string, beanValue reflect.Value) {
	values := paramStringsMap[strings.ToLower(root)]
	if values != nil && len(values) > 0 {
		decodeHttpParamsBaseType_Int_nest(values[0], beanValue)
	}
}

func decodeHttpParamsBaseType_Int_nest(value string, beanValue reflect.Value) bool {
	n, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		beanValue.SetInt(n)
		return true
	} else {
		return false
	}
}

func decodeHttpParamsBaseType_Float(paramStringsMap map[string][]string, root string, beanValue reflect.Value) {
	values := paramStringsMap[strings.ToLower(root)]
	if values != nil && len(values) > 0 {
		decodeHttpParamsBaseType_Float_nest(values[0], beanValue)
	}
}

func decodeHttpParamsBaseType_Float_nest(value string, beanValue reflect.Value) bool {
	f, err := strconv.ParseFloat(value, 64)
	if err == nil {
		beanValue.SetFloat(f)
		return true
	} else {
		return false
	}
}

func decodeHttpParamsBaseType_Bool(paramStringsMap map[string][]string, root string, beanValue reflect.Value) {
	values := paramStringsMap[strings.ToLower(root)]
	if values != nil && len(values) > 0 {
		decodeHttpParamsBaseType_Bool_nest(values[0], beanValue)
	}
}

func decodeHttpParamsBaseType_Bool_nest(value string, beanValue reflect.Value) bool {
	b, err := strconv.ParseBool(value)
	if err == nil {
		beanValue.SetBool(b)
		return true
	} else {
		return false
	}
}

func decodeHttpParamsBaseType_String(paramStringsMap map[string][]string, root string, beanValue reflect.Value) {
	values := paramStringsMap[strings.ToLower(root)]
	if values != nil && len(values) > 0 {
		decodeHttpParamsBaseType_String_nest(values[0], beanValue)
	}
}

func decodeHttpParamsBaseType_String_nest(value string, beanValue reflect.Value) bool {
	beanValue.SetString(value)
	return true
}

func decodeHttpParamsBaseType_nest(value string, kind reflect.Kind, beanValue reflect.Value) bool {
	flag := false
	if kind == reflect.Int || kind == reflect.Int32 || kind == reflect.Int64 {
		flag = decodeHttpParamsBaseType_Int_nest(value, beanValue)
	} else if kind == reflect.Float32 || kind == reflect.Float64 {
		flag = decodeHttpParamsBaseType_Float_nest(value, beanValue)
	} else if kind == reflect.Bool {
		flag = decodeHttpParamsBaseType_Bool_nest(value, beanValue)
	} else if kind == reflect.String {
		flag = decodeHttpParamsBaseType_String_nest(value, beanValue)
	}
	return flag
}

/*
func decodeHttpParamsStruct1(paramStringsMap map[string][]string, root string, paramBean interface{}) {
	data := reflect.ValueOf(paramBean)
	if data.Kind() == reflect.Ptr {
		data = data.Elem()
	}

	if data.Kind() == reflect.Struct {
		dataType := data.Type()
		for i := 0; i < data.NumField(); i++ {
			if data.Field(i).Kind() == reflect.Slice {
				if data.Field(i).Type().Elem().Kind() == reflect.Int {
					name := getHttpParamAlias(dataType.Field(i), root)
					value := paramStringsMap[name]
					result := make([]int, 0, 1)
					if value != nil && len(value) > 0 {
						for _, v := range value {
							n, err := strconv.Atoi(v)
							if err == nil {
								result = append(result, n)
							}
						}
					}
					if len(result) > 0 {
						data.Field(i).Set(reflect.ValueOf(result))
					}
				} else if data.Field(i).Type().Elem().Kind() == reflect.Int32 {
					name := getHttpParamAlias(dataType.Field(i), root)
					value := paramStringsMap[name]
					result := make([]int32, 0, 1)
					if value != nil && len(value) > 0 {
						for _, v := range value {
							n, err := strconv.Atoi(v)
							if err == nil {
								result = append(result, int32(n))
							}
						}
					}
					if len(result) > 0 {
						data.Field(i).Set(reflect.ValueOf(result))
					}
				} else if data.Field(i).Type().Elem().Kind() == reflect.Int64 {
					name := getHttpParamAlias(dataType.Field(i), root)
					value := paramStringsMap[name]
					result := make([]int64, 0, 1)
					if value != nil && len(value) > 0 {
						for _, v := range value {
							n, err := strconv.ParseInt(v, 10, 64)
							if err == nil {
								result = append(result, (n))
							}
						}
					}
					if len(result) > 0 {
						data.Field(i).Set(reflect.ValueOf(result))
					}
				} else if data.Field(i).Type().Elem().Kind() == reflect.Float32 {
					name := getHttpParamAlias(dataType.Field(i), root)
					value := paramStringsMap[name]
					result := make([]float32, 0, 1)
					if value != nil && len(value) > 0 {
						for _, v := range value {
							f, err := strconv.ParseFloat(v, 32)
							if err == nil {
								result = append(result, float32(f))
							}
						}
					}
					if len(result) > 0 {
						data.Field(i).Set(reflect.ValueOf(result))
					}
				} else if data.Field(i).Type().Elem().Kind() == reflect.Float64 {
					name := getHttpParamAlias(dataType.Field(i), root)
					value := paramStringsMap[name]
					result := make([]float64, 0, 1)
					if value != nil && len(value) > 0 {
						for _, v := range value {
							f, err := strconv.ParseFloat(v, 64)
							if err == nil {
								result = append(result, f)
							}
						}
					}
					if len(result) > 0 {
						data.Field(i).Set(reflect.ValueOf(result))
					}
				} else if data.Field(i).Type().Elem().Kind() == reflect.Bool {
					name := getHttpParamAlias(dataType.Field(i), root)
					value := paramStringsMap[name]
					result := make([]bool, 0, 1)
					if value != nil && len(value) > 0 {
						for _, v := range value {
							b, err := strconv.ParseBool(v)
							if err == nil {
								result = append(result, b)
							}
						}
					}
					if len(result) > 0 {
						data.Field(i).Set(reflect.ValueOf(result))
					}
				} else if data.Field(i).Type().Elem().Kind() == reflect.String {
					name := getHttpParamAlias(dataType.Field(i), root)
					value := paramStringsMap[name]
					if value != nil && len(value) > 0 {
						data.Field(i).Set(reflect.ValueOf(value))
					}
				}
			} else if data.Field(i).Kind() == reflect.Struct {
			} else if data.Field(i).Kind() == reflect.Ptr {

			} else if data.Field(i).Kind() == reflect.Bool {
				name := getHttpParamAlias(dataType.Field(i), root)
				value := paramStringsMap[name]
				if value != nil && len(value) > 0 {
					b, err := strconv.ParseBool(value[0])
					if err == nil {
						data.Field(i).SetBool(b)
					}
				}
			} else if data.Field(i).Kind() == reflect.Float32 || data.Field(i).Kind() == reflect.Float64 {
				name := getHttpParamAlias(dataType.Field(i), root)
				value := paramStringsMap[name]
				if value != nil && len(value) > 0 {
					f, err := strconv.ParseFloat(value[0], 64)
					if err == nil {
						data.Field(i).SetFloat(f)
					}
				}
			} else if data.Field(i).Kind() == reflect.Int ||
				data.Field(i).Kind() == reflect.Int32 ||
				data.Field(i).Kind() == reflect.Int64 {
				name := getHttpParamAlias(dataType.Field(i), root)
				value := paramStringsMap[name]
				if value != nil && len(value) > 0 {
					n, err := strconv.ParseInt(value[0], 10, 64)
					if err == nil {
						data.Field(i).SetInt(n)
					}
				}
			} else if data.Field(i).Kind() == reflect.String {
				name := getHttpParamAlias(dataType.Field(i), root)
				value := paramStringsMap[name]
				if value != nil && len(value) > 0 {
					data.Field(i).SetString(value[0])
				}
			}
		}
	}
}

*/
