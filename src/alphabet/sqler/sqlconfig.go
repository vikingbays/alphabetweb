// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sqler

/*
SqlConfig为入口类，记录所有SqlContainer的配置
*/
type SqlConfig struct {
	SqlMap map[string]SqlMetas
}

type SqlMetas struct {
	SqlMetaMap map[string]SqlMeta
}

type SqlMeta struct {
	Key string
	Sql string
}
