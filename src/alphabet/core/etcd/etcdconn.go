// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package etcd

import (
	"alphabet/log4go"
	"alphabet/log4go/message"
	"context"
	"errors"
	"time"

	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/connectivity"
)

//
// ETCD 的操作类，一个EtcdClient实例对象，只能用于创建一个连接对象处理
// 1、连接的创建和管理
//     CreateConnection()
//     ....
//     Close()
// 2、支持事务处理方式，他的结构是：
//     StartTrans()
//     Put()
//     ...
//     Delete()
//     ...
//     CommitTrans() / RollbackTrans()
//
// 3. 支撑TTL 数据到期设置
//     TTL_SetGrant(...)  设置到期时间，返回一个TTL对象
//     TTL_Put(...)   把TTL对象应用到具体的数据上
//    如果想在过程中保持数据有效期，
//     TTL_SetKeepAliveOnce()   每次调用一次，就把到期时间重置，如果超过一个到期时间调用，实际上无效。
//    如果想知道到期时间情况
//     TTL_GetTimeToLive()    能够获取到期时间设置 和 实际剩余的有效时间。
//
// 4. 支持监听，感知etcd数据变化，支持 增删改的变化。
//     WatchWithPrefix(...)    监听操作，采用 go func(){}() 协程执行
//
// 5. 支持Context的处理
//     ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)  // 设置5秒超时
//     SetContext(ctx)
//     ....         // 支持：StartTrans() ， Put() ，Del() , Get() , TTL_SetGrant,TTL_Put ,TTL_SetKeepAliveOnce
//
type EtcdConnection struct {
	cli           *clientv3.Client
	ctx           context.Context
	txn           clientv3.Txn
	useTrans      bool
	datasForTrans []dataStoreForTrans
}

type dataStoreForTrans struct {
	oper    int // 操作类型， 0 表示PUT，1表示DELETE ，2表示DELETE所有相似的前缀
	key     string
	value   string
	respTTL *clientv3.LeaseGrantResponse
}

func (ec *EtcdConnection) SetContext(ctx0 context.Context) {
	ec.ctx = ctx0
}

func (ec *EtcdConnection) GetContext() context.Context {
	if ec.ctx == nil {
		ec.ctx = context.TODO()
	}
	return ec.ctx
}

func (ec *EtcdConnection) initTrans() {
	ec.datasForTrans = nil
	ec.datasForTrans = make([]dataStoreForTrans, 0, 1)
	ec.useTrans = false
}

// 创建连接
func (ec *EtcdConnection) CreateConnection(endpoints []string, dialTimeout time.Duration, username string, passwd string) {
	var err error
	ec.initTrans()
	if ec.cli != nil {
		ec.cli.Close()
	}

	ec.cli, err = clientv3.New(
		clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: dialTimeout,
			Username:    username,
			Password:    passwd,
		})
	if err != nil {
		log4go.ErrorLog(message.ERR_CORE_39010, err)
	}

	log4go.DebugLog(message.DEG_CORE_79001, ec.GetState())

}

// 关闭连接
func (ec *EtcdConnection) Close() {
	if ec.useTrans { // 如果事务还是使用状态，就做commit提交
		ec.CommitTrans()
	}
	if ec.cli != nil {
		err := ec.cli.Close()
		if err != nil {
			log4go.ErrorLog(message.ERR_CORE_39011, err)
		}
		ec.cli = nil
	}
	ec.initTrans()
	log4go.DebugLog(message.DEG_CORE_79001, ec.GetState())
}

// 判断连接是否正常，
func (ec *EtcdConnection) IsOpened() bool {
	isOpenState := false
	if ec.cli != nil {
		conn := ec.cli.ActiveConnection()
		if conn.GetState() == connectivity.Ready {
			isOpenState = true
		}
	}
	return isOpenState
}

// 获取连接状态
func (ec *EtcdConnection) GetState() string {
	state := ""
	if ec.cli == nil {
		state = "Conn_Object_IsNil"
	} else {
		state = ec.cli.ActiveConnection().GetState().String()
	}
	return state
}

// 开启事务
func (ec *EtcdConnection) StartTrans() (err error) {
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		err = errors.New(message.ERR_CORE_39012.String())
	} else {
		ec.useTrans = true
		ec.txn = ec.cli.Txn(ec.GetContext())
	}
	return
}

