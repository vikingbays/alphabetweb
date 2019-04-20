// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sqler

import (
	//"alphabet/log4go"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	STRING_FRAGMET = iota
	VAR_OF_PREPARED
	VAR_OF_REPLACE
)

const (
	FLAG_POUND_SIGN  = "#{"
	FLAG_DOLLAR_SIGN = "${"
	FLAG_END_SIGN    = "}"
)

//记录解析后的sql模板信息，采用树状结构存储
//
//.Part： 存储片段信息
//
//.Type： `STRING_FRAGMET` 表示字符串类型。 `VAR_OF_PREPARED` 表示预处理变量，记录｀#{}｀定义的变量，当成｀?｀处理 。 `VAR_OF_REPLACE` 表示替换型变量，记录｀${}｀ 定义的变量，直接当字符串替换。
//
//.next: 下一个SqlLink
//
type SqlLink struct {
	Part string // 存储片段信息
	Type int
	next *SqlLink
}

/*
通过sqlstringTemplate和param 生成可执行sql(变量sqlstring) 及 传参(变量params)

@param param  参数对象

@param isPreparedSqlStandard  判断是否是标准sql，如果true，就采用｀?｀方式预处理，如果false，就采用｀$1｀方式预处理

@return sql  生成可执行sql

@return paramOut  生成预处理变量的参数数组
*/
func (sqlLink *SqlLink) ConvertSql(param interface{}, isPreparedSqlStandard bool) (sql string, paramOut []interface{}) {
	if param == nil {
		if sqlLink.Part != "" {
			sql = sqlLink.Part
		}
		nextSqlLink := sqlLink.next

		for nextSqlLink != nil {
			sql = sql + nextSqlLink.Part
			nextSqlLink = nextSqlLink.next
		}

	} else {
		paramValue := reflect.ValueOf(param)
		paramType := reflect.TypeOf(param)

		fieldsMap := make(map[string]string)

		if paramType.Kind() == reflect.Ptr {
			paramType = paramType.Elem()
			paramValue = paramValue.Elem()
		}
		for i := 0; i < paramType.NumField(); i++ {
			fieldsMap[strings.ToLower(paramType.Field(i).Name)] = paramType.Field(i).Name
		}

		varIndex := 1

		if sqlLink.Part != "" {
			if sqlLink.Type == STRING_FRAGMET {
				sql = sqlLink.Part
			} else if sqlLink.Type == VAR_OF_PREPARED {
				sql = sqlLink.Part
			} else if sqlLink.Type == VAR_OF_REPLACE {
				sql = sqlLink.Part
			}
		}
		nextSqlLink := sqlLink.next

		for nextSqlLink != nil {

			if nextSqlLink.Type == STRING_FRAGMET {
				sql = sql + nextSqlLink.Part
			} else if nextSqlLink.Type == VAR_OF_PREPARED {
				if isPreparedSqlStandard {
					sql = sql + "?"
				} else {
					sql = sql + "$" + strconv.Itoa(varIndex)
				}

				varIndex++

				if fieldsMap[nextSqlLink.Part] != "" {
					value1 := paramValue.FieldByName(fieldsMap[nextSqlLink.Part])
					switch value1.Kind() {
					case reflect.String:
						paramOut = append(paramOut, value1.String())
					case reflect.Int:
						paramOut = append(paramOut, strconv.FormatInt(value1.Int(), 10))
					case reflect.Int32:
						paramOut = append(paramOut, strconv.FormatInt(value1.Int(), 10))
					case reflect.Int64:
						paramOut = append(paramOut, strconv.FormatInt(value1.Int(), 10))
					case reflect.Bool:
						paramOut = append(paramOut, strconv.FormatBool(value1.Bool()))
					case reflect.Float32:
						paramOut = append(paramOut, strconv.FormatFloat(value1.Float(), 'f', -1, 32))
					case reflect.Float64:
						paramOut = append(paramOut, strconv.FormatFloat(value1.Float(), 'f', -1, 64))
					}
				} else {
					// 在params中找不到可对应 nextSqlLink.Part 的参数
					// 没有参数，应该抛出异常未处理
					// todo
				}

			} else if nextSqlLink.Type == VAR_OF_REPLACE {

				if fieldsMap[nextSqlLink.Part] != "" {

					value1 := paramValue.FieldByName(fieldsMap[nextSqlLink.Part])

					switch value1.Kind() {
					case reflect.String:
						sql = sql + value1.String()
					case reflect.Int:
						sql = sql + strconv.FormatInt(value1.Int(), 10)
					case reflect.Int32:
						sql = sql + strconv.FormatInt(value1.Int(), 10)
					case reflect.Int64:
						sql = sql + strconv.FormatInt(value1.Int(), 10)
					case reflect.Bool:
						sql = sql + strconv.FormatBool(value1.Bool())
					case reflect.Float32:
						sql = sql + strconv.FormatFloat(value1.Float(), 'f', -1, 32)
					case reflect.Float64:
						sql = sql + strconv.FormatFloat(value1.Float(), 'f', -1, 64)
					}

				} else {
					// 在params中找不到可对应 nextSqlLink.Part 的参数
					// 没有参数，应该抛出异常未处理
					// todo
				}

			}

			nextSqlLink = nextSqlLink.next
		}
	}
	return
}

