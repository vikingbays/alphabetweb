// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sqler

import (
	"alphabet/env"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var METHOD_OF_FILTER_ROWS string = "FilterRows" // 每行数据渲染的方法，在result的struct中定义方法。

type DataFinder struct {
	tx             *sql.Tx
	db             *sql.DB
	cp             *ConnectionPool
	connectionName string
}

/*
获取一个连接，并返回DataFinder对象。

获取连接的方法是：
  dataFinder, err := sqler.GetConnectionForDataFinder("pg_test1")
  defer dataFinder.ReleaseConnection(0)
如果需要回滚调用：
  panic(&sqler.RollbackError{err})
或
  panic(&sqler.RollbackError{fmt.Errorf("error"))

@param connectionName  数据库连接池别名

@return 返回DataFinder对象
*/
func GetConnectionForDataFinder(connectionName string) (*DataFinder, error) {
	dataFinder := &DataFinder{}
	dataFinder.connectionName = connectionName
	dataFinder.cp = GetConnectionMap(connectionName)
	if dataFinder.cp == nil {
		log4go.ErrorLog(message.ERR_SQLR_39067, connectionName)
		return dataFinder, &ConnectionError{fmt.Errorf(message.ERR_SQLR_39067.String(), connectionName)}
	}

	dataFinder.db = dataFinder.cp.GetConnection()

	if dataFinder.db == nil {
		log4go.ErrorLog(message.ERR_SQLR_39068, connectionName)
		return dataFinder, &ConnectionError{fmt.Errorf(message.ERR_SQLR_39068.String(), connectionName)}
	}

	return dataFinder, nil

}

/*
完成操作后释放当前的连接

@param type 设置事务关闭条件，如果 0 ，表示抛出异常即可回滚关闭；如果 1 ，表示抛出异常必须是“RollbackError”的情况下才能回滚关闭。
例如：当执行 panic(&RollbackError("sql error.....")) ，就会在 ReleaseConnection中触发 rollback操作并释放连接。
*/
func (df *DataFinder) ReleaseConnection(typeOfError int) {
	if df.cp != nil && df.db != nil { //当前数据库连接正常。如果数据库连接异常，连接池不需要回收
		if df.tx != nil {
			successFlag := true
			if r := recover(); r != nil {
				if typeOfError == 0 {
					successFlag = false
					//log4go.ErrorLog("%v \n %s", r, env.GetInheritCodeInfoAlls())
					log4go.ErrorLog(r, env.GetInheritCodeInfoAlls("panic.go", true))
				} else {
					resultType := reflect.TypeOf(r)
					if resultType.Kind() == reflect.Ptr {
						if resultType.Elem().Name() == "RollbackError" {
							rollbackerror := r.(*RollbackError)
							//log4go.ErrorLog("%v \n %s", rollbackerror, env.GetInheritCodeInfoAlls())
							log4go.ErrorLog(rollbackerror)
							successFlag = false
						}
					} else {
						if resultType.Name() == "RollbackError" {
							rollbackerror := r.(RollbackError)
							//	log4go.ErrorLog("%v \n %s", rollbackerror, env.GetInheritCodeInfoAlls())
							log4go.ErrorLog(rollbackerror)
							successFlag = false
						}
					}
				}
			}

			if successFlag {
				df.tx.Commit()
			} else {
				df.tx.Rollback()
			}
		} else {
			if r := recover(); r != nil {
				log4go.ErrorLog(r)
			}
		}
		df.cp.ReleaseConnection(df.db)
	} else {
		if r := recover(); r != nil {
			log4go.ErrorLog(r)
		}
	}
	df.tx = nil
}

/*
开启一个数据库事务。事务使用完成后，需要关闭(关闭有两种可能，一种提交，一种回滚)。因此是配对出现的
df.StartTrans()
df.CommitTrans() 或  df.RollbackTrans()

*/
func (df *DataFinder) StartTrans() {
	tx, err := df.db.Begin()
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39069.String(), err.Error()))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	} else {
		df.tx = tx
	}
}

