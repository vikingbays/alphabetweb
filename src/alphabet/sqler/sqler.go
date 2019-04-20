// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sqler

import (
	"alphabet/core/toml"
	"alphabet/core/utils"
	"alphabet/env"
	"alphabet/log4go"
	"io/ioutil"
)

type Dbs_Config_List_Toml struct {
	Dbs []Dbs_Config_Toml
}

type Dbs_Config_Toml struct {
	Name           string
	DriverName     string
	DataSourceName string
	MaxPoolSize    int
}

type Db_SqlConfig_List_Toml struct {
	Db []Db_SqlConfig_Toml
}

type Db_SqlConfig_Toml struct {
	Name         string
	Sql          string // 通用sql语句
	Sql_postgres string // 针对postgresql数据库的sql语句，特殊处理
	Sql_mysql    string // 针对mysql     数据库的sql语句，特殊处理
	Sql_sqlite3  string // 针对sqlite3   数据库的sql语句，特殊处理
}

var connectionMap map[string]*ConnectionPool

var sqlConfig SqlConfig

// 初始化数据库连接，sql注册等
func Init() {
	connectionMap = make(map[string]*ConnectionPool)
	sqlConfig.SqlMap = make(map[string]SqlMetas)

	for _, appName := range env.Env_Project_Resource_Apps {
		var dbsConfigListToml Dbs_Config_List_Toml
		pathconfig := env.Env_Project_Resource + "/" + appName + "/" + env.GetSysconfigPath() + env.Env_Project_Resource_DbConfig
		if _, err := toml.DecodeFile(pathconfig, &dbsConfigListToml); err != nil {
			log4go.ErrorLog(err)
		} else {
			initConnectionMap(dbsConfigListToml)
			initSqlMap(dbsConfigListToml, appName)
		}
	}

}

/*
 根据Dbs_Config_List_Toml配置信息，生成连接池
 Dbs_Config_List_Toml 对应的配置文件是 apps/dbsconfig.toml ，可以配置多个数据库连接
*/
func initConnectionMap(dbsConfigListToml Dbs_Config_List_Toml) {
	//connectionMap = make(map[string]*ConnectionPool)
	for _, dbs_config := range dbsConfigListToml.Dbs {
		////	fmt.Println(dbs_config)
		var cp ConnectionPool
		cp.NewConnectionPool(dbs_config.MaxPoolSize, dbs_config.Name, dbs_config.DriverName, dbs_config.DataSourceName)
		connectionMap[dbs_config.Name] = &cp
	}
}

//
// 根据 Db_SqlConfig_List_Toml 配置信息，生成sql配置信息
// Db_SqlConfig_List_Toml 对应的配置文件是 apps/app1~n/config/db/xxxx.toml ， 每个app下的配置信息都被加载。
// 其中： xxxx名称是 Dbs_Config_Toml.Name的名字
//
func initSqlMap(dbsConfigListToml Dbs_Config_List_Toml, appName string) {
	dirs, err := ioutil.ReadDir(env.Env_Project_Resource + "/" + appName)
	if err != nil {
		log4go.ErrorLog(err)
	} else {
		//sqlConfig.SqlMap = make(map[string]SqlMetas)

		for _, dir := range dirs {
			for _, dbs_config := range dbsConfigListToml.Dbs {
				appname := dir.Name()
				pathsqlconfig := env.Env_Project_Resource + "/" + appName + "/" + appname + "/" + env.Env_Project_Resource_Apps_SqlConfig + "/" + dbs_config.Name + ".toml"
				//// fmt.Println(pathsqlconfig)
				if utils.ExistFile(pathsqlconfig) {
					var dbSqlConfigListToml Db_SqlConfig_List_Toml
					if _, err := toml.DecodeFile(pathsqlconfig, &dbSqlConfigListToml); err != nil {
						log4go.ErrorLog("Reading file(%s) is error . err: %s .", pathsqlconfig, err.Error())
					} else {
						//// fmt.Println(dbSqlConfigListToml)
						var sqlMetas SqlMetas
						if len(sqlConfig.SqlMap[dbs_config.Name].SqlMetaMap) > 0 {
							sqlMetas = sqlConfig.SqlMap[dbs_config.Name]
						} else {
							sqlMetas.SqlMetaMap = make(map[string]SqlMeta)
						}
						for _, dbSqlConfig := range dbSqlConfigListToml.Db {
							sqlMeta := SqlMeta{dbSqlConfig.Name, useCurrentSql(dbSqlConfig, dbs_config.DriverName)}
							sqlMetas.SqlMetaMap[appname+"."+dbSqlConfig.Name] = sqlMeta
						}
						sqlConfig.SqlMap[dbs_config.Name] = sqlMetas
					}
				}
			}
		}

		if log4go.IsFinestLevel() {
			log4go.FinestLog("sqlConfig:  %v \n", sqlConfig.SqlMap)
		}

	}
	/*
		for _, dbs_config := range dbsConfigListToml.Dbs {
			fmt.Println(dbs_config)
		}
	*/
}

func useCurrentSql(dbSqlConfig Db_SqlConfig_Toml, driverName string) string {
	if driverName == "postgres" {
		if dbSqlConfig.Sql_postgres != "" {
			return dbSqlConfig.Sql_postgres
		}
	} else if driverName == "mysql" {
		if dbSqlConfig.Sql_mysql != "" {
			return dbSqlConfig.Sql_mysql
		}
	} else if driverName == "sqlite3" {
		if dbSqlConfig.Sql_sqlite3 != "" {
			return dbSqlConfig.Sql_sqlite3
		}
	}
	return dbSqlConfig.Sql
}

// 根据在dbconfig.toml中配置的数据库别名，获取该数据库连接池
func GetConnectionMap(name string) (cp *ConnectionPool) {
	cp = connectionMap[name]
	return
}

func GetSqlConfig() SqlConfig {
	return sqlConfig
}