func (sqlLink *SqlLink) String() (sql string) {
	sql = "SqlLink toString : \n"
	if sqlLink.Part != "" {
		sql = sql + fmt.Sprintf(" Part : %s , Type : %d \n", sqlLink.Part, sqlLink.Type)
	}
	nextSqlLink := sqlLink.next
	for nextSqlLink != nil {
		sql = sql + fmt.Sprintf(" Type : %d , Part : %s  \n", nextSqlLink.Type, nextSqlLink.Part)
		nextSqlLink = nextSqlLink.next
	}
	return
}

/*
根据配置的sql模板信息，解析出变量信息，并生成SqlLink的结构。该方法才用递归模式，逐级解析。

@param sqlpart  sql模板信息
*/
func (sqlLink *SqlLink) BuildSqlLink(sqlpart string) {
	child, nextSqlpart := sqlLink.nextSqlPart(sqlLink, sqlpart)
	for nextSqlpart != "" {
		child, nextSqlpart = sqlLink.nextSqlPart(child, nextSqlpart)
	}
}

func (sqlLink *SqlLink) nextSqlPart(parent *SqlLink, sqlpart string) (child *SqlLink, nextSqlpart string) {
	poundSign := strings.Index(sqlpart, FLAG_POUND_SIGN)
	dollarSign := strings.Index(sqlpart, FLAG_DOLLAR_SIGN)
	if poundSign == -1 && dollarSign == -1 { // 两个变量都不存在
		childSqlLink := new(SqlLink)
		childSqlLink.Part = sqlpart
		childSqlLink.Type = STRING_FRAGMET
		childSqlLink.next = nil
		child = childSqlLink
		nextSqlpart = ""
		parent.next = childSqlLink
	} else if poundSign == -1 { // 只有"${" 变量，按照  "${" 变量处理
		child, nextSqlpart = sqlLink.getVarFromSqlPart(parent, sqlpart, FLAG_DOLLAR_SIGN)
	} else if dollarSign == -1 { // 只有"#{" 变量，按照 "#{" 变量处理
		child, nextSqlpart = sqlLink.getVarFromSqlPart(parent, sqlpart, FLAG_POUND_SIGN)
	} else if poundSign < dollarSign { //按照 "#{" 变量处理
		child, nextSqlpart = sqlLink.getVarFromSqlPart(parent, sqlpart, FLAG_POUND_SIGN)
	} else if poundSign > dollarSign { //按照 "${" 变量处理
		child, nextSqlpart = sqlLink.getVarFromSqlPart(parent, sqlpart, FLAG_DOLLAR_SIGN)
	}
	return
}

func (sqlLink *SqlLink) getVarFromSqlPart(parent *SqlLink, sqlpart string, sign string) (child *SqlLink, nextSqlpart string) {
	childLeftSqlLink := new(SqlLink)
	childRightSqlLink := new(SqlLink)
	sqlpartArray := strings.SplitN(sqlpart, sign, 2)
	childLeftSqlLink.Part = sqlpartArray[0]
	childLeftSqlLink.Type = STRING_FRAGMET
	childLeftSqlLink.next = childRightSqlLink
	parent.next = childLeftSqlLink
	subSqlpartArray := strings.SplitN(sqlpartArray[1], FLAG_END_SIGN, 2)
	if len(subSqlpartArray) == 2 {

		hasSpace := strings.Contains(subSqlpartArray[0], " ")
		if !hasSpace { // 判断没有空串，说明是变量
			childRightSqlLink.Part = strings.ToLower(subSqlpartArray[0])
			childRightSqlLink.Type = sqlLink.getSignVar(sign)
		} else {
			childRightSqlLink.Part = sign + subSqlpartArray[0] + FLAG_END_SIGN
			childRightSqlLink.Type = STRING_FRAGMET
		}
		nextSqlpart = subSqlpartArray[1]
	} else if len(subSqlpartArray) == 1 {
		childRightSqlLink.Part = sign + subSqlpartArray[0]
		childRightSqlLink.Type = STRING_FRAGMET
		childRightSqlLink.next = nil
		nextSqlpart = ""
	}
	child = childRightSqlLink
	return
}

func (sqlLink *SqlLink) getSignVar(sign string) int {
	if sign == FLAG_POUND_SIGN {
		return VAR_OF_PREPARED
	} else {
		return VAR_OF_REPLACE
	}
}