/*
提交一个数据库事务。事务必须先开启，才能提交，因此是配对出现的：
df.StartTrans()
df.CommitTrans()

*/
func (df *DataFinder) CommitTrans() {
	err := df.tx.Commit()
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39070.String(), err.Error()))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}
}

/*
回滚一个数据库事务。事务必须先开启，才能回滚，因此是配对出现的：
df.StartTrans()
df.RollbackTrans()

*/
func (df *DataFinder) RollbackTrans() {
	err := df.tx.Rollback()
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39071.String(), err.Error()))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}
}

/*
测试使用
*/
func (df *DataFinder) startTransErr() {
	panic(&RollbackError{errors.New(fmt.Sprintf("It is error when call df.StartTrans(), errinfo=%s", "err.Error()"))})
}

/*
测试使用
*/
func (df *DataFinder) commitTransErr() {
	panic(&RollbackError{errors.New(fmt.Sprintf("It is error when call df.CommitTrans(), errinfo=%s", "err.Error()"))})
}

/*
测试使用
*/
func (df *DataFinder) rollbackTransErr() {
	panic(&RollbackError{errors.New(fmt.Sprintf("It is error when call df.RollbackTrans(), errinfo=%s", "err.Error()"))})
}

/*
获取 sql.Tx 对象，用于执行原生sql，具体用法，参考如下例子：
  dataFinder, err := sqler.GetConnectionForDataFinder("pg_test1")
  defer dataFinder.ReleaseConnection(1)
  if err != nil {
    log4go.ErrorLog(err)
	panic(&sqler.RollbackError{err})
  } else {
	tx := dataFinder.GetTx()

	for i := 1; i <= count; i++ {
	  stmt, err := tx.Prepare("insert into user1(usrid,name,nanjing,money) values($1,$2,$3,$4 ) ")
	  if err != nil {
		log4go.ErrorLog(err)
		panic(&sqler.RollbackError{err})
	  } else {
		rs, _ := stmt.Exec(3000+i, "Ikkk_"+time.Now().String(), false, -99.6)
		affect, _ := rs.RowsAffected() //影响多少行记录
		log4go.DebugLog(affect)
	  }
	  log4go.DebugLog(">>>>>>>>>>>> insert single .... .... ")
    }
  }

@return (*sql.Tx) 返回sql.Tx 可以进行数据库操作

*/
func (df *DataFinder) GetTx() *sql.Tx {
	return df.tx
}

/*
查询一组对象存储到List中，该对象可以是struct 也可以是基础类型。例如：

  sql配置文件

  [[db]]
  name="getUserList"
  sql="""
      select * From user1 where usrid>#{minuSrid} and usrid<${maxusrID} and name like '${name}' and nanjing = #{nanjing}
      """

  查询代码：

  dataFinder, err := sqler.GetConnectionForDataFinder("pg_test1")
  defer dataFinder.ReleaseConnection(1)
  if err != nil {
	log4go.ErrorLog(err)
  } else {
	dataFinder.StartTrans()
	paramUser1 := new(User1)
	paramUser1.MinUsrid = minUsrId
	paramUser1.MaxUsrid = maxUsrId
	paramUser1.Name = name
	paramUser1.Nanjing = false
	resultList := make([]User1, 0, 1)
	dataFinder.DoQueryList("getUserList", *paramUser1, &resultList)
	for num, v := range resultList {
	  v1 := v.(*User1)
    if v1.Usrid > minUsrId && v1.Usrid < maxUsrId && v1.Nanjing == false {
		  log4go.DebugLog((num + 1), v)
	  }
	}
	dataFinder.CommitTrans()
}

如果需要对查询的结果数据每行数据进行渲染，即二次加工，只需要对结果对象的类型（struct）定义 FilterRows 方法。
例如：在每个结果对象中都会通过FilterRows方法：
    1）、可以渲染该行数据，（当前在Name属性前加上‘Apple_’） ；
	  2）、也可以过滤，是否不写入该行数据，返回是false 或者 int的方式，表示不加载 。
   type User1 struct {
	    Usrid    int
	    Name     string
	    Nanjing  bool
	    Money    float32
	    Hello    string
	    MinUsrid int
	    MaxUsrid int
   }

	 func (usr *User1) FilterRows() bool {
	 	usr.Name = "Apple_" + usr.Name
	 	if usr.Usrid%2 != 0 {
	 		return false
	 	} else {
	 		return true
	 	}
	 }

	说明： FilterRows 的返回值可以三种：
	 1）、无返回值，那么只渲染，不过滤
	 2）、返回值 bool ，支持过滤， 如果false，表示不写入该数据
	 2）、返回值 int ，支持过滤， 如果0，表示不写入该数据

   如果查询结果是指定的是基础类型，例如：string，int，float32，bool等，
   需要注意获取的结束数据类型对应关系是：
       string      ->   string
       int         ->   int
       int64       ->   int64
       float32     ->   float64
       float64     ->   float64
       bool        ->   bool
   其中数值类型都是转换成64位，如果需要强制转换
       int         ->   int64     :  int64(var1)
       int64       ->   int       :  int(var1)
       interface{} ->   int64     :  var1.(int64)


@param sqlname sql别名

@param param   参数对象

@param resultList  返回结果数据集

@return (error) 报错信息

*/

