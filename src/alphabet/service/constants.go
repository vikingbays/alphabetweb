// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*
 定义微服务的常量信息
*/

package service

import "time"

// 监听服务异常退出后，隔多长时间重试，单位是 :秒 .
var SERVICE_MS_WATCH_EXIT_SLEEP time.Duration = time.Duration(1)

// 设置绑定到具体地址上的票据（ticket）更新的间隔时间，单位是 :分钟 .
var SERVICE_MS_SEND_ADDR_TICKET_DURATION time.Duration = time.Duration(30)

// 设置服务端的超时设置，单位：秒。
var SERVICE_MS_SERVER_ADDR_TIMEOUT int = 10

// 设置服务端的权重
var SERVICE_MS_SERVER_ADDR_WEIGHT int = 10

//  设置监听服务存储的 消息chan的缓冲区大小，默认：200
var SERVICE_MS_CHAN_WATCH_STORE_LENGTH int = 200

// 设置 addr ttl 的有效时间。单位：秒
var SERVICE_MS_ADDR_TTL_TIME int64 = 5

// 设置 addr ttl 的sleep间隔时间。单位：秒
var SERVICE_MS_ADDR_TTL_SLEEP_TIME int = 3

// 设置监听服务的处理时，每次批处理的间隔时间，单位是秒
var SERVICE_MS_WATCH_DEAL_BATCH_DURATION time.Duration = time.Duration(1)

//  设置BaseServiceManager.typeOfService 的服务端信息标示
var SERVICE_MS_BASE_TYPE_SERVER string = "server"

//  设置BaseServiceManager.typeOfService 的客户端信息标示
var SERVICE_MS_BASE_TYPE_CLIENT string = "client"

// 设置service连接池的重试次数
var SERVICE_MS_POOL_TRY_TIMES int = 5

// 设置请求处理时候，存储在header中的票据名称
var SERVICE_MS_REQ_HEADER_TICKET_NAME = "ticket_current"

// 设置请求处理时候，存储在header中的票据名称(上一个)
var SERVICE_MS_REQ_HEADER_TICKET_LAST_NAME = "ticket_last"
