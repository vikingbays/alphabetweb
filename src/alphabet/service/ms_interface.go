// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*
定义抽象服务端管理模块的接口
*/

package service

import (
	"encoding/json"
)

/*
  服务注册环节：
	服务提供者：provider ，服务消费者：consumer ，服务管理注册中心， register
  1、provider 启动时：
	  [RegGroupAndRpc]    首先调用RegGroupAndRpc， 注册group 和 rpc，如果已经存在，就不要注册了
		[RegAddr]           然后调用RegAddr，注册Addr
	2、provider 运行过程中，不断汇报给register 状态信息，
	  [ttlAddrByProvider] 如果状态失效（provider停止，或者网络不通），register将移除此provider的Addr： ttlAddrByProvider
    [sendTicketByProvider] 定期更新票据
	3、consumer 启动时或者重新连接，连接到register，监听provider，获取注册信息
	  [GetGroupById]         如果没有group信息，那么说明是新启动，直接获取group信息
	4、consumer 运行过程中，监控状态
	   监控Addr的变化，包括票据的变化，更新Addr信息。
		[watchVersionWithTimestampByConsumer] 监控group的Version变化，如果变化，说明group被更新了，更新group信息。
		[watchTicketByConsumer]  监控group的Ticket变化


*/
type IServiceManager interface {

	// 在服务提供者启动的时候调用，用于注册group信息和RPC信息，不包含Addr
	// 第一步：先查询当前注册中心（register）该group 的信息：
	//    调用方法： queryGroupById
	//      如果 没有找到该group  ，那么就进行group信息和RPC信息的注册
	//      如果 发现没有服务地址（Addrs），那么就进行group信息和RPC信息的注册（更新group信息）
	//      如果 有一个服务地址存在，那么就不需要注册。
	// 第二步：如果需要注册，就进行注册。
	//    调用方法：delGroup(groupId string)
	//    调用方法：saveGroupWithoutAddr(group *ServiceManagerGroup)
	RegGroupAndRpc(group *ServiceManagerGroup)

	// 添加一个服务地址ip和端口
	// 调用方法： saveAddr
	RegAddr(groupId string, addr *ServiceManagerGroupAddr, rootAppname string)

	// 获取RPC配置信息
	GetGroupById(groupId string) *ServiceManagerGroup

	// 获取所有的groupId信息
	GetGroupIds() []string
	// 返回可以调用的服务
	//GetRpcURIInfo(groupId string, rpcId string) (protocol ProtocolType, ip string, port int, url string, timeout int)

	existGroup(groupId string) bool // 判断这个group 是否存在
	queryGroupById(groupId string) (*ServiceManagerGroup, error)
	queryGroupIds() []string
	delGroup(groupId string) error
	saveGroupWithoutAddr(group *ServiceManagerGroup) error // 不包含addr的所有存储

	saveAddr(groupId string, addr *ServiceManagerGroupAddr) error
	delAddr(groupId string, ip string, port int)

	queryAddr(groupId string, ip string, port int) *ServiceManagerGroupAddr

	ttlAddrByProvider(groupId string, ip string, port int, timeout int64)               // provider 发起地址存活监控
	sendTicketByProvider(groupId string, ip string, port int, rootAppname string) error // provider 定期发送票据
	watchByConsumer()                                                                   // consumer (1)、监听 Group的Timestamp有没有变化，（2）、监听 Group的Ticket有没有变化

	UpdateClientGroup(groupId string, delFlag bool)                  // 重新更新客户端的groupid
	UpdateClientGroupAddr(groupId string, addr string, delFlag bool) //重新更新客户端的某一个group的某一地址
}

// Group信息
type ServiceManagerGroup struct {
	GroupId     string                              `json:"groupId"` //  groupid唯一标示，一旦创建不能修改
	Addrs       map[string]*ServiceManagerGroupAddr `json:"addrs"`
	Rpcs        map[string]*ServiceManagerGroupRpc  `json:"rpcs"`
	Timestamp   string                              `json:"timestamp"`   //  时间戳，用于识别版本信息
	MaxPoolSize int                                 `json:"maxPoolSize"` // 针对客户端的每个addr的最大连接数
	ReqPerConn  int                                 `json:"reqPerConn"`  // 针对每个连接的并发请求数
}

func (smg *ServiceManagerGroup) Json() string {
	if smg == nil {
		return ""
	} else {
		b, _ := json.Marshal(smg)
		return string(b)
	}
}

// 记录IP和端口
type ServiceManagerGroupAddr struct {
	IpType        string       `json:"ipType"`        // IP地址类型： ipv4 或 ipv6
	Ip            string       `json:"ip"`            // IP地址
	Port          int          `json:"port"`          // 端口
	Weight        int          `json:"weight"`        // 权重，默认设置10，数据范围在【1到10之间】，用于设置访问占比 。 默认是10
	Ticket        string       `json:"ticket"`        // 票据信息，每各30分钟变一次。
	TicketForLast string       `json:"ticketForLast"` // 上一个票据信息。  默认情况下 使用 TicketForLast 和 Ticket 都有效
	Timeout       int          `json:"timeout"`       //  超时时间，单位是：秒 ，默认是30秒
	WebContext    string       `json:"webContext"`    // web 上下文名称
	Protocol      ProtocolType `json:"protocol"`      // 使用的协议，当前主要的协议是：  rpc_unix ,rpc_tcp ,rpc_tcp_ssl
	//	Group         *ServiceManagerGroup //  指定group信息
}

func (smga *ServiceManagerGroupAddr) Json() string {
	if smga == nil {
		return ""
	} else {
		b, _ := json.Marshal(smga)
		return string(b)
	}
}

func (smga *ServiceManagerGroupAddr) GetProtocolInfo() string {
	switch smga.Protocol {
	case PROTOCOL_RPC_UNIX:
		return "rpc_unix"
	case PROTOCOL_RPC_TCP:
		return "rpc_tcp"
	case PROTOCOL_RPC_TCP_SSL:
		return "rpc_tcp_ssl"
	}
	return ""
}

// 记录每个服务的信息
type ServiceManagerGroupRpc struct {
	RpcId     string //  服务唯一标示
	Url       string //  服务实际的url地址，不包含 IP和Port
	Desc      string //  服务描述
	Available bool   //  是否可用
	//	Group     *ServiceManagerGroup //  指定group信息
}

func (smgr *ServiceManagerGroupRpc) Json() string {
	if smgr == nil {
		return ""
	} else {
		b, _ := json.Marshal(smgr)
		return string(b)
	}
}

type ProtocolType int // 定义三个协议类型：rpc_unix ,rpc_tcp ,rpc_tcp_ssl

const (
	PROTOCOL_RPC_UNIX ProtocolType = 1 + iota
	PROTOCOL_RPC_TCP
	PROTOCOL_RPC_TCP_SSL
)

const (
	GROUP_DEFAULT_TIMEOUT int = 30 // 缺省 30秒
)