func (df *DataFinder) DoQueryList(sqlname string, param interface{}, resultList interface{}) error {
	resultType := reflect.TypeOf(resultList)
	sqlname = env.GetAppInfoFromRuntime(2) + "." + sqlname
	if resultType.Kind() == reflect.Ptr && resultType.Elem().Kind() == reflect.Slice {
		if resultType.Elem().Elem().Kind() == reflect.Int || resultType.Elem().Elem().Kind() == reflect.Int32 ||
			resultType.Elem().Elem().Kind() == reflect.Int64 || resultType.Elem().Elem().Kind() == reflect.String ||
			resultType.Elem().Elem().Kind() == reflect.Float32 || resultType.Elem().Elem().Kind() == reflect.Float64 ||
			resultType.Elem().Elem().Kind() == reflect.Bool {
			df.queryListForBaseType_nest(df.connectionName, sqlname, param, resultList)
		} else {
			df.queryListForStruct_nest(df.connectionName, sqlname, param, resultList)
		}
		return nil
	} else {
		return errors.New(message.ERR_SQLR_39072.String())
	}

}

func (df *DataFinder) queryListForBaseType_nest(appName string, sqlname string, param interface{}, resultList interface{}) {
	resultValue := reflect.ValueOf(resultList)
	resultType := reflect.TypeOf(resultList)
	resultFieldsType := resultType.Elem().Elem().Kind()

	if resultType.Kind() != reflect.Ptr {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39073.String(), appName, sqlname))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}

	if resultValue.Elem().CanSet() {
		resultValue.Elem().SetLen(0)
		resultValue.Elem().SetCap(1)
	}

	sqlConfig := GetSqlConfig()
	sqlstringTemplate := sqlConfig.SqlMap[appName].SqlMetaMap[sqlname].Sql //获取配置文件中定义的Sql模板信息

	sqlLink := new(SqlLink)
	sqlLink.BuildSqlLink(sqlstringTemplate)
	sqlstring, params := sqlLink.ConvertSql(param, df.cp.PreparedSqlStandard) // 通过sqlstringTemplate和param 生成可执行sql(变量sqlstring) 及 传参(变量params)

	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79006, appName, sqlname, sqlstringTemplate, sqlstring, params)
		log4go.DebugLog(message.DEG_SQLR_79007, "QueryList", time.Now().Format(env.FORMAT_TIME_YMDHMS_NS))
	}
	stmt, err := df.tx.Prepare(sqlstring)
	defer stmt.Close()
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}
	rows, err := stmt.Query(params...)
	defer rows.Close()
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}
	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79008, "QueryList", time.Now().Format(env.FORMAT_TIME_YMDHMS_NS))
	}
	columns, err := rows.Columns() // 获取查询结果字段名信息
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	} else if len(columns) > 1 {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), " length of result_column is not 1 .", appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {

		err := rows.Scan(scanArgs...)
		if err != nil {
			err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
			log4go.ErrorLog(err99)
			panic(&RollbackError{err99})

		}

		col := values[0]

		switch resultFieldsType {
		case reflect.String:
			if col != nil {
				resultValue.Elem().Set(reflect.Append(resultValue.Elem(), reflect.ValueOf(string(col))))
			}
		case reflect.Int:
			if col != nil {
				intValue, err := strconv.Atoi(string(col))
				if err != nil {
					err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "int", err, appName, sqlname, sqlstring, params))
					log4go.ErrorLog(err99)
					panic(&RollbackError{err99})
				} else {
					resultValue.Elem().Set(reflect.Append(resultValue.Elem(), reflect.ValueOf(intValue)))
				}
			}
		case reflect.Int32:
			if col != nil {
				intValue, err := strconv.Atoi(string(col))
				if err != nil {
					err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "int32", err, appName, sqlname, sqlstring, params))
					log4go.ErrorLog(err99)
					panic(&RollbackError{err99})
				} else {
					resultValue.Elem().Set(reflect.Append(resultValue.Elem(), reflect.ValueOf(int32(intValue))))
				}
			}
		case reflect.Int64:
			if col != nil {
				intValue, err := strconv.ParseInt(string(col), 10, 0)
				if err != nil {
					err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "int64", err, appName, sqlname, sqlstring, params))
					log4go.ErrorLog(err99)
					panic(&RollbackError{err99})
				} else {
					resultValue.Elem().Set(reflect.Append(resultValue.Elem(), reflect.ValueOf(intValue)))
				}
			}
		case reflect.Bool:
			if col != nil {
				boolValue, err := strconv.ParseBool(string(col))
				if err != nil {
					err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "bool", err, appName, sqlname, sqlstring, params))
					log4go.ErrorLog(err99)
					panic(&RollbackError{err99})
				} else {
					resultValue.Elem().Set(reflect.Append(resultValue.Elem(), reflect.ValueOf(boolValue)))
				}
			}
		case reflect.Float32:
			if col != nil {
				floatValue, err := strconv.ParseFloat(string(col), 32)
				if err != nil {
					err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "float32", err, appName, sqlname, sqlstring, params))
					log4go.ErrorLog(err99)
					panic(&RollbackError{err99})
				} else {
					resultValue.Elem().Set(reflect.Append(resultValue.Elem(), reflect.ValueOf(floatValue)))
				}
			}
		case reflect.Float64:
			if col != nil {
				floatValue, err := strconv.ParseFloat(string(col), 64)
				if err != nil {
					err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "float64", err, appName, sqlname, sqlstring, params))
					log4go.ErrorLog(err99)
					panic(&RollbackError{err99})
				} else {
					resultValue.Elem().Set(reflect.Append(resultValue.Elem(), reflect.ValueOf(floatValue)))
				}
			}
		}

	}

	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79008, "QueryList_data_comlete", time.Now().Format(env.FORMAT_TIME_YMDHMS_NS))
	}

}

