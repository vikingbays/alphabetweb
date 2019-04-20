// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sqler

import (
	_ "alphabet/core/mysql"
	_ "alphabet/core/pq"
	"alphabet/log4go/message"
	//_ "alphabet/core/sqlite3"
	"alphabet/core/utils"
	"alphabet/log4go"
	"database/sql"
)

/*
初始化：
  var cp ConnectionPool
  cp.NewConnectionPool(db_config.MaxPoolSize, db_config.Name, db_config.DriverName, db_config.DataSourceName)
获取一个连接：
  db := cp.GetConnection()
释放一个连接
  cp.ReleaseConnection(db)

*/
type ConnectionPool struct {
	utils.AbstractPool
	PreparedSqlStandard bool
}

func (c *ConnectionPool) NewConnectionPool(maxPoolSize int, name string, driverName string, dataSourceName string) {
	c.Init()
	c.MaxPoolSize = maxPoolSize
	c.ObjectFactory = &DBObjectFactory{name, driverName, dataSourceName}
	if driverName == "postgres" || driverName == "postgresql" {
		c.PreparedSqlStandard = false
	} else {
		c.PreparedSqlStandard = true
	}
	c.TryTimes = SQLER_CONNECTION_POOL_TRY_TIMES
	count := c.CreateObjects(maxPoolSize)
	log4go.InfoLog(message.INF_SQLR_09017, name, driverName, count, dataSourceName)
}

//从连接池中获取数据库连接
func (c *ConnectionPool) GetConnection() *sql.DB {
	obj := c.Get()
	if obj != nil {
		return obj.(*sql.DB)
	}
	return nil
}

//释放数据库连接到连接池
func (c *ConnectionPool) ReleaseConnection(db *sql.DB) {
	c.Release(db)
}

type DBObjectFactory struct {
	Name           string
	DriverName     string
	DataSourceName string
}

// 创建对象，例如：数据库连接
func (c *DBObjectFactory) Create() interface{} {
	log4go.DebugLog(message.DEG_SQLR_79012, c.DriverName, c.DataSourceName)

	db, err := sql.Open(c.DriverName, c.DataSourceName)
	if err != nil {
		log4go.ErrorLog(message.ERR_SQLR_39066, c.DriverName, c.DataSourceName, err)
		return nil
	} else {
		db.SetMaxIdleConns(1)
		db.SetMaxOpenConns(1)
		return db
	}
}

// 验证对象是否有效（或运行正常）
func (c *DBObjectFactory) Valid(obj interface{}) bool {
	if obj != nil {
		db := obj.(*sql.DB)
		err := db.Ping()
		if err != nil {
			return false
		} else {
			return true
		}

	}
	return false
}

// AbstractPool的Release方法，释放对象前调用
func (c *DBObjectFactory) ReleaseStart(obj interface{}) {

}

// AbstractPool的Release方法，释放对象后调用
func (c *DBObjectFactory) ReleaseEnd(obj interface{}) {

}

// AbstractPool的Get方法，获取对象前调用
func (c *DBObjectFactory) GetStart() {

}

// AbstractPool的Get方法，获取对象后调用
func (c *DBObjectFactory) GetEnd(obj interface{}) {

}