// 事务提交，实际进行数据写入，（去重复key，只执行最后一个）
func (ec *EtcdConnection) CommitTrans() (err error) {
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		err = errors.New(message.ERR_CORE_39012.String())
	} else {
		if len(ec.datasForTrans) > 0 {
			ops := make([]clientv3.Op, 0, 10)
			keyDuplicateMap := make(map[string]int)
			for _, data := range ec.datasForTrans { // 用于处理重复key的问题，默认最后一个key是有效的
				num := keyDuplicateMap[data.key]
				keyDuplicateMap[data.key] = num + 1
			}
			for _, data := range ec.datasForTrans {
				if keyDuplicateMap[data.key] > 1 {
					keyDuplicateMap[data.key] = keyDuplicateMap[data.key] - 1
					log4go.DebugLog(message.DEG_CORE_79002, data.key, data.value, data.oper, data.respTTL)
					continue
				} else {
					keyDuplicateMap[data.key] = keyDuplicateMap[data.key] - 1
					log4go.DebugLog(message.DEG_CORE_79003, data.key, data.value, data.oper, data.respTTL)
				}

				if data.oper == 0 { // put 操作
					if data.respTTL == nil {
						ops = append(ops, clientv3.OpPut(data.key, data.value))
					} else {
						ops = append(ops, clientv3.OpPut(data.key, data.value, clientv3.WithLease(data.respTTL.ID)))
					}
				} else if data.oper == 1 { // delete操作
					ops = append(ops, clientv3.OpDelete(data.key))
				} else if data.oper == 2 {
					ops = append(ops, clientv3.OpDelete(data.key, clientv3.WithPrefix()))
				}
			}
			//log4go.Debug(ops[0])
			//ec.txn.Then(ops[0], ops[1])

			//ec.txn.Then(ops[0:len(ops)]...)

			ec.txn.Then(ops...)
		}
		_, err0 := ec.txn.Commit()
		if err0 != nil {
			err = err0
			if err0.Error() == "etcdserver: duplicate key given in txn request" { // 存在重复的key
				log4go.ErrorLog(message.ERR_CORE_39013, err0.Error())
			} else {
				log4go.ErrorLog(message.ERR_CORE_39013, err.Error())
			}
		}
	}
	ec.initTrans()
	return
}

// 事务回滚
func (ec *EtcdConnection) RollbackTrans() (err error) {
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		err = errors.New(message.ERR_CORE_39012.String())
	} else {
		ec.txn.Commit() // 只提交不做任何处理
	}
	ec.initTrans() // 数据全部清除
	return
}

// 写入key，value数据
func (ec *EtcdConnection) Put(key string, value string) error {
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		return errors.New(message.ERR_CORE_39012.String())
	}
	if ec.useTrans {
		ec.datasForTrans = append(ec.datasForTrans, dataStoreForTrans{oper: 0, key: key, value: value, respTTL: nil})
	} else {
		_, err := ec.cli.Put(ec.GetContext(), key, value)
		if err != nil {
			log4go.ErrorLog(message.ERR_CORE_39017, err.Error())
			return err
		}
	}
	return nil
}

// 删除一个key及其数据
func (ec *EtcdConnection) Del(key string) error {
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		return errors.New(message.ERR_CORE_39012.String())
	}
	if ec.useTrans {
		ec.datasForTrans = append(ec.datasForTrans, dataStoreForTrans{oper: 1, key: key, value: "", respTTL: nil})
	} else {
		_, err := ec.cli.Delete(ec.GetContext(), key)
		if err != nil {
			log4go.ErrorLog(message.ERR_CORE_39018, err.Error())
			return err
		}
	}
	return nil
}

// 根据key的前缀，删除所有匹配的key及其数据
func (ec *EtcdConnection) DelWithPrefix(key string) error {
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		return errors.New(message.ERR_CORE_39012.String())
	}
	if ec.useTrans {
		ec.datasForTrans = append(ec.datasForTrans, dataStoreForTrans{oper: 2, key: key, value: "", respTTL: nil})
	} else {
		_, err := ec.cli.Delete(ec.GetContext(), key, clientv3.WithPrefix())
		if err != nil {
			log4go.ErrorLog(message.ERR_CORE_39018, err.Error())
			return err
		}
	}
	return nil
}

// 根据key获取value数据
// @return value  返回数据，如果没有，就为空‘’
func (ec *EtcdConnection) Get(key string) (value string, err error) {
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		return "", errors.New(message.ERR_CORE_39012.String())
	}
	resp, err := ec.cli.Get(ec.GetContext(), key)
	if err == nil {
		if resp.Count > 0 {
			value = string(resp.Kvs[0].Value)
		}
	}
	return
}