func (df *DataFinder) queryListForStruct_nest(appName string, sqlname string, param interface{}, resultList interface{}) {

	resultValue := reflect.ValueOf(resultList)
	resultType := reflect.TypeOf(resultList)
	resultFieldsMap := make(map[string]string)
	resultFieldsTypeMap := make(map[string]reflect.Kind)

	if resultType.Kind() != reflect.Ptr {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39073.String(), appName, sqlname))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}
	if resultValue.Elem().CanSet() {
		resultValue.Elem().SetLen(0)
		resultValue.Elem().SetCap(1)
	}

	resultType_Struct := resultType.Elem().Elem()
	func1 := reflect.New(resultType_Struct).MethodByName(METHOD_OF_FILTER_ROWS)
	resultFilterRows := func1.IsValid()
	for i := 0; i < resultType_Struct.NumField(); i++ {
		lowerName := strings.ToLower(resultType_Struct.Field(i).Name)
		resultFieldsMap[lowerName] = resultType_Struct.Field(i).Name
		resultFieldsTypeMap[lowerName] = resultType_Struct.Field(i).Type.Kind()
	}
	sqlConfig := GetSqlConfig()
	sqlstringTemplate := sqlConfig.SqlMap[appName].SqlMetaMap[sqlname].Sql //获取配置文件中定义的Sql模板信息
	sqlLink := new(SqlLink)
	sqlLink.BuildSqlLink(sqlstringTemplate)
	sqlstring, params := sqlLink.ConvertSql(param, df.cp.PreparedSqlStandard) // 通过sqlstringTemplate和param 生成可执行sql(变量sqlstring) 及 传参(变量params)
	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79006, appName, sqlname, sqlstringTemplate, sqlstring, params)
		log4go.DebugLog(message.DEG_SQLR_79007, "QueryList", time.Now().Format(env.FORMAT_TIME_YMDHMS_NS))
	}
	stmt, err := df.tx.Prepare(sqlstring)
	defer stmt.Close()
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}
	rows, err := stmt.Query(params...)
	defer rows.Close()
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}
	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79008, "QueryList", time.Now().Format(env.FORMAT_TIME_YMDHMS_NS))
	}
	columns, err := rows.Columns() // 获取查询结果字段名信息
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})

	} else {
		// 转换成小写，用于匹配结果对象
		for i, column := range columns {
			columns[i] = strings.ToLower(column)
		}
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {

		err := rows.Scan(scanArgs...)
		if err != nil {
			err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
			log4go.ErrorLog(err99)
			panic(&RollbackError{err99})
		}

		rowStruct := reflect.New(resultType_Struct)
		mutableRowStruct := rowStruct.Elem()

		flag_returnValue_resultFilterRows := true // 如果是true表示这条记录有效，如果是false表示这条记录需要被过滤

		for i, col := range values {

			fieldResult := mutableRowStruct.FieldByName(resultFieldsMap[columns[i]])

			switch resultFieldsTypeMap[columns[i]] {
			case reflect.String:
				if col != nil {
					fieldResult.SetString(string(col))
				}
			case reflect.Int:
				if col != nil {
					intValue, err := strconv.ParseInt(string(col), 10, 0)
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "int", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})

					} else {
						fieldResult.SetInt(intValue)
					}
				}
			case reflect.Int32:
				if col != nil {
					intValue, err := strconv.ParseInt(string(col), 10, 0)
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "int32", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})
					} else {
						fieldResult.SetInt(intValue)
					}
				}
			case reflect.Int64:
				if col != nil {
					intValue, err := strconv.ParseInt(string(col), 10, 0)
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "int64", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})
					} else {
						fieldResult.SetInt(intValue)
					}
				}
			case reflect.Bool:
				if col != nil {
					boolValue, err := strconv.ParseBool(string(col))
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "bool", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})
					} else {
						fieldResult.SetBool(boolValue)
					}
				}
			case reflect.Float32:
				if col != nil {
					floatValue, err := strconv.ParseFloat(string(col), 32)
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "float32", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})
					} else {
						fieldResult.SetFloat(floatValue)
					}
				}
			case reflect.Float64:
				if col != nil {
					floatValue, err := strconv.ParseFloat(string(col), 64)
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "float64", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})
					} else {
						fieldResult.SetFloat(floatValue)
					}
				}
			}

		}

		if resultFilterRows { // 调用行过滤方法，处理这行数据
			returnValue_resultFilterRows := rowStruct.MethodByName(METHOD_OF_FILTER_ROWS).Call([]reflect.Value{})
			if len(returnValue_resultFilterRows) == 1 {
				if returnValue_resultFilterRows[0].Kind() == reflect.Bool {
					flag_returnValue_resultFilterRows = returnValue_resultFilterRows[0].Bool()
				} else if returnValue_resultFilterRows[0].Kind() == reflect.Int {
					flag_returnValue_resultFilterRows = returnValue_resultFilterRows[0].Int() > 0
				}
			}

		}
		if flag_returnValue_resultFilterRows {
			resultValue.Elem().Set(reflect.Append(resultValue.Elem(), rowStruct.Elem()))
		}
	}

	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79008, "QueryList_data_complete", time.Now().Format(env.FORMAT_TIME_YMDHMS_NS))
	}

}

