// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sqler

import (
	"alphabet/core/abtest"
	"alphabet/core/toml"
	"alphabet/log4go"
	"fmt"
	"testing"
)

var dbsconfig_sqler_df000 string = `
[[dbs]]
name="db_test1"
driverName="postgres"
dataSourceName="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&application_name=alphabet"
maxPoolSize=2
`

var sqlconfig_sqler_df000 string = `

[[db]]
name="case1_exist"
sql="""
  select 1 from pg_class where relname='case1_user1'::name and relkind='r'
    """

[[db]]
name="case1_drop"
sql="""
drop table case1_user1
    """

[[db]]
name="case1_create"
sql="""
  create table case1_user1(
    usrid integer,
    name character varying(100),
    nanjing boolean,
    money numeric(20,2),
    hello character varying(100)
  )
    """

[[db]]
name="case1_insert"
sql="""
  insert into case1_user1(usrid,name,nanjing,money,hello) values (#{usrid},#{name},#{nanjing},#{money},#{hello})
    """

[[db]]
name="case1_query"
sql="""
  select * from case1_user1 where usrid>#{minuSrid} and usrid<${maxusrID}
    """

[[db]]
name="case1_query_userid"
sql="""
  select UsrId from case1_user1 where usrid>#{minuSrid} and usrid<${maxusrID}
    """

`

type Case1_User1_df000 struct {
	Usrid    int
	Name     string
	Nanjing  bool
	Money    float64
	Hello    string
	MinUsrid int
	MaxUsrid int
}

var dbsConfigListToml_sqler_df000 Dbs_Config_List_Toml

var flag_create_init bool = true

func test_init_sqler_df000(t *testing.T) *DataFinder {

	test_init_sqler_df000_only_register(t)
	dataFinder, err := GetConnectionForDataFinder(dbsConfigListToml_sqler_df000.Dbs[0].Name)
	if err != nil {
		t.Error(err)
	}
	return dataFinder
}

func test_initSqlMap_sqler_df000(dbsConfigListToml Dbs_Config_List_Toml) {
	sqlConfig.SqlMap = make(map[string]SqlMetas)

	var dbSqlConfigListToml Db_SqlConfig_List_Toml
	if _, err := toml.Decode(sqlconfig_sqler_df000, &dbSqlConfigListToml); err != nil {
		log4go.ErrorLog("Reading sqlconfig is error . err: %s . \n %s", err.Error(), sqlconfig_sqler_df000)
	} else {
		//// fmt.Println(dbSqlConfigListToml)
		var sqlMetas SqlMetas
		dbs_config := dbsConfigListToml.Dbs[0]
		if len(sqlConfig.SqlMap[dbs_config.Name].SqlMetaMap) > 0 {

			sqlMetas = sqlConfig.SqlMap[dbs_config.Name]
		} else {
			sqlMetas.SqlMetaMap = make(map[string]SqlMeta)
		}
		for _, dbSqlConfig := range dbSqlConfigListToml.Db {
			sqlMeta := SqlMeta{dbSqlConfig.Name, useCurrentSql(dbSqlConfig, dbs_config.DriverName)}
			//sqlMetas.SqlMetaMap[appname+"."+dbSqlConfig.Name] = sqlMeta
			sqlMetas.SqlMetaMap["."+dbSqlConfig.Name] = sqlMeta
		}
		sqlConfig.SqlMap[dbs_config.Name] = sqlMetas
	}

}

func test_sqler_df000_case1_initData(dataFinder *DataFinder, createRowNum int, t *testing.T) {
	if !flag_create_init {
		return
	}

	count := createRowNum
	//paramUser := User1{MinUsrid: minUserId, MaxUsrid: maxUserId}
	paramUser := new(Case1_User1_df000)
	paramUser.MinUsrid = 0
	paramUser.MaxUsrid = 199

	dataFinder.StartTrans()

	var flag bool = false
	flags := make([]int, 0, 1)
	dataFinder.DoQueryList("case1_exist", nil, &flags)
	if len(flags) == 1 && flags[0] == 1 {
		flag = true
	}

	if flag {
		dataFinder.DoExec("case1_drop", nil)
	}
	dataFinder.DoExec("case1_create", nil)

	dataFinder.CommitTrans()

	{
		dataFinder.StartTrans()

		for i := 1; i <= count; i++ {
			paramInsertUser := &Case1_User1_df000{Usrid: i, Name: fmt.Sprintf("a_%d", i), Nanjing: true, Money: (66 + float64(i)), Hello: "hi"}
			dataFinder.DoExec("case1_insert", *paramInsertUser)
		}
		dataFinder.CommitTrans()

	}
}

func test_init_sqler_df000_only_register(t *testing.T) {

	abtest.InitTestEnv("")
	abtest.InitTestLog("")

	if _, err := toml.Decode(dbsconfig_sqler_df000, &dbsConfigListToml_sqler_df000); err != nil {
		t.Error(err)
		log4go.ErrorLog(err)
	} else {

		connectionMap = make(map[string]*ConnectionPool)
		sqlConfig.SqlMap = make(map[string]SqlMetas)

		initConnectionMap(dbsConfigListToml_sqler_df000)
		test_initSqlMap_sqler_df000(dbsConfigListToml_sqler_df000)
	}
}
