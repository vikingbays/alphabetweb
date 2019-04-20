// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*
用于定义服务端管理的基础模块，基类，例如：基于ETCD的管理
*/

package service

import (
	"alphabet/env"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"fmt"
	"strconv"
	"strings"
)

// 服务管理类入口，
type BaseServiceManager struct {
	typeOfService string // 如果是服务端 typeOfService="server" ，如果是客户端 typeOfService="client"
	configName    string // 配置文件名称
	i             IServiceManager
}

func (b *BaseServiceManager) RegGroupAndRpc(group *ServiceManagerGroup) {
	regFlag := false // 是否需要发起注册, 如果是true，表示需要注册
	gId := group.GroupId
	storedGroup, err := b.i.queryGroupById(gId)
	if storedGroup == nil || err != nil {
		regFlag = true //  groupId在注册中心不存在，需要注册
	} else if len(storedGroup.Addrs) == 0 {
		regFlag = true // 没有地址
	}

	if regFlag {

		err0 := b.i.delGroup(gId)

		if err0 != nil {
			log4go.ErrorLog(message.ERR_MSSV_39021, err0.Error())
		}

		err1 := b.i.saveGroupWithoutAddr(group)
		if err1 != nil {
			log4go.ErrorLog(message.ERR_MSSV_39022, err1.Error())
		}

	}
}

func (b *BaseServiceManager) RegAddr(groupId string, addr *ServiceManagerGroupAddr, rootAppname string) {
	err := b.i.saveAddr(groupId, addr)
	if err == nil {
		b.i.ttlAddrByProvider(groupId, addr.Ip, addr.Port, SERVICE_MS_ADDR_TTL_TIME)
		b.i.sendTicketByProvider(groupId, addr.Ip, addr.Port, rootAppname)
	}
}

func (b *BaseServiceManager) GetGroupById(groupId string) *ServiceManagerGroup {
	storedGroup, err := b.i.queryGroupById(groupId)
	if err != nil {
		log4go.ErrorLog(message.ERR_MSSV_39023, err.Error())
	}
	return storedGroup
}

func (b *BaseServiceManager) GetGroupIds() []string {
	return b.i.GetGroupIds()
}

func (b *BaseServiceManager) UpdateClientGroup(groupId string, delFlag bool) {
	for _, clients := range clientsGlobalMap {
		if delFlag { // 删除操作
			flexServicePoolGroup.DeleteServicePoolGroup(groupId)
			delete(clients.dataMap, groupId)
		} else { // 更新操作
			if clients.managerMap[clients.groupidToFilenameMap[groupId]] != nil {

				serviceManagerGroup := clients.managerMap[clients.groupidToFilenameMap[groupId]].GetGroupById(groupId)
				if serviceManagerGroup != nil {
					log4go.InfoLog(message.INF_MSSV_09004, serviceManagerGroup.Json())
					serviceManagerGroup.MaxPoolSize = clients.configMap[clients.groupidToFilenameMap[groupId]].Client[0].MaxPoolSize
					serviceManagerGroup.ReqPerConn = clients.configMap[clients.groupidToFilenameMap[groupId]].Client[0].ReqPerConn
					clients.dataMap[groupId] = serviceManagerGroup
					for _, addr := range serviceManagerGroup.Addrs {
						log4go.InfoLog(message.INF_MSSV_09005, serviceManagerGroup.GroupId, addr.Ip, addr.Port)
						if addr.GetProtocolInfo() != "rpc_tcp_ssl" {
							if serviceManagerGroup.ReqPerConn <= 0 {
								serviceManagerGroup.ReqPerConn = 1
							}
							maxConn := serviceManagerGroup.MaxPoolSize * serviceManagerGroup.ReqPerConn
							if maxConn > env.Env_MS_MaxConn {
								maxConn = env.Env_MS_MaxConn
							}
							serviceManagerGroup.MaxPoolSize = maxConn
							serviceManagerGroup.ReqPerConn = 1
							log4go.WarnLog(message.WAR_MSSV_69006, fmt.Sprintf("%s ,%s:%d", addr.GetProtocolInfo(), addr.Ip, addr.Port),
								clients.configMap[clients.groupidToFilenameMap[groupId]].Client[0].MaxPoolSize,
								clients.configMap[clients.groupidToFilenameMap[groupId]].Client[0].ReqPerConn,
								serviceManagerGroup.MaxPoolSize, serviceManagerGroup.ReqPerConn)

						}
						flexServicePoolGroup.UpdateServicePool(serviceManagerGroup.MaxPoolSize, serviceManagerGroup.ReqPerConn,
							serviceManagerGroup.GroupId, addr)
					}
				}
			} else {
				log4go.ErrorLog(message.ERR_MSSV_39024, groupId, clients.groupidToFilenameMap[groupId])
			}
		}
	}

}

func (b *BaseServiceManager) UpdateClientGroupAddr(groupId string, addr string, delFlag bool) {
	for _, clients := range clientsGlobalMap {
		if delFlag { // 删除操作
			flexServicePoolGroup.DeleteServicePool(groupId, clients.dataMap[groupId].Addrs[addr])
			delete(clients.dataMap[groupId].Addrs, addr)
		} else { // 更新操作,只是更新这个地址信息
			idx_ipAndPort := strings.LastIndex(addr, "_")
			ipAndPort := []string{addr[0:idx_ipAndPort], addr[(idx_ipAndPort + 1):len(addr)]}
			//ipAndPort := strings.Split(addr, "_")
			port, _ := strconv.Atoi(ipAndPort[1])
			serviceManagerGroupAddr := b.i.queryAddr(groupId, ipAndPort[0], port)
			if serviceManagerGroupAddr != nil {
				//			serviceManagerGroupAddr.Group = clients.dataMap[groupId]
				if clients.managerMap[clients.groupidToFilenameMap[groupId]] != nil {
					MaxPoolSize := clients.configMap[clients.groupidToFilenameMap[groupId]].Client[0].MaxPoolSize
					ReqPerConn := clients.configMap[clients.groupidToFilenameMap[groupId]].Client[0].ReqPerConn
					if serviceManagerGroupAddr.GetProtocolInfo() != "rpc_tcp_ssl" {
						if ReqPerConn <= 0 {
							ReqPerConn = 1
						}
						maxConn := MaxPoolSize * ReqPerConn
						if maxConn > env.Env_MS_MaxConn {
							maxConn = env.Env_MS_MaxConn
						}
						MaxPoolSize = maxConn
						ReqPerConn = 1

						log4go.WarnLog(message.WAR_MSSV_69006, fmt.Sprintf("%s ,%s:%d", serviceManagerGroupAddr.GetProtocolInfo(), serviceManagerGroupAddr.Ip, serviceManagerGroupAddr.Port),
							clients.configMap[clients.groupidToFilenameMap[groupId]].Client[0].MaxPoolSize,
							clients.configMap[clients.groupidToFilenameMap[groupId]].Client[0].ReqPerConn,
							MaxPoolSize, ReqPerConn)

					}

					flexServicePoolGroup.UpdateServicePool(MaxPoolSize, ReqPerConn, groupId, serviceManagerGroupAddr)
					clients.dataMap[groupId].Addrs[addr] = serviceManagerGroupAddr
				} else {
					log4go.ErrorLog(message.ERR_MSSV_39024, groupId, clients.groupidToFilenameMap[groupId])
				}
			}
		}
	}
}