/*
查询一组对象存储到Map中，存储信息是 key－value 模式的键值对 。其中：key 为某字段，value可以是某字段，也可以是该行纪录  ，例如：

  sql配置文件

  [[db]]
  name="getUserList"
  sql="""
      select * From user1 where usrid>#{minuSrid} and usrid<${maxusrID} and name like '${name}' and nanjing = #{nanjing}
      """

  查询代码：

  dataFinder, err := sqler.GetConnectionForDataFinder("pg_test1")
  defer dataFinder.ReleaseConnection(1)
  if err != nil {
	log4go.ErrorLog(err)
  } else {
	paramUser1 := new(User1)
	paramUser1.MinUsrid = minUsrId
	paramUser1.MaxUsrid = maxUsrId
	paramUser1.Name = name
	paramUser1.Nanjing = false
	dataFinder.StartTrans()
	userMap := make(map[string]User1)
	dataFinder.DoQueryMap("case1_query", paramUser, userMap, "Name", "")
	num := 0
	for k, v := range resultMap {
	  v1 := v.(*User1)
	  if v1.Usrid > minUsrId && v1.Usrid < maxUsrId && v1.Nanjing == false {
		num++
		log4go.DebugLog(num, k, v)
	  }
	}
	dataFinder.CommitTrans()
  }

其他解释请参考QueryList 说明

@param sqlname sql别名

@param param 参数

@param resultMap       返回结果集，必须是 map 结构

@param keycolumn    表示一个结果字段名，用于作为Map的Key

@param valuecolumn  表示一个结果字段名或者一行记录，用于作为Map的Value，如果valuecolumn为空串（""），表示是一行记录

@return (map[string]interface{})  返回结果对象是一个Map对象，必须要注意的是key必然被转换成string类型
*/
func (df *DataFinder) DoQueryMap(sqlname string,
	param interface{}, resultMap interface{},
	keycolumn string, valuecolumn string) {

	appName := df.connectionName

	sqlname = env.GetAppInfoFromRuntime(2) + "." + sqlname

	resultValue := reflect.ValueOf(resultMap)
	resultType := reflect.TypeOf(resultMap)

	if resultType.Kind() != reflect.Map {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39076.String(), appName, sqlname))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	} else if resultType.Elem().Kind() != reflect.Struct && valuecolumn == "" {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39077.String(), appName, sqlname))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}

	resultFieldsMap := make(map[string]string)
	resultFieldsTypeMap := make(map[string]reflect.Kind)

	var resultFilterRows bool = false
	// 判断如果struct，那么判断是否有这个方法
	if resultType.Elem().Kind() == reflect.Struct {
		func1 := reflect.New(resultType.Elem()).MethodByName(METHOD_OF_FILTER_ROWS)
		resultFilterRows = func1.IsValid()
	}

	keycolumn = strings.ToLower(keycolumn)
	valuecolumn = strings.ToLower(valuecolumn)

	keyIndex, valueIndex, valueStructIndex := -1, -1, -1

	if resultType.Elem().Kind() == reflect.Struct {
		for i := 0; i < resultType.Elem().NumField(); i++ {
			lowerName := strings.ToLower(resultType.Elem().Field(i).Name)
			resultFieldsMap[lowerName] = resultType.Elem().Field(i).Name
			resultFieldsTypeMap[lowerName] = resultType.Elem().Field(i).Type.Kind()
			if valuecolumn == lowerName {
				valueStructIndex = i
			}
		}
	} else { // 如果是非结构体，那么记录 keycolumn 和 valuecolumn 定义的信息，
		lowerkeyName := strings.ToLower(keycolumn)
		resultFieldsMap[lowerkeyName] = keycolumn
		resultFieldsTypeMap[lowerkeyName] = reflect.String
		lowerValueName := strings.ToLower(valuecolumn)
		resultFieldsMap[lowerValueName] = valuecolumn
		resultFieldsTypeMap[lowerValueName] = resultType.Elem().Kind()
	}
	//resultMap := make(map[string]interface{}) // 返回结果存储的List对象

	sqlConfig := GetSqlConfig()
	sqlstringTemplate := sqlConfig.SqlMap[appName].SqlMetaMap[sqlname].Sql //获取配置文件中定义的Sql模板信息

	sqlLink := new(SqlLink)
	sqlLink.BuildSqlLink(sqlstringTemplate)
	sqlstring, params := sqlLink.ConvertSql(param, df.cp.PreparedSqlStandard) // 通过sqlstringTemplate和param 生成可执行sql(变量sqlstring) 及 传参(变量params)
	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79006, appName, sqlname, sqlstringTemplate, sqlstring, params)
		log4go.DebugLog(message.DEG_SQLR_79007, "QueryMap", time.Now().Format(env.FORMAT_TIME_YMDHMS_NS))
	}
	stmt, err := df.tx.Prepare(sqlstring)
	defer stmt.Close()
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})

	}
	rows, err := stmt.Query(params...)
	defer rows.Close()
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}
	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79008, "QueryMap", time.Now().Format(env.FORMAT_TIME_YMDHMS_NS))
	}
	columns, err := rows.Columns() // 获取查询结果字段名信息
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	} else {
		// 转换成小写，用于匹配结果对象
		for i, column := range columns {
			columns[i] = strings.ToLower(column)
			if columns[i] == keycolumn {
				keyIndex = i
			}

			if columns[i] == valuecolumn {
				valueIndex = i
			}
		}
	}

	if keyIndex == -1 {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), " keycolumn [%s] is not found .", appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})

	}

	if valuecolumn != "" && valueIndex == -1 {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), "  valuecolumn [%s] is not found .", appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {

		err := rows.Scan(scanArgs...)
		if err != nil {
			err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
			log4go.ErrorLog(err99)
			panic(&RollbackError{err99})
		}

		var rowStruct reflect.Value
		if resultType.Elem().Kind() == reflect.Struct {
			rowStruct = reflect.New(resultType.Elem())
		} else { // 构建一个临时的Map
			rowStruct = reflect.New(resultType.Elem())
		}
		mutableRowStruct := rowStruct.Elem()
		var key1 string

		flag_returnValue_resultFilterRows := true // 如果是true表示这条记录有效，如果是false表示这条记录需要被过滤

		for i, col := range values {

			if col != nil {
				if keyIndex == i {
					key1 = string(col)
				}
			}

			var fieldResult reflect.Value
			if resultType.Elem().Kind() == reflect.Struct {
				fieldResult = mutableRowStruct.FieldByName(resultFieldsMap[columns[i]])
			} else {
				if valueIndex != i {
					continue
				}
				fieldResult = mutableRowStruct
			}

			switch resultFieldsTypeMap[columns[i]] {
			case reflect.String:
				if col != nil {
					fieldResult.SetString(string(col))
				}
			case reflect.Int:
				if col != nil {
					intValue, err := strconv.ParseInt(string(col), 10, 0)
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "int", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})
					} else {
						fieldResult.SetInt(intValue)
					}
				}
			case reflect.Int32:
				if col != nil {
					intValue, err := strconv.ParseInt(string(col), 10, 0)
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "int32", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})
					} else {
						fieldResult.SetInt(intValue)
					}
				}
			case reflect.Int64:
				if col != nil {
					intValue, err := strconv.ParseInt(string(col), 10, 0)
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "int64", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})
					} else {
						fieldResult.SetInt(intValue)
					}
				}
			case reflect.Bool:
				if col != nil {
					boolValue, err := strconv.ParseBool(string(col))
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "bool", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})

					} else {
						fieldResult.SetBool(boolValue)
					}
				}
			case reflect.Float32:
				if col != nil {
					floatValue, err := strconv.ParseFloat(string(col), 32)
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "float32", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})
					} else {
						fieldResult.SetFloat(floatValue)
					}
				}
			case reflect.Float64:
				if col != nil {
					floatValue, err := strconv.ParseFloat(string(col), 64)
					if err != nil {
						err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39075.String(), string(col), "float64", err, appName, sqlname, sqlstring, params))
						log4go.ErrorLog(err99)
						panic(&RollbackError{err99})

					} else {
						fieldResult.SetFloat(floatValue)
					}
				}
			}

		}

		if resultFilterRows { // 调用行过滤方法，处理这行数据
			rowStruct.MethodByName(METHOD_OF_FILTER_ROWS).Call([]reflect.Value{})
		}

		if resultFilterRows { // 调用行过滤方法，处理这行数据
			returnValue_resultFilterRows := rowStruct.MethodByName(METHOD_OF_FILTER_ROWS).Call([]reflect.Value{})
			if len(returnValue_resultFilterRows) == 1 {
				if returnValue_resultFilterRows[0].Kind() == reflect.Bool {
					flag_returnValue_resultFilterRows = returnValue_resultFilterRows[0].Bool()
				} else if returnValue_resultFilterRows[0].Kind() == reflect.Int {
					flag_returnValue_resultFilterRows = returnValue_resultFilterRows[0].Int() > 0
				}
			}

		}
		if flag_returnValue_resultFilterRows {
			if resultType.Elem().Kind() == reflect.Struct {
				if valuecolumn == "" {
					resultValue.SetMapIndex(reflect.ValueOf(key1), mutableRowStruct)
				} else {
					resultValue.SetMapIndex(reflect.ValueOf(key1), mutableRowStruct.Field(valueStructIndex))
				}
			} else {
				resultValue.SetMapIndex(reflect.ValueOf(key1), mutableRowStruct)
			}
		}

	}

	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79008, "QueryMap_data_complete", time.Now().Format(env.FORMAT_TIME_YMDHMS_NS))
	}
}

