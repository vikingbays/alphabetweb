// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package utils

import (
	//	"alphabet/log4go"

	"container/list"
	"sync"
)

/*
AbstractPool 是一个抽象对象，用于构建对象池，例如：数据库连接池等。

使用：以创建数据库连接池举例

1、首先，需要实现IObjectFactory接口
  type IObjectFactory interface {
		// 创建对象，例如：数据库连接
		Create() interface{}

		// 验证对象是否有效（或运行正常）
		Valid(obj interface{}) bool

		// AbstractPool的Release方法，释放对象前调用
		ReleaseStart(obj interface{})

		// AbstractPool的Release方法，释放对象后调用
		ReleaseEnd(obj interface{})

		// AbstractPool的Get方法，获取对象前调用
		GetStart(obj interface{})

		// AbstractPool的Get方法，获取对象后调用
		GetEnd(obj interface{})
  }

  实现：
	type DBObjectFactory struct {
	  Name           string
	  DriverName     string
	  DataSourceName string
	}

	func (c *DBObjectFactory) Create() interface{} {}
	func (c *DBObjectFactory) Valid(obj interface{}) bool{}
	func (c *DBObjectFactory) ReleaseStart(obj interface{}){}
	func (c *DBObjectFactory) ReleaseEnd(obj interface{}){}
	func (c *DBObjectFactory) GetStart(obj interface{}){}
	func (c *DBObjectFactory) GetEnd(obj interface{}){}


2、创建数据库连接池对象，采用golang语言特有的嵌套struct的方式继承AbstractPool，另外可以加入特有属性
  type ConnectionPool struct{
    AbstractPool
  }

3、为ConnectionPool创建初始化函数
  func NewConnectionPool(maxPoolSize int,Name string ,DriverName string,DataSourceName string)  {
    c.MaxPoolSize = maxPoolSize
    c.UsedObjectCount = 0
    c.Lock = new(sync.Mutex)
    c.Cond = sync.NewCond(c.Lock)
	c.ObjectStack = utils.NewStack()
	c.ObjectFactory = &DBObjectFactory{name, driverName, dataSourceName}
	c.TryTimes = 3
	c.CreateObjects(maxPoolSize)
  }

4、使用
  初始化工作，创建连接
  NewConnectionPool(10,"PG1","postgres","postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

*/

type AbstractPool2 struct {
	Parent            interface{}
	Lock              *sync.Mutex
	Cond              *sync.Cond
	MaxPoolSize       int // 连接池大小
	ReqPerConn        int // 每个连接可接受的请求数
	maxReqCount       int // 总并发数是 MaxPoolSize * ReqPerConn
	ObjectStack       *Stack
	usedObjectCounter *PoolCounter
	ObjectFactory     IConnectionObjectFactory
	TryTimes          int
	flagComplete      bool
}

type IConnectionObjectFactory interface {
	Create() IConnectionObject // 不能为nil
	Valid(obj IConnectionObject) bool

	ReleaseStart(obj IConnectionObject)
	ReleaseEnd(obj IConnectionObject)

	GetStart()
	GetEnd(obj IConnectionObject)
}

type IConnectionObject interface {
	ReConn()
}

func (pool *AbstractPool2) Init() {
	pool.flagComplete = false
	pool.MaxPoolSize = 10
	pool.ReqPerConn = 1
	pool.usedObjectCounter = new(PoolCounter)
	pool.Lock = new(sync.Mutex)
	pool.Cond = sync.NewCond(pool.Lock)
	pool.ObjectStack = NewStack()
}

func (pool *AbstractPool2) CreateObjects() (connCount int, reqCount int) {
	pool.Lock.Lock()
	defer pool.Lock.Unlock()

	pool.maxReqCount = pool.MaxPoolSize * pool.ReqPerConn
	objectList := list.New()
	for i := 0; i < pool.MaxPoolSize; i++ { // 第一批是初始化的对象
		obj := pool.ObjectFactory.Create() // 不能为nil
		objectList.PushBack(obj)
		pool.ObjectStack.Push(obj)
	}

	for j := 1; j < pool.ReqPerConn; j++ { // 第二批开始，都是第一批的对象指针
		iterator_poolSize := 0
		for e := objectList.Front(); e != nil && iterator_poolSize < pool.MaxPoolSize; e = e.Next() {
			pool.ObjectStack.Push(e.Value)
		}
	}

	connCount = objectList.Len()
	reqCount = pool.ObjectStack.Len()

	objectList = nil

	pool.flagComplete = true

	return
}

/*
 * 从对象池中获取对象
 * 如果获取的对象是nil，表示没有从连接池获取成功。无需再释放
 *
 */
func (pool *AbstractPool2) Get() IConnectionObject {
	var obj IConnectionObject
	pool.Lock.Lock()
	defer pool.Lock.Unlock()
	if !pool.flagComplete {
		return nil
	}

	pool.ObjectFactory.GetStart()

	var flag bool

	if pool.ObjectStack.Len() == 0 {
		pool.Cond.Wait()
	}

	if pool.ObjectStack.Len() > 0 {
		o := pool.ObjectStack.Pop()
		pool.usedObjectCounter.add()
		if o != nil {
			obj = o.(IConnectionObject)
		} else {
			return nil
		}

		flag = pool.ObjectFactory.Valid(obj)
		if !flag { //判断获取对象是否有效 , flag=true 是有效的。
			obj.ReConn()
		}
	}

	pool.ObjectFactory.GetEnd(obj)

	return obj
}

/**
 * 使用完，释放会对象池
 */
func (pool *AbstractPool2) Release(obj IConnectionObject) {
	pool.Lock.Lock()
	defer func() {
		pool.Lock.Unlock()
		pool.Cond.Signal() // 每次只通知一个等待 ， 那么只会触发一个运行。不需要在锁期间执行。
	}()

	if !pool.flagComplete {
		return
	}

	pool.ObjectFactory.ReleaseStart(obj)
	if pool.usedObjectCounter.get() > 0 {
		pool.usedObjectCounter.sub()
	}
	pool.ObjectStack.Push(obj)

	pool.ObjectFactory.ReleaseEnd(obj)

}
