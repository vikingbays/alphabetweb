// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*
用于定义服务端管理模块的具体实现，以etcd方式实现。
*/

package service

import (
	"alphabet/core/etcd"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type EtcdServiceManager struct {
	BaseServiceManager
	endpoints          []string
	username           string
	password           string
	timeout            time.Duration
	root               string
	sendTicketDuration time.Duration
	chanWatchStore     chan WatchGroupInfoStore // 内部，不需要初始化 ， 客户端监听时使用
	arrayWatchStore    []WatchGroupInfoStore    // 内部，不需要初始化 ， 客户端监听时使用
	lockForWS          *sync.Mutex              // 内部，不需要初始化 ， 客户端监听时使用
}

func (e *EtcdServiceManager) queryGroupById(groupId string) (*ServiceManagerGroup, error) {
	//e.saveAddr(groupId, ip, port, weight, false)

	etcdConnection := etcd.EtcdConnection{}
	etcdConnection.CreateConnection(e.endpoints, e.timeout, e.username, e.password)
	defer etcdConnection.Close()
	groupIdPrefix := e.root + ".groups" + "." + groupId
	rpcsPrefix := groupIdPrefix + ".rpcs"
	addrsPrefix := groupIdPrefix + ".addrs"

	kvMaps, err := etcdConnection.GetWithPrefix(groupIdPrefix)

	if err == nil && len(kvMaps) > 0 {
		serviceManagerGroup := &ServiceManagerGroup{}

		{ // 初始化：serviceManagerGroup
			serviceManagerGroup.Addrs = make(map[string]*ServiceManagerGroupAddr)
			serviceManagerGroup.Rpcs = make(map[string]*ServiceManagerGroupRpc)

			serviceManagerGroup.GroupId = groupId

			serviceManagerGroup.Timestamp = kvMaps[groupIdPrefix+".timestamp"]
		}

		{ // 初始化serviceManagerGroup.Rpcs
			str_rpcs := kvMaps[rpcsPrefix]
			if str_rpcs != "" {
				rpcs := strings.Split(str_rpcs, ",")

				for _, rpcId := range rpcs {
					serviceManagerGroupRpc := &ServiceManagerGroupRpc{}
					//					serviceManagerGroupRpc.Group = serviceManagerGroup

					serviceManagerGroupRpc.RpcId = rpcId
					serviceManagerGroupRpc.Desc = kvMaps[rpcsPrefix+"."+rpcId+".desc"]
					serviceManagerGroupRpc.Url = kvMaps[rpcsPrefix+"."+rpcId+".url"]
					available, _ := strconv.ParseBool(kvMaps[rpcsPrefix+"."+rpcId+".available"])
					serviceManagerGroupRpc.Available = available
					serviceManagerGroup.Rpcs[rpcId] = serviceManagerGroupRpc

				}
			}

		}

		{ // 初始化serviceManagerGroup.Addrs

			//			serviceManagerGroupAddr.Group = serviceManagerGroup

			str_addrs := kvMaps[addrsPrefix]
			if str_addrs != "" {
				addrs := strings.Split(str_addrs, ",")
				for _, addr := range addrs {
					addrSamplePrefix := fmt.Sprintf("%s%s", addrsPrefix, ".["+addr+"]")
					if kvMaps[addrSamplePrefix] != "" {
						serviceManagerGroupAddr := &ServiceManagerGroupAddr{}
						serviceManagerGroupAddr.IpType = kvMaps[addrSamplePrefix+".iptype"]
						timeout, _ := strconv.Atoi(kvMaps[addrSamplePrefix+".timeout"])
						serviceManagerGroupAddr.Timeout = timeout
						protocol, _ := strconv.Atoi(kvMaps[addrSamplePrefix+".protocol"])
						serviceManagerGroupAddr.Protocol = ProtocolType(protocol)
						serviceManagerGroupAddr.Ip = kvMaps[addrSamplePrefix+".ip"]
						port, _ := strconv.Atoi(kvMaps[addrSamplePrefix+".port"])
						serviceManagerGroupAddr.Port = port
						weight, _ := strconv.Atoi(kvMaps[addrSamplePrefix+".weight"])
						serviceManagerGroupAddr.Weight = weight
						serviceManagerGroupAddr.Ticket = kvMaps[addrSamplePrefix+".ticket"]
						serviceManagerGroupAddr.WebContext = kvMaps[addrSamplePrefix+".webcontext"]
						serviceManagerGroup.Addrs[fmt.Sprintf("%s_%d", serviceManagerGroupAddr.Ip, serviceManagerGroupAddr.Port)] = serviceManagerGroupAddr
					}
				}
			}
		}
		return serviceManagerGroup, nil
	}

	return nil, err
}

func (e *EtcdServiceManager) delGroup(groupId string) error {
	etcdConnection := etcd.EtcdConnection{}
	etcdConnection.CreateConnection(e.endpoints, e.timeout, e.username, e.password)
	groupsPrefix := e.root + ".groups"
	groupIdPrefix := groupsPrefix + "." + groupId
	err := etcdConnection.DelWithPrefix(groupIdPrefix)

	vGroup, vgErr := etcdConnection.Get(groupsPrefix)
	if vgErr == nil {
		groupIds := strings.Split(vGroup, ",")
		newGroupIds := make([]string, 0, 1)
		for _, gid := range groupIds {
			if gid != groupId {
				newGroupIds = append(newGroupIds, gid)
			}
		}
		etcdConnection.Put(groupsPrefix, strings.Join(newGroupIds, ","))
	}

	etcdConnection.Close()
	return err
}

func (e *EtcdServiceManager) saveGroupWithoutAddr(group *ServiceManagerGroup) error {
	etcdConnection := etcd.EtcdConnection{}
	etcdConnection.CreateConnection(e.endpoints, e.timeout, e.username, e.password)

	defer etcdConnection.Close()

	groupsPrefix := e.root + ".groups"
	rpcsPrefix := e.root + ".groups" + "." + group.GroupId + ".rpcs"
	vGroup, vgErr := etcdConnection.Get(groupsPrefix)
	if vgErr == nil {
		if vGroup == "" {
			etcdConnection.Put(groupsPrefix, group.GroupId)
		} else {
			groupIds := strings.Split(vGroup, ",")
			sameFlag := false
			for _, gid := range groupIds {
				if gid == group.GroupId {
					sameFlag = true
					break
				}
			}
			if !sameFlag {
				groupIds = append(groupIds, group.GroupId)
			}
			etcdConnection.Put(groupsPrefix, strings.Join(groupIds, ","))
		}
	} else {
		return vgErr
	}

	err := etcdConnection.StartTrans()
	if err != nil {
		etcdConnection.RollbackTrans()
		return err
	}
	etcdConnection.Put(groupsPrefix+"."+group.GroupId, group.GroupId)
	etcdConnection.Put(groupsPrefix+"."+group.GroupId+".timestamp", group.Timestamp)
	etcdConnection.Put(groupsPrefix+"."+group.GroupId+".protocol_info", "1:RPC_UNIX,2:RPC_TCP,3:RPC_TCP_SSL")
	rpcIds := make([]string, 0, len(group.Rpcs))
	for _, rpc := range group.Rpcs {
		etcdConnection.Put(rpcsPrefix+"."+rpc.RpcId, rpc.RpcId)
		etcdConnection.Put(rpcsPrefix+"."+rpc.RpcId+".url", rpc.Url)
		etcdConnection.Put(rpcsPrefix+"."+rpc.RpcId+".desc", rpc.Desc)
		etcdConnection.Put(rpcsPrefix+"."+rpc.RpcId+".available", strconv.FormatBool(rpc.Available))
		rpcIds = append(rpcIds, rpc.RpcId)
	}
	etcdConnection.Put(rpcsPrefix, strings.Join(rpcIds, ","))

	err1 := etcdConnection.CommitTrans()
	if err1 != nil {
		etcdConnection.RollbackTrans()
		return err1
	}

	return nil
}

func (e *EtcdServiceManager) delAddr(groupId string, ip string, port int) {
	etcdConnection := etcd.EtcdConnection{}
	etcdConnection.CreateConnection(e.endpoints, e.timeout, e.username, e.password)
	defer etcdConnection.Close()

	addrsPrefix := e.root + ".groups." + groupId + ".addrs"

	addrSamplePrefix := fmt.Sprintf("%s%s%d%s", addrsPrefix, ".["+ip+"_", port, "]")

	str_addrs, _ := etcdConnection.Get(addrsPrefix)

	addrs := strings.Split(str_addrs, ",")
	new_addrs := make([]string, 0, 1)
	for _, addr_info := range addrs {
		if addr_info != fmt.Sprintf("%s_%d", ip, port) {
			new_addrs = append(new_addrs, addr_info)
		}
	}
	etcdConnection.StartTrans()
	etcdConnection.Put(addrsPrefix, strings.Join(new_addrs, ","))
	etcdConnection.DelWithPrefix(addrSamplePrefix)
	etcdConnection.CommitTrans()

}

func (e *EtcdServiceManager) existGroup(groupId string) bool {
	etcdConnection := etcd.EtcdConnection{}
	etcdConnection.CreateConnection(e.endpoints, e.timeout, e.username, e.password)
	defer etcdConnection.Close()
	groupsPrefix := e.root + ".groups"
	v1, _ := etcdConnection.Get(groupsPrefix + "." + groupId)
	if v1 != "" {
		return true
	} else {
		return false
	}
}

func (e *EtcdServiceManager) queryGroupIds() []string {

	etcdConnection := etcd.EtcdConnection{}
	etcdConnection.CreateConnection(e.endpoints, e.timeout, e.username, e.password)
	defer etcdConnection.Close()
	groupsPrefix := e.root + ".groups"
	v1, _ := etcdConnection.Get(groupsPrefix)
	return strings.Split(v1, ",")
}

func (e *EtcdServiceManager) saveAddr(groupId string, addr *ServiceManagerGroupAddr) error {

	etcdConnection := etcd.EtcdConnection{}
	addrsPrefix := e.root + ".groups." + groupId + ".addrs"
	etcdConnection.CreateConnection(e.endpoints, e.timeout, e.username, e.password)
	defer etcdConnection.Close()

	// addr的存储格式： "192.168.1.22_9000,fe80::cbd:8cd5:ec3a:160b%en0_8000"
	// vAddr 记录所有地址列表
	vAddr, vaErr := etcdConnection.Get(addrsPrefix)
	if vaErr == nil {
		if vAddr == "" {
			vAddr = fmt.Sprintf("%s_%d", addr.Ip, addr.Port)
		} else {
			addrs := strings.Split(vAddr, ",")
			sameFlag := false
			for _, addrInfo := range addrs {
				if addrInfo == fmt.Sprintf("%s_%d", addr.Ip, addr.Port) {
					sameFlag = true
					break
				}
			}
			if !sameFlag {
				addrs = append(addrs, fmt.Sprintf("%s_%d", addr.Ip, addr.Port))
			}
			vAddr = strings.Join(addrs, ",")
		}
	}

	err := etcdConnection.StartTrans()
	if err != nil {
		etcdConnection.RollbackTrans()
		return err
	}

	addrSamplePrefix := fmt.Sprintf("%s%s%d%s", addrsPrefix, ".["+addr.Ip+"_", addr.Port, "]")

	etcdConnection.Put(addrSamplePrefix+".ip", addr.Ip)
	etcdConnection.Put(addrSamplePrefix+".iptype", addr.IpType)
	etcdConnection.Put(addrSamplePrefix+".port", strconv.Itoa(addr.Port))
	etcdConnection.Put(addrSamplePrefix+".protocol", fmt.Sprintf("%d", addr.Protocol))
	etcdConnection.Put(addrSamplePrefix+".ticket", addr.Ticket)
	etcdConnection.Put(addrSamplePrefix+".ticketforlast", addr.Ticket) // 初始化时候，票据信息都一致
	etcdConnection.Put(addrSamplePrefix+".webcontext", addr.WebContext)
	etcdConnection.Put(addrSamplePrefix+".timeout", strconv.Itoa(addr.Timeout))
	etcdConnection.Put(addrSamplePrefix+".weight", strconv.Itoa(addr.Weight))
	etcdConnection.Put(addrSamplePrefix, "add") // 设置 add 表示新增，这个必须放在最后，说明地址新增完毕。
	etcdConnection.Put(addrsPrefix, vAddr)

	err1 := etcdConnection.CommitTrans()

	return err1

}

func (e *EtcdServiceManager) queryAddr(groupId string, ip string, port int) *ServiceManagerGroupAddr {
	addrsPrefix := e.root + ".groups." + groupId + ".addrs"
	addrSamplePrefix := fmt.Sprintf("%s%s%d%s", addrsPrefix, ".["+ip+"_", port, "]")
	etcdConnection := etcd.EtcdConnection{}
	smgAddr := new(ServiceManagerGroupAddr)
	etcdConnection.CreateConnection(e.endpoints, e.timeout, e.username, e.password)
	defer etcdConnection.Close()
	kvMap, err1 := etcdConnection.GetWithPrefix(addrSamplePrefix)
	if kvMap[addrSamplePrefix] == "" {
		return nil
	}
	if err1 == nil {
		for k, v := range kvMap {
			switch k {
			case addrSamplePrefix + ".ip":
				smgAddr.Ip = v
			case addrSamplePrefix + ".iptype":
				smgAddr.IpType = v
			case addrSamplePrefix + ".port":
				port, err1 := strconv.Atoi(v)
				if err1 == nil {
					smgAddr.Port = port
				}
			case addrSamplePrefix + ".protocol":
				protocol, err1 := strconv.Atoi(v)
				if err1 == nil {
					switch protocol {
					case 1:
						smgAddr.Protocol = PROTOCOL_RPC_UNIX
					case 2:
						smgAddr.Protocol = PROTOCOL_RPC_TCP
					case 3:
						smgAddr.Protocol = PROTOCOL_RPC_TCP_SSL
					}
				}
			case addrSamplePrefix + ".ticket":
				smgAddr.Ticket = v
			case addrSamplePrefix + ".ticketforlast":
				smgAddr.TicketForLast = v
			case addrSamplePrefix + ".webcontext":
				smgAddr.WebContext = v
			case addrSamplePrefix + ".timeout":
				timeout, err1 := strconv.Atoi(v)
				if err1 == nil {
					smgAddr.Timeout = timeout
				}

			case addrSamplePrefix + ".weight":
				weight, err1 := strconv.Atoi(v)
				if err1 == nil {
					smgAddr.Weight = weight
				}

			}
		}
	}
	return smgAddr
}

func (e *EtcdServiceManager) ttlAddrByProvider(groupId string, ip string, port int, timeout int64) {
	go func() {

		for true {
			func() {
				etcdConnection := etcd.EtcdConnection{}
				etcdConnection.CreateConnection(e.endpoints, e.timeout, e.username, e.password)
				defer etcdConnection.Close()

				addrsPrefix := e.root + ".groups." + groupId + ".addrs"
				addrSamplePrefix := fmt.Sprintf("%s%s%d%s", addrsPrefix, ".["+ip+"_", port, "]")

				respTTL, err := etcdConnection.TTL_SetGrant(timeout)
				if err == nil {
					value, _ := etcdConnection.Get(addrSamplePrefix)
					//判断，如果是空串，那么就需要打新增标示
					if value == "" {
						etcdConnection.TTL_Put(addrSamplePrefix, "add", respTTL)
					} else {
						etcdConnection.TTL_Put(addrSamplePrefix, "ttl", respTTL)
					}
					time.Sleep(time.Duration(SERVICE_MS_ADDR_TTL_SLEEP_TIME) * time.Second)

				}
			}()
		}
	}()
}

func (e *EtcdServiceManager) sendTicketByProvider(groupId string, ip string, port int, rootAppname string) error {

	go func() {

		for true {

			time.Sleep(e.sendTicketDuration)
			ticket := generateTicket()
			e.sendTicketByProvider_once(groupId, ip, port, ticket, rootAppname)
		}
	}()
	return nil

}

func (e *EtcdServiceManager) sendTicketByProvider_once(groupId string, ip string, port int, ticket string, rootAppname string) {

	defer func() { // 用于捕获panic异常，不影响整个服务运行
		if err := recover(); err != nil {
			log4go.ErrorLog(message.ERR_MSSV_39025, err)
		}
	}()

	etcdConnection := etcd.EtcdConnection{}
	addrsPrefix := e.root + ".groups." + groupId + ".addrs"
	etcdConnection.CreateConnection(e.endpoints, e.timeout, e.username, e.password)
	defer etcdConnection.Close()

	addrSamplePrefix := fmt.Sprintf("%s%s%d%s", addrsPrefix, ".["+ip+"_", port, "]")

	ticketForList, _ := etcdConnection.Get(addrSamplePrefix + ".ticket") // 获取上一个票据
	// 如果上一个票据没获取到，那么就没法更新。 更新顺序是必须把ticket更新放在最后
	if ticketForList != "" {
		etcdConnection.StartTrans()
		etcdConnection.Put(addrSamplePrefix+".ticketforlast", ticketForList) // 把上一个票据存储到ticketforlast中
		etcdConnection.Put(addrSamplePrefix+".ticket", ticket)               // 存储当前票据
		etcdConnection.CommitTrans()
		serverGlobalMap[rootAppname].ticketForLast = ticketForList
		serverGlobalMap[rootAppname].ticket = ticket
	}
}

// 消费者监听程序入口，主要监听注册中心的变化。
// 程序调用流程  watchByConsumer -> watchByConsumer_this 启动监听，监听一般会捕获数据更新（watchByConsumer_putEvent）、数据删除（watchByConsumer_deleteEvent）、监听退出（watchByConsumer_exitEvent）等事件。
//             watchByConsumer -> watchByConsumer_deal_getChanDatas  把watchByConsumer_this监听到的chanWatchStore数据放到一个数组 arrayWatchStore 中
//             watchByConsumer -> watchByConsumer_deal_batchAll  针对arrayWatchStore 中的数据，定期批处理一次，进行消费者（客户端）更新。
//
// 使用watchByConsumer_this进行监听，目的是：如果异常退出，触发了退出事件（watchByConsumer_exitEvent），会自动重启开启监听。
// 使用watchByConsumer_deal_batchAll的处理方式，采用异步化，避免对注册中心产生压力。在处理时，第一步会分析重复操作步骤，去重。第二步再按顺序更新。
func (e *EtcdServiceManager) watchByConsumer() {

	e.lockForWS = new(sync.Mutex)
	e.arrayWatchStore = make([]WatchGroupInfoStore, 0, 1)

	groupsPrefix := e.root + ".groups"

	e.chanWatchStore = make(chan WatchGroupInfoStore, SERVICE_MS_CHAN_WATCH_STORE_LENGTH)

	e.watchByConsumer_this(groupsPrefix)

	go e.watchByConsumer_deal_getChanDatas()
	go e.watchByConsumer_deal_batchAll()
}

func (e *EtcdServiceManager) watchByConsumer_this(groupsPrefix string) {

	ctx, _ := context.WithCancel(context.Background())

	go func() {
		etcdConnection := etcd.EtcdConnection{}
		etcdConnection.CreateConnection(e.endpoints, e.timeout, e.username, e.password)
		defer etcdConnection.Close()
		etcdConnection.WatchWithPrefix(ctx, groupsPrefix, e.watchByConsumer_putEvent, e.watchByConsumer_deleteEvent, e.watchByConsumer_exitEvent)
	}()

}

func (e *EtcdServiceManager) watchByConsumer_deal_getChanDatas() {
	for ws := range e.chanWatchStore {
		e.lockForWS.Lock()
		e.arrayWatchStore = append(e.arrayWatchStore, ws)
		e.lockForWS.Unlock()
	}
}

func (e *EtcdServiceManager) watchByConsumer_deal_batchAll() {
	for true {

		time.Sleep(SERVICE_MS_WATCH_DEAL_BATCH_DURATION * time.Second)

		e.lockForWS.Lock()
		lengthOfDatas := len(e.arrayWatchStore)
		arrayDatas := make([]WatchGroupInfoStore, lengthOfDatas, lengthOfDatas)
		copy(arrayDatas, e.arrayWatchStore)
		e.arrayWatchStore = make([]WatchGroupInfoStore, 0, 1)
		e.lockForWS.Unlock()

		groupMap := make(map[string][]WatchGroupInfoStore)
		for _, v1 := range arrayDatas {
			if groupMap[v1.groupId] == nil {
				groupMap[v1.groupId] = make([]WatchGroupInfoStore, 0, 1)
				groupMap[v1.groupId] = append(groupMap[v1.groupId], v1)
			} else {
				groupMap[v1.groupId] = append(groupMap[v1.groupId], v1)
			}
		}

		// 11:  修改 ticket （增量，只更新该地址的票据） ， 12:  修改  timestamp (全量，全量更新该group) ，  13: 新增地址 （增量，只新增地址）
		//  22:  删除了  timestamp （全量，本地删除group）  ， 23： 删除了服务地址 （增量，只删除地址，本地删除）

		// 根据groupid 顺序存储这段时间的数据进入处理队列。每个groupid队列里， 具体更新步骤如下：
		//     如果 在 11 后面有 12，22，那么该数据无效（invalid）
		//     如果 在 11 后面有 11，13，23，还需要判断 他们的ip addr 是否一样，如果一样，那么该数据无效（invalid）
		//     如果 在 12 后面有 12、22，那么该数据无效（invalid）
		//     如果 在 13 后面有 12、22，那么该数据无效（invalid）
		//     如果 在 13 后面有 13、23 ，还需要判断 他们的ip addr 是否一样，如果一样，那么该数据无效（invalid）
		//     如果 在 22 后面有 12、22，那么该数据无效（invalid）
		//     如果 在 23 后面有 12、22，那么该数据无效（invalid）
		//     如果 在 23 后面有 11、13、23 ，还需要判断 他们的ip addr 是否一样，如果一样，那么该数据无效（invalid）

		for k1, v1 := range groupMap {
			for idx01, data01 := range v1 {
				for i := idx01 + 1; i < len(v1); i++ {
					if v1[i].watchType == 12 || v1[i].watchType == 22 {
						//     如果 在 11 后面有 12，22，那么该数据无效（invalid）
						//     如果 在 12 后面有 12、22，那么该数据无效（invalid）
						//     如果 在 13 后面有 12、22，那么该数据无效（invalid）
						//     如果 在 22 后面有 12、22，那么该数据无效（invalid）
						//     如果 在 23 后面有 12、22，那么该数据无效（invalid）
						data01.invalid = true
						break
					} else if (data01.watchType == 11 || data01.watchType == 23) &&
						(v1[i].watchType == 11 || v1[i].watchType == 13 || v1[i].watchType == 23) &&
						data01.ip == v1[i].ip && data01.port == v1[i].port {
						//     如果 在 11 后面有 11，13，23，还需要判断 他们的ip addr 是否一样，如果一样，那么该数据无效（invalid）
						//     如果 在 23 后面有 11、13、23 ，还需要判断 他们的ip addr 是否一样，如果一样，那么该数据无效（invalid）
						data01.invalid = true
						break
					} else if (data01.watchType == 13) &&
						(v1[i].watchType == 13 || v1[i].watchType == 23) &&
						data01.ip == v1[i].ip && data01.port == v1[i].port {
						//     如果 在 13 后面有 13、23 ，还需要判断 他们的ip addr 是否一样，如果一样，那么该数据无效（invalid）
						data01.invalid = true
						break
					}
				}
			}

			for _, data01 := range v1 {
				if !data01.invalid { // 非无效（有效）
					switch data01.watchType {
					case 11:
						log4go.InfoLog(message.INF_MSSV_09007, k1, data01)
						e.UpdateClientGroupAddr(data01.groupId, fmt.Sprintf("%s_%d", data01.ip, data01.port), false)
					case 12:
						log4go.InfoLog(message.INF_MSSV_09008, k1, data01)
						e.UpdateClientGroup(data01.groupId, false)
					case 13:
						log4go.InfoLog(message.INF_MSSV_09009)
						e.UpdateClientGroupAddr(data01.groupId, fmt.Sprintf("%s_%d", data01.ip, data01.port), false)
					case 22:
						log4go.InfoLog(message.INF_MSSV_09010, k1, data01)
						e.UpdateClientGroup(data01.groupId, true)
					case 23:
						log4go.InfoLog(message.INF_MSSV_09011, k1, data01)
						e.UpdateClientGroupAddr(data01.groupId, fmt.Sprintf("%s_%d", data01.ip, data01.port), true)
					}

					for _, clients := range clientsGlobalMap {
						log4go.InfoLog(message.INF_MSSV_09012, data01.groupId, data01.groupId, clients.dataMap[data01.groupId].Json())
					}

				}
			}
		}

	}
}

func (e *EtcdServiceManager) watchByConsumer_putEvent(key string, value string) {
	//fmt.Printf("putEvent   key: %s , value: %s\n", key, value)
	if strings.HasSuffix(key, "ticket") { // 判断票据更新了 或  上一个票据更新了
		// 参考格式：  e.root + ".groups." + groupId + ".addrs"+".["+addr.Ip+"_", addr.Port, "]"+".ticket"
		index0 := strings.Index(key, ".groups.") + 8
		index1 := strings.Index(key, ".addrs.[")
		index2 := index1 + 8
		index3 := strings.Index(key, "].ticket")
		groupId := key[index0:index1]
		addr := key[index2:index3]
		idx_ipAndPort := strings.LastIndex(addr, "_")
		ipAndPort := []string{addr[0:idx_ipAndPort], addr[(idx_ipAndPort + 1):len(addr)]}
		//ipAndPort := strings.Split(addr, "_")
		port, _ := strconv.Atoi(ipAndPort[1])
		w := WatchGroupInfoStore{groupId: groupId, ip: ipAndPort[0], port: port, key: key, value: value, watchType: 11}
		e.chanWatchStore <- w
	} else if strings.HasSuffix(key, "timestamp") { // 判断时间戳更新，说明 groupId更新了
		// 参考格式：  e.root + ".groups." + groupId + ".timestamp"
		index0 := strings.Index(key, ".groups.") + 8
		index1 := len(key) - 9 - 1
		groupId := key[index0:index1]
		w := WatchGroupInfoStore{groupId: groupId, key: key, value: value, watchType: 12}
		e.chanWatchStore <- w
		//e.UpdateClientGroup(GetClientServiceGroupMap(), groupId, false)
	} else if strings.HasSuffix(key, "]") && value == "add" { // 判断是不是新增地址
		// 参考格式：  e.root + ".groups." + groupId + ".addrs"+".["+addr.Ip+"_", addr.Port, "]"
		index0 := strings.Index(key, ".groups.") + 8
		index1 := strings.Index(key, ".addrs.[")
		index2 := index1 + 8
		index3 := len(key) - 1
		groupId := key[index0:index1]
		addr := key[index2:index3]

		idx_ipAndPort := strings.LastIndex(addr, "_")
		ipAndPort := []string{addr[0:idx_ipAndPort], addr[(idx_ipAndPort + 1):len(addr)]}
		//ipAndPort := strings.Split(addr, "_")
		port, _ := strconv.Atoi(ipAndPort[1])
		w := WatchGroupInfoStore{groupId: groupId, ip: ipAndPort[0], port: port, key: key, value: value, watchType: 13}
		e.chanWatchStore <- w
	}

}

func (e *EtcdServiceManager) watchByConsumer_deleteEvent(key string, value string) {
	//fmt.Printf("deleteEvent   key: %s , value: %s\n", key, value)
	if strings.HasSuffix(key, "timestamp") { // 判断时间戳被删除，说明 groupId不存在了
		// 参考格式：  e.root + ".groups." + groupId + ".timestamp"
		index0 := strings.Index(key, ".groups.") + 8
		index1 := strings.Index(key, ".timestamp")
		groupId := key[index0:index1]
		w := WatchGroupInfoStore{groupId: groupId, key: key, value: value, watchType: 22}
		e.chanWatchStore <- w
	} else if strings.HasSuffix(key, "]") { // 判断是不是addr被删除了，需要及时更新地址信息
		// 参考格式：  e.root + ".groups." + groupId + ".addrs"+".["+addr.Ip+"_", addr.Port, "]"
		index0 := strings.Index(key, ".groups.") + 8
		index1 := strings.Index(key, ".addrs.[")
		index2 := index1 + 8
		index3 := len(key) - 1
		groupId := key[index0:index1]
		addr := key[index2:index3]
		idx_ipAndPort := strings.LastIndex(addr, "_")
		ipAndPort := []string{addr[0:idx_ipAndPort], addr[(idx_ipAndPort + 1):len(addr)]}
		//ipAndPort := strings.Split(addr, "_")
		port, _ := strconv.Atoi(ipAndPort[1])
		w := WatchGroupInfoStore{groupId: groupId, ip: ipAndPort[0], port: port, key: key, value: value, watchType: 23}
		e.chanWatchStore <- w

	}
}

// 监听： 退出事件，那么重新启动，并且
func (e *EtcdServiceManager) watchByConsumer_exitEvent(path string, err error) {
	log4go.WarnLog(message.WAR_MSSV_69003, path)
	time.Sleep(SERVICE_MS_WATCH_EXIT_SLEEP * time.Second)
	e.watchByConsumer_this(path)
	for _, clients := range clientsGlobalMap {
		for _, gid := range clients.configMap[e.configName].Client[0].GroupIds { // 重新更新这个group信息
			e.UpdateClientGroup(gid, false)
		}
	}
}

type WatchGroupInfoStore struct {
	groupId   string
	ip        string
	port      int
	key       string
	value     string
	watchType int // 11:  修改 ticket （增量，只更新该地址的票据） ， 12:  修改  timestamp (全量，全量更新该group) ，  13: 判断是不是新增地址 （增量，只新增地址）
	//  22:  删除了  timestamp （全量，本地删除group）  ， 23： 删除了服务地址 （增量，只删除地址，本地删除）
	invalid bool
}
