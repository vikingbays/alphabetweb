// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package service

import (
	"alphabet/core/utils"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"sync"
)

// 计数器，用于识别下一个连接指针
type fspg_counter struct {
	num int
}

// rpc服务的连接池管理类，支撑 多个 地址的连接管理，每个地址的子连接池对应的是ServicePool
type FlexServicePoolGroup struct {
	poolMap    map[string]map[string]*ServicePool2 // poolMap[groupId][ipAndPort]=new(ServicePool)
	counterMap map[string]*fspg_counter            // 根据groupId 设置计数器信息
	lockMap    map[string]*sync.Mutex              // 根据groupId 设置锁信息
	pointerMap map[string][]string                 // 设置每个groupId的所有addr的排序，因为poolMap是无序状态，所以在做轮询时候有问题
}

// 初始化 NewFlexServicePoolGroup 对象
func (fspg *FlexServicePoolGroup) NewFlexServicePoolGroup() {
	fspg.poolMap = make(map[string]map[string]*ServicePool2)
	fspg.counterMap = make(map[string]*fspg_counter)
	fspg.lockMap = make(map[string]*sync.Mutex)
	fspg.pointerMap = make(map[string][]string)
}

// 添加或者更新一个ServicePool
func (fspg *FlexServicePoolGroup) UpdateServicePool(maxPoolSize int, reqPerConn int, groupId string, addr *ServiceManagerGroupAddr) {
	ipAndPort := merge_ip_and_port(addr.Ip, addr.Port)
	fspg.DeleteServicePool(groupId, addr)

	var servicePool2 *ServicePool2 = new(ServicePool2)
	servicePool2.NewServicePool(maxPoolSize, reqPerConn, groupId, addr)
	if fspg.poolMap[groupId] == nil {
		fspg.poolMap[groupId] = make(map[string]*ServicePool2)
	}

	fspg.poolMap[groupId][ipAndPort] = servicePool2

	if fspg.counterMap[groupId] == nil {
		fspg.counterMap[groupId] = new(fspg_counter)
	}
	if fspg.lockMap[groupId] == nil {
		fspg.lockMap[groupId] = new(sync.Mutex)
	}

	if fspg.pointerMap[groupId] == nil {
		fspg.pointerMap[groupId] = make([]string, 0, 1)
	}
	fspg.pointerMap[groupId] = append(fspg.pointerMap[groupId], ipAndPort)

}

// 根据addr移除一个ServicePool
func (fspg *FlexServicePoolGroup) DeleteServicePool(groupId string, addr *ServiceManagerGroupAddr) {
	if addr == nil {
		return
	}
	ipAndPort := merge_ip_and_port(addr.Ip, addr.Port)
	if fspg.pointerMap[groupId] != nil { // 移除无效的addr
		addrs := fspg.pointerMap[groupId]
		newAddrs := make([]string, 0, 1)
		for _, value := range addrs {
			if value != ipAndPort {
				newAddrs = append(newAddrs, value)
			}
		}
		fspg.pointerMap[groupId] = newAddrs
	}
	if fspg.poolMap[groupId] != nil {
		{ // 销毁原来的服务连接
			var oldServicePool *ServicePool2 = nil
			if fspg.poolMap[groupId][ipAndPort] != nil {
				oldServicePool = fspg.poolMap[groupId][ipAndPort]
			}
			if oldServicePool != nil {
				lenObjectStack := oldServicePool.ObjectStack.Len()
				for i := 0; i < lenObjectStack; i++ {
					obj := oldServicePool.ObjectStack.Pop()
					if obj != nil {
						obj.(*FlexRpcClientWrapper).Close()
					}
					obj = nil
				}
				oldServicePool.ObjectStack = nil
				oldServicePool.Cond = nil
				oldServicePool.Addr = nil
				oldServicePool.ObjectFactory = nil
				oldServicePool.Parent = nil
				oldServicePool.Lock = nil
				oldServicePool = nil
			}
		}

		delete(fspg.poolMap[groupId], ipAndPort)
	}

}

// 根据groupId移除他的所有的ServicePool
func (fspg *FlexServicePoolGroup) DeleteServicePoolGroup(groupId string) {
	delete(fspg.poolMap, groupId)
}

