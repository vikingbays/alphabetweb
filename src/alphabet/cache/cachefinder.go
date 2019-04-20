// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package cache

import (
	"alphabet/core/redis"
	//	"alphabet/log4go"
	"fmt"
	"time"
)

type CacheFinder struct {
	cache redis.Conn
	cp    *CachePool
}

/*
获取一个连接。

获取连接的方法是：
  cacheFinder1, err := cache.GetCacheFinder("cache1")
  defer cacheFinder1.Close()

@param connectionName  缓存库连接池别名

@return 返回DataFinder对象
*/
func GetCacheFinder(connectionName string) (*CacheFinder, error) {
	cacheFinder := &CacheFinder{}
	cacheFinder.cp = GetCacheConnectionMap(connectionName)
	if cacheFinder.cp == nil {
		return cacheFinder, &ConnectionError{fmt.Errorf("connection error...")}
	} else {
		//cacheFinder.cp.StartTime = time.Now().UnixNano()
		//log4go.InfoLog("GetCache()  [start]>>> cache duration: %d", cacheFinder.cp.StartTime)
		cacheFinder.cache = cacheFinder.cp.GetCache()
		//log4go.InfoLog("GetCache()  [  end]>>> cache duration: %d", (time.Now().UnixNano() - cacheFinder.cp.StartTime))
		if cacheFinder.cache == nil {
			return cacheFinder, &ConnectionError{fmt.Errorf("connection error...")}
		}
	}

	return cacheFinder, nil
}

/*
完成操作后释放连接
*/
func (cf *CacheFinder) Close() {
	if cf.cp != nil && cf.cache != nil { //当前数据库连接正常。如果数据库连接异常，连接池不需要回收
		cf.cp.ReleaseCache(cf.cache)

		cf.cp.EndTime = time.Now().UnixNano()
		//log4go.InfoLog("cache duration: %d", (cf.cp.EndTime - cf.cp.StartTime))
	}
}

//直接获取redis对象，可以进行原生操作，例如:
//  cf.GetCache().Do(...)
func (cf *CacheFinder) GetCache() redis.Conn {
	return cf.cache
}

//
//设置Map信息一条记录：field:value
//
func (cf *CacheFinder) SetMap(mapKeyId string, field string, value string) {
	cf.cache.Do("HSET", mapKeyId, field, value)
}

//设置失效时间，时间单位是秒
func (cf *CacheFinder) Expire(key string, timesecond int) {
	cf.cache.Do("EXPIRE", key, timesecond)
}

//
//根据map对象唯一标示（mapKeyId），获取该map信息
//
func (cf *CacheFinder) GetMap(mapKeyId string) map[string]string {

	values, _ := redis.Values(cf.cache.Do("HGETALL", mapKeyId))
	map1 := make(map[string]string)
	key1 := ""
	value1 := ""
	for index, v := range values {
		if index%2 == 0 {
			key1 = string(v.([]byte))
		} else {
			value1 = string(v.([]byte))
			if key1 != "" && value1 != "" {
				map1[key1] = value1
			}
		}
	}

	return map1
}

//
// 根据map对象中的field获取对应的value值。
//
//
//
//
func (cf *CacheFinder) GetMapField(mapKeyId string, field string) string {
	value := ""
	data, err := cf.cache.Do("HGET", mapKeyId, field)
	if err == nil && data != nil {
		value = string(data.([]byte))
	}

	return value
}

//
// 删除该map信息
//
func (cf *CacheFinder) DelMap(mapKeyId string) {
	cf.cache.Do("DEL", mapKeyId)

}

//
// 删除该map中其中一个field
//
func (cf *CacheFinder) DelMapField(mapKeyId string, field string) {
	cf.cache.Do("HDEL", mapKeyId, field)
}
