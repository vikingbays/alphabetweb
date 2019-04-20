// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package utils

import (
	//	"alphabet/log4go"

	"alphabet/log4go"
	"alphabet/log4go/message"
	"sync"
	"time"
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

type IObjectFactory interface {
	Create() interface{}
	Valid(obj interface{}) bool

	ReleaseStart(obj interface{})
	ReleaseEnd(obj interface{})

	GetStart()
	GetEnd(obj interface{})
}

type AbstractPool struct {
	Parent            interface{}
	Lock              *sync.Mutex
	Cond              *sync.Cond
	MaxPoolSize       int
	ObjectStack       *Stack
	usedObjectCounter *PoolCounter
	ObjectFactory     IObjectFactory
	TryTimes          int
}

type PoolCounter struct {
	count int
}

func (pc *PoolCounter) add() {
	pc.count++
}

func (pc *PoolCounter) sub() {
	pc.count--
}

func (pc *PoolCounter) get() int {
	return pc.count
}

func (pool *AbstractPool) Init() {
	pool.MaxPoolSize = 100
	pool.usedObjectCounter = new(PoolCounter)
	pool.Lock = new(sync.Mutex)
	pool.Cond = sync.NewCond(pool.Lock)
	pool.ObjectStack = NewStack()
}

func (pool *AbstractPool) CreateObjects(countOfObject int) int {
	pool.Lock.Lock()
	defer pool.Lock.Unlock()
	return pool._CreateObjects_NonLock(countOfObject)
}

func (pool *AbstractPool) _CreateObjects_NonLock(countOfObject int) int {
	num := 0
	for i := 0; i < countOfObject; i++ {
		obj := pool.ObjectFactory.Create()
		if obj != nil {
			pool.ObjectStack.Push(obj)
			num = num + 1
		}
	}
	return num
}

/*
 * 从对象池中获取对象
 * 如果获取的对象是nil，表示没有从连接池获取成功。无需再释放
 *
 */
func (pool *AbstractPool) Get() interface{} {
	var obj interface{}
	pool.Lock.Lock()
	defer pool.Lock.Unlock()
	pool.ObjectFactory.GetStart()

	var flag bool
	var tryTimesCurr int = 0

	for !flag && tryTimesCurr <= pool.TryTimes {
		if pool.ObjectStack.Len() == 0 {
			if pool.MaxPoolSize > pool.usedObjectCounter.get() { // 判断实际使用数如果少于容器数量，那么就需要创建对象
				pool._CreateObjects_NonLock(pool.MaxPoolSize - pool.usedObjectCounter.get())
			} else { // 如果对象已经用完，那么就需要等待
				pool.Cond.Wait()
				continue
			}
		}
		objFlag := false
		if pool.ObjectStack.Len() != 0 {
			obj = pool.ObjectStack.Pop()
			flag = pool.ObjectFactory.Valid(obj)
			if flag { //判断获取对象是否有效 , flag=true 是有效的。
				pool.usedObjectCounter.add()

				pool.ObjectFactory.GetEnd(obj)
			}
			objFlag = true
		}
		if !flag { // 当获取的对象无效时，暂停0.5s后再尝试
			//fmt.Println("retry connection>>>>>>>>>")

			if !objFlag { // 说明对象池中没有对象了。
				log4go.DebugLog(message.WAR_CORE_69001)
				// 每次休眠，都比上一次休眠多一倍时间。
				time.Sleep(time.Duration(UTILS_POOL_TRY_TIMES_PER_SLEEPTIME*tryTimesCurr) * time.Millisecond)
			} else { // 如果对象池有对象，但是对象无效，那么就立即重试。
				log4go.DebugLog(message.WAR_CORE_69002)
			}
		}
		tryTimesCurr++
	}
	if flag {
		return obj
	} else {
		//log4go.ErrorLog("GetObject in pool is error  , current object is nil . (try get is %d times.)", tryTimesCurr)
		log4go.ErrorLog(message.ERR_CORE_39001, tryTimesCurr)
		return nil
	}
}

/**
 * 使用完，释放会对象池
 */
func (pool *AbstractPool) Release(obj interface{}) {
	pool.Lock.Lock()
	defer func() {
		pool.Lock.Unlock()
		pool.Cond.Signal() // 每次只通知一个等待 ， 那么只会触发一个运行。不需要在锁期间执行。
	}()

	pool.ObjectFactory.ReleaseStart(obj)
	if pool.usedObjectCounter.get() > 0 {
		pool.usedObjectCounter.sub()
	}
	pool.ObjectStack.Push(obj)

	pool.ObjectFactory.ReleaseEnd(obj)

	//pool.Cond.Broadcast()  ／／直接通知全部等待 ， 那么就全部触发运行

}