// 根据groupId获取一个数据库连接
func (fspg *FlexServicePoolGroup) GetConnection(groupId string) *FlexRpcClientWrapper {
	if fspg.poolMap[groupId] == nil {
		return nil
	}
	if len(fspg.pointerMap[groupId]) == 0 { // 如果没有可使用的连接池对象
		return nil
	}
	fspg.lockMap[groupId].Lock()
	poolGroupLength := len(fspg.pointerMap[groupId])
	currNum := fspg.counterMap[groupId].num % poolGroupLength
	if currNum == 0 && fspg.counterMap[groupId].num != 0 {
		fspg.counterMap[groupId].num = 1
	} else {
		fspg.counterMap[groupId].num = fspg.counterMap[groupId].num + 1
	}
	fspg.lockMap[groupId].Unlock()
	if len(fspg.poolMap[groupId]) > currNum && len(fspg.pointerMap[groupId]) > currNum {
		if fspg.poolMap[groupId][fspg.pointerMap[groupId][currNum]] == nil {
			log4go.ErrorLog(message.ERR_MSSV_39026, fspg.pointerMap[groupId][currNum])
			return nil
		} else {
			log4go.DebugLog(message.DEG_CORE_79004, poolGroupLength, currNum, groupId, fspg.pointerMap[groupId][currNum])
			return fspg.poolMap[groupId][fspg.pointerMap[groupId][currNum]].GetConnection()
		}
	} else {
		log4go.ErrorLog(message.ERR_MSSV_39027, len(fspg.poolMap[groupId]), len(fspg.pointerMap[groupId]), currNum)
		return nil
	}

}

// 释放一个rpc连接
func (sspg *FlexServicePoolGroup) ReleaseConnection(fpc *FlexRpcClientWrapper) {
	ipAndPort := merge_ip_and_port(fpc.Ip, fpc.Port)
	sspg.poolMap[fpc.GroupId][ipAndPort].Release(fpc)
}

// 相比较 SimpleRpcClient 增加 groupId的管理
type FlexRpcClientWrapper struct {
	SimpleRpcClient
	GroupId string
}

func (frc *FlexRpcClientWrapper) ReConn() {
	frc.Close()
	err := frc.Connect()
	if err != nil {
		log4go.ErrorLog(message.ERR_MSSV_39028, err.Error())
	} else {
		frc.head_nest("/", true)
	}
}

func (frc *FlexRpcClientWrapper) CloneRpcClient() *SimpleRpcClient {
	sc := new(SimpleRpcClient)
	sc.Ip = frc.Ip
	sc.Protocol = frc.Protocol
	sc.Port = frc.Port
	sc.prefix = frc.prefix
	sc.client = frc.client
	return sc
}

type ServicePool2 struct {
	utils.AbstractPool2
	GroupId   string
	IpAndPort string
	Addr      *ServiceManagerGroupAddr
}

func (sp *ServicePool2) NewServicePool(maxPoolSize int, reqPerConn int, groupId string, addr *ServiceManagerGroupAddr) {
	sp.Init()
	sp.MaxPoolSize = maxPoolSize
	sp.ReqPerConn = reqPerConn
	sp.ObjectFactory = &RpcObjectFactory2{groupId, addr}
	sp.TryTimes = SERVICE_MS_POOL_TRY_TIMES
	sp.GroupId = groupId
	sp.IpAndPort = merge_ip_and_port(addr.Ip, addr.Port)
	sp.Addr = addr
	connCount, reqCount := sp.CreateObjects()
	log4go.InfoLog(message.INF_MSSV_09016, groupId, reqCount, connCount, addr.GetProtocolInfo(), sp.IpAndPort)
}

//从连接池中获取一个RPC连接
func (sp *ServicePool2) GetConnection() *FlexRpcClientWrapper {
	obj := sp.Get()
	if obj != nil {
		return obj.(*FlexRpcClientWrapper)
	}
	return nil
}

//释放RPC连接到连接池
func (sp *ServicePool2) ReleaseConnection(spc *FlexRpcClientWrapper) {
	sp.Release(spc)
}

type RpcObjectFactory2 struct {
	GroupId string
	Addr    *ServiceManagerGroupAddr
}

// 创建对象，例如：数据库连接
func (rpc *RpcObjectFactory2) Create() utils.IConnectionObject {
	//log4go.InfoLog("  Create()   ........   GroupId=%s, Addr=%s   \n", rpc.GroupId, rpc.Addr.Json())
	log4go.DebugLog(message.DEG_MSSV_79011, rpc.GroupId, rpc.Addr.Json())

	sc := &FlexRpcClientWrapper{}
	sc.Protocol = rpc.Addr.GetProtocolInfo()
	sc.Ip = rpc.Addr.Ip
	sc.Port = rpc.Addr.Port
	sc.GroupId = rpc.GroupId

	sc.ReConn()
	return sc
}

// 验证对象是否有效（或运行正常）
func (rpc *RpcObjectFactory2) Valid(obj utils.IConnectionObject) bool {
	if obj != nil {
		sc := obj.(*FlexRpcClientWrapper)
		return sc.Ping()
	}
	return false
}

// AbstractPool的Release方法，释放对象前调用
func (rpc *RpcObjectFactory2) ReleaseStart(obj utils.IConnectionObject) {

}

// AbstractPool的Release方法，释放对象后调用
func (rpc *RpcObjectFactory2) ReleaseEnd(obj utils.IConnectionObject) {

}

// AbstractPool的Get方法，获取对象前调用
func (rpc *RpcObjectFactory2) GetStart() {

}

// AbstractPool的Get方法，获取对象后调用
func (rpc *RpcObjectFactory2) GetEnd(obj utils.IConnectionObject) {

}
