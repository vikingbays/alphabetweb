// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package cache

import (
	"alphabet/core/toml"
	"alphabet/env"
	"alphabet/log4go"
	"alphabet/log4go/message"
)

type Caches_Config_List_Toml struct {
	Caches []Caches_Config_Toml
}

type Caches_Config_Toml struct {
	Name           string
	DataSourceName string
	MaxPoolSize    int
}

var cacheConnectionMap map[string]*CachePool

// 初始化redis的cache连接，sql注册等
func Init() {
	cacheConnectionMap = make(map[string]*CachePool)

	var cachesConfigListToml Caches_Config_List_Toml
	for _, appName := range env.Env_Project_Resource_Apps {
		pathconfig := env.Env_Project_Resource + "/" + appName + "/" + env.GetSysconfigPath() + env.Env_Project_Resource_CacheConfig
		if _, err := toml.DecodeFile(pathconfig, &cachesConfigListToml); err != nil {
			//log4go.ErrorLog(err)
			log4go.ErrorLog(message.ERR_CACH_39002, pathconfig, err.Error())
		} else {
			initCacheConnectionMap(cachesConfigListToml)
		}
	}

}

/*
 根据Caches_Config_List_Toml配置信息，生成连接池
 Caches_Config_List_Toml 对应的配置文件是 apps/cachesconfig.toml ，可以配置多个缓存库连接
*/
func initCacheConnectionMap(cachesConfigListToml Caches_Config_List_Toml) {
	for _, caches_config := range cachesConfigListToml.Caches {
		////	fmt.Println(dbs_config)
		var cp CachePool
		cp.NewCachePool(caches_config.MaxPoolSize, caches_config.Name, caches_config.DataSourceName)
		cacheConnectionMap[caches_config.Name] = &cp
	}
}

// 根据在cacheconfig.toml中配置的缓存库别名，获取该数据库连接池
//
// @param name 缓存库别名
//
// @return 返回一个连接
func GetCacheConnectionMap(name string) (cp *CachePool) {
	cp = cacheConnectionMap[name]
	return
}