// 根据key的前缀，获取value数据
// @return  kvMaps 返回数据，如果报错或者没有数据，为nil
func (ec *EtcdConnection) GetWithPrefix(key string) (kvMaps map[string]string, err error) {
	kvMaps = nil
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		return kvMaps, errors.New(message.ERR_CORE_39012.String())
	}
	resp, err := ec.cli.Get(ec.GetContext(), key, clientv3.WithPrefix())

	if err == nil {
		if resp.Count > 0 {
			kvMaps = make(map[string]string)
			for _, datas := range resp.Kvs {
				kvMaps[string(datas.Key)] = string(datas.Value)
			}
		}
	}
	return
}

// 获取一批数据的创建版本信息
// 通过这种方式实现分布式锁
func (ec *EtcdConnection) GetCreateVersionWithPrefix(key string) (kvMaps map[string]int64, err error) {
	kvMaps = nil
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		return kvMaps, errors.New(message.ERR_CORE_39012.String())
	}
	resp, err := ec.cli.Get(ec.GetContext(), key, clientv3.WithPrefix())

	if err == nil {
		if resp.Count > 0 {
			kvMaps = make(map[string]int64)
			for _, datas := range resp.Kvs {
				kvMaps[string(datas.Key)] = datas.CreateRevision
			}
		}
	}
	return
}

// 设置多少秒以后到期,返回TTL标示
// @param  second  设置到期时间，单位：秒
// @return 返回TTL对象，用于后续的设置中。
func (ec *EtcdConnection) TTL_SetGrant(second int64) (*clientv3.LeaseGrantResponse, error) {
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		return nil, errors.New(message.ERR_CORE_39012.String())
	}
	espTTL, err := ec.cli.Grant(ec.GetContext(), second)
	if err != nil {
		log4go.ErrorLog(message.ERR_CORE_39019, err.Error())
	}
	return espTTL, err
}

// 在写入数据时候，设置TTL标示，到期后数据被清除
func (ec *EtcdConnection) TTL_Put(key string, value string, respTTL *clientv3.LeaseGrantResponse) error {
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		return errors.New(message.ERR_CORE_39012.String())
	}
	if ec.useTrans {
		ec.datasForTrans = append(ec.datasForTrans, dataStoreForTrans{oper: 0, key: key, value: value, respTTL: respTTL})
	} else {
		_, err := ec.cli.Put(ec.GetContext(), key, value, clientv3.WithLease(respTTL.ID))
		return err
	}
	return nil
}

// 保持数据有效一次，也就是把到期时间重置
func (ec *EtcdConnection) TTL_SetKeepAliveOnce(respTTL *clientv3.LeaseGrantResponse) error {
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		return errors.New(message.ERR_CORE_39012.String())
	} else {
		_, err := ec.cli.KeepAliveOnce(ec.GetContext(), respTTL.ID)
		return err
	}

}

// 获取当前有效时间，
// 返回值： 第一个，表示设置的到期时间，单位秒
//         第二个，实际现在还有多少秒到期
func (ec *EtcdConnection) TTL_GetTimeToLive(respTTL *clientv3.LeaseGrantResponse) (int64, int64, error) {
	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		return 0, 0, errors.New(message.ERR_CORE_39012.String())
	} else {
		timeResp, err := ec.cli.TimeToLive(ec.GetContext(), respTTL.ID)
		return timeResp.GrantedTTL, timeResp.TTL, err
	}
}

// 监听数据变化情况，在协程中运行
// 需要定义 put 事件和delete事件处理的方法
// 返回操作函数：取消函数，如果执行他，那么监听就被取消。
/*
func (ec *EtcdConnection) WatchWithPrefix(prefixPath string,
	putEventFunc func(key string, value string),
	deleteEventFunc func(key string, value string),
	exitFunc func(path string, err error)) (context.CancelFunc, error) {

	if !ec.IsOpened() {
		return nil, errors.New("Etcd Connection maybe closed . It's not working!")
	}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		w := ec.cli.Watch(ctx, prefixPath, clientv3.WithPrefix())

		for action := range w {
			errWatch := action.Err()
			if errWatch != nil { //如果报错，就直接退出
				log4go.Error(errWatch)
				cancel()
				exitFunc(prefixPath, errWatch)
				return
			}
			for _, ev := range action.Events {
				if ev.Type == clientv3.EventTypePut {
					putEventFunc(string(ev.Kv.Key), string(ev.Kv.Value))
				} else if ev.Type == clientv3.EventTypeDelete {
					deleteEventFunc(string(ev.Kv.Key), string(ev.Kv.Value))
				}
			}
		}

		exitFunc(prefixPath, errors.New("Exit Error"))
	}()
	return cancel, nil
}
*/

