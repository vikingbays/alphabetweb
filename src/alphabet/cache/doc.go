// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

//
// 提供cache处理能力。
//
// 获取连接的方法是：
//  cacheFinder1, err := cache.GetCacheFinder("cache1")
//  defer cacheFinder1.Close()
//
// 设置Map信息一条记录：field:value
//  cacheFinder1.SetMap(mapId, fieldId, value)
//
// 设置失效时间，时间单位是秒
//  cacheFinder1.Expire(mapId , timeoutSecond)
//
// 根据map对象唯一标示（mapKeyId），获取该map信息
//  cacheFinder1.GetMap(mapId)
//
// 根据map对象中的field获取对应的value值。
//  cacheFinder1.GetMapField(mapId, fieldId)
//
// 删除该map信息
//  cacheFinder1.DelMap(mapId)
//
// 删除该map中其中一个field
//  cacheFinder1.DelMapField(mapId,fieldId)
//
package cache