/*
执行所有DML的操作，包括：Create、Drop、Insert、Delete、Update等操作，例如：

  sql配置文件

  [[db]]
  name="insertUser1"
  sql="""
      insert into user1(usrid,name,nanjing,money) values(#{uSrid},'${name}',#{nanjing},#{money} )
      """

  查询代码：

  dataFinder, err := sqler.GetConnectionForDataFinder("pg_test1")
  defer dataFinder.ReleaseConnection(1)
  if err != nil {
	log4go.ErrorLog(err)
	panic(&sqler.RollbackError{err})
  } else {
	dataFinder.StartTrans()
	for i := 1; i <= count; i++ {
		paramUser1 := new(User1)
		paramUser1.Usrid = 100 + i
		paramUser1.Name = "Ikkk_" + time.Now().String()
		paramUser1.Nanjing = false
		paramUser1.Money = -99.6
		dataFinder.Exec( "insertUser1", *paramUser1)
	}
	dataFinder.CommitTrans()
  }

其他解释请参考QueryList 说明

@param sqlname sql别名

@param param   参数对象

@return int64 如果没有错误，那么就返回本次执行影响的行数。如果有错误直接抛出异常

*/
func (df *DataFinder) DoExec(sqlname string, param interface{}) int64 {
	appName := df.connectionName
	sqlname = env.GetAppInfoFromRuntime(2) + "." + sqlname
	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79009, appName, sqlname)
	}
	sqlConfig := GetSqlConfig()
	sqlstringTemplate := sqlConfig.SqlMap[appName].SqlMetaMap[sqlname].Sql //获取配置文件中定义的Sql模板信息
	sqlLink := new(SqlLink)
	sqlLink.BuildSqlLink(sqlstringTemplate)
	sqlstring, params := sqlLink.ConvertSql(param, df.cp.PreparedSqlStandard) // 通过sqlstringTemplate和param 生成可执行sql(变量sqlstring) 及 传参(变量params)
	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79006, appName, sqlname, sqlstringTemplate, sqlstring, params)
		log4go.DebugLog(message.DEG_SQLR_79007, "Execute SQL", time.Now().Format(env.FORMAT_TIME_YMDHMS_NS))
	}
	stmt, err := df.tx.Prepare(sqlstring)
	defer stmt.Close()
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		//panic(&RollbackError{err99})
		panic(err99.Error())
	}
	rs, err := stmt.Exec(params...)

	if log4go.IsDebugLevel() {
		log4go.DebugLog(message.DEG_SQLR_79008, "Execute SQL", time.Now().Format(env.FORMAT_TIME_YMDHMS_NS))
	}
	if err != nil {
		err99 := errors.New(fmt.Sprintf(message.ERR_SQLR_39074.String(), err, appName, sqlname, sqlstring, params))
		log4go.ErrorLog(err99)
		panic(&RollbackError{err99})

	} else {
		//id, _ := rs.LastInsertId()
		affect, _ := rs.RowsAffected() //影响多少行记录
		return affect
	}
}
