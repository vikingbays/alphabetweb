// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package cache

/*
cache.toml 配置cache信息

设置session采用方式 sessionstore="cache001" , 如果为 "" 就是存储到mem中，如果是具体值"cache001"表示名称是"cache001"的redis连接池


*/

import (
	"alphabet/core/redis"
	"alphabet/core/utils"
	"alphabet/log4go"
	"alphabet/log4go/message"
)

/*
初始化：
  var cp CachePool
  cp.NewCachePool(db_config.MaxPoolSize, db_config.Name, db_config.DataSourceName)
获取一个连接：
  db := cp.GetCache()
释放一个连接
  cp.ReleaseCache(db)

*/
type CachePool struct {
	utils.AbstractPool
	StartTime int64
	EndTime   int64
}

func (c *CachePool) NewCachePool(maxPoolSize int, name string, dataSourceName string) {
	c.Init()
	c.MaxPoolSize = maxPoolSize
	c.ObjectFactory = &CacheObjectFactory{name, dataSourceName}
	c.TryTimes = CACHE_POOL_TRY_TIMES
	count := c.CreateObjects(maxPoolSize)
	log4go.InfoLog(message.INF_CACH_09018, name, count, dataSourceName)
}

//从连接池中获取数据库连接
func (c *CachePool) GetCache() redis.Conn {
	obj := c.Get()
	if obj != nil {
		return obj.(redis.Conn)
	}
	return nil
}

//释放数据库连接到连接池
func (c *CachePool) ReleaseCache(cacheObj redis.Conn) {
	c.Release(cacheObj)
}

type CacheObjectFactory struct {
	Name           string
	DataSourceName string // redis://user:abcdef123456@127.0.0.1:6379/0   采用tcp协议，密码是abcdef123456 ， 数据库是0
	//redis://user:abcdef123456@/tmp/redisserv/redis.sock/1   采用unix协议，密码是abcdef123456 ， 数据库是1
}

// 创建对象，例如：数据库连接
func (c *CacheObjectFactory) Create() interface{} {
	cacheObject, err := redis.DialURL(c.DataSourceName) //  redis.Conn
	if err != nil {
		log4go.ErrorLog(message.ERR_CACH_39003, err.Error())
		return nil
	} else {
		return cacheObject
	}
}

// 验证对象是否有效（或运行正常）
func (c *CacheObjectFactory) Valid(obj interface{}) bool {
	if obj != nil {
		cacheObject := obj.(redis.Conn)

		data, err1 := cacheObject.Do("PING")
		if err1 == nil && data != nil && data == "PONG" {
			return true
		} else {
			log4go.ErrorLog(message.ERR_CACH_39004, err1.Error())
			return false
		}
	}
	return false
}

// AbstractPool的Release方法，释放对象前调用
func (c *CacheObjectFactory) ReleaseStart(obj interface{}) {

}

// AbstractPool的Release方法，释放对象后调用
func (c *CacheObjectFactory) ReleaseEnd(obj interface{}) {

}

// AbstractPool的Get方法，获取对象前调用
func (c *CacheObjectFactory) GetStart() {

}

// AbstractPool的Get方法，获取对象后调用
func (c *CacheObjectFactory) GetEnd(obj interface{}) {

}