// 监听数据变化情况，在协程中运行
// 需要定义 put 事件和delete事件处理的方法
// 返回操作函数：取消函数，如果执行他，那么监听就被取消，会触发exitFunc函数。
/*
func (ec *EtcdConnection) Watch(path string,
	putEventFunc func(key string, value string),
	deleteEventFunc func(key string, value string),
	exitFunc func(path string, err error)) (context.CancelFunc, error) {

	if !ec.IsOpened() {
		return nil, errors.New("Etcd Connection maybe closed . It's not working!")
	}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {

		w := ec.cli.Watch(ctx, path)

		for action := range w {
			fmt.Println("watch:  action.Canceled = " + strconv.FormatBool(action.Canceled))
			errWatch := action.Err()
			if errWatch != nil { //如果报错，就直接退出
				log4go.Error(errWatch)
				cancel()
				exitFunc(path, errWatch)
				return
			}
			for _, ev := range action.Events {
				if ev.Type == clientv3.EventTypePut {
					putEventFunc(string(ev.Kv.Key), string(ev.Kv.Value))
				} else if ev.Type == clientv3.EventTypeDelete {
					deleteEventFunc(string(ev.Kv.Key), string(ev.Kv.Value))
				}
			}
		}
		exitFunc(path, errors.New("Exit Error"))
	}()
	return cancel, nil
}
*/

// 监听数据变化情况，在协程中运行
// ctx, cancel := context.WithCancel(context.Background())
// 需要定义 put 事件和delete事件处理的方法
// 返回操作函数：取消函数，如果执行他，那么监听就被取消，会触发exitFunc函数。
func (ec *EtcdConnection) Watch(ctx context.Context, path string,
	putEventFunc func(key string, value string),
	deleteEventFunc func(key string, value string),
	exitFunc func(path string, err error)) error {

	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		err := errors.New(message.ERR_CORE_39012.String())

		exitFunc(path, err)
		return err
	}

	w := ec.cli.Watch(ctx, path)

	for action := range w {
		errWatch := action.Err()
		if errWatch != nil { //如果报错，就直接退出
			log4go.ErrorLog(message.ERR_CORE_39020, errWatch.Error())
			exitFunc(path, errWatch)
			return errWatch
		}
		for _, ev := range action.Events {
			if ev.Type == clientv3.EventTypePut {
				putEventFunc(string(ev.Kv.Key), string(ev.Kv.Value))
			} else if ev.Type == clientv3.EventTypeDelete {
				deleteEventFunc(string(ev.Kv.Key), string(ev.Kv.Value))
			}
		}
	}

	err := errors.New("Exit Error")
	exitFunc(path, err)
	return err
}

// 监听数据变化情况，在协程中运行
// 需要定义 put 事件和delete事件处理的方法
// 返回操作函数：取消函数，如果执行他，那么监听就被取消。
func (ec *EtcdConnection) WatchWithPrefix(ctx context.Context, prefixPath string,
	putEventFunc func(key string, value string),
	deleteEventFunc func(key string, value string),
	exitFunc func(path string, err error)) error {

	if !ec.IsOpened() {
		log4go.ErrorLog(message.ERR_CORE_39012)
		err := errors.New(message.ERR_CORE_39012.String())
		exitFunc(prefixPath, err)
		return err
	}

	w := ec.cli.Watch(ctx, prefixPath, clientv3.WithPrefix())

	for action := range w {
		errWatch := action.Err()
		if errWatch != nil { //如果报错，就直接退出
			log4go.ErrorLog(message.ERR_CORE_39020, errWatch.Error())
			exitFunc(prefixPath, errWatch)
			return errWatch
		}
		for _, ev := range action.Events {
			if ev.Type == clientv3.EventTypePut {
				putEventFunc(string(ev.Kv.Key), string(ev.Kv.Value))
			} else if ev.Type == clientv3.EventTypeDelete {
				deleteEventFunc(string(ev.Kv.Key), string(ev.Kv.Value))
			}
		}
	}

	exitFunc(prefixPath, errors.New("Exit Error"))

	return nil
}
