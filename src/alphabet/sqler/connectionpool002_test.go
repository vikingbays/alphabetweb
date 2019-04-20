// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*
该文件只测试连接池并发请求性能，包含连接异常情况下的测试

执行：
go test -v -test.run Test_connectionpool001

判断如果 count %10 ==0 情况下，就会关闭数据库连接。

*/

package sqler

import (
	"alphabet/core/abtest"
	"database/sql"
	"sync"
	"testing"
	"time"
)

/*
CREATE TABLE public.alphabet_db_user1
(
  usrid integer,
  name character varying(100),
  nanjing boolean,
  money numeric(20,2)
)
*/

// 初始化测试连接池参数参数
func test_connectionpool_initparam002() (dbMaxPoolSize int, requestCount int, exitSleep int, sqlstring string, dbName string,
	dbDriverName string, dbDataSourceName string) {

	dbMaxPoolSize = 100 //数据库初始连接池
	requestCount = 500  // 同时并发请求多少次
	exitSleep = 10      // 退出测试程序前，是否需要等待，单位秒

	sqlstring = "select * from alphabet_db_user1 where usrid>0 and usrid<199" // 设置测试的sql

	dbName = "db_test1"
	dbDriverName = "postgres"
	dbDataSourceName = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&application_name=alphabet"

	return

}

//go test -v -test.run Test_connectionpool002
//测试数据库连接池的并发请求下的执行情况。
func Test_connectionpool002(t *testing.T) {

	abtest.InitTestEnv("")
	abtest.InitTestLog("")

	dbMaxPoolSize, requestCount, exitSleep, sqlstring, dbName, dbDriverName, dbDataSourceName := test_connectionpool_initparam002()

	abtest.InfoTestLog(t, "开始测试！")

	wg := sync.WaitGroup{}
	wg.Add(requestCount)

	var cp ConnectionPool
	cp.NewConnectionPool(dbMaxPoolSize, dbName, dbDriverName, dbDataSourceName)

	dbs := make([]*sql.DB, dbMaxPoolSize)

	for i := 0; i < dbMaxPoolSize; i++ {
		dbs[i] = cp.GetConnection()
	}

	for i := 0; i < dbMaxPoolSize; i++ {
		cp.ReleaseConnection(dbs[i])
	}

	abtest.InfoTestLog(t, "Test_connectionpool002..........................第一次")
	lock := new(sync.Mutex)
	counter := new(test_connectionpool_counter002)
	tq1 := time.Now()
	for i := 0; i < requestCount; i++ {
		num := i
		go func() {
			defer wg.Done()
			test_connectionpool_query002(cp, lock, counter, sqlstring, (num%10 == 0), t)
		}()
	}

	wg.Wait()

	tq2 := time.Now()

	abtest.InfoTestLog(t, "Test_connectionpool002..........................第一次...... 完成")

	abtest.InfoTestLog(t, "最大数据库连接数(dbMaxPoolSize =%d), 计划测试的并发请求数（requestCount: %d）", dbMaxPoolSize, requestCount)
	abtest.InfoTestLog(t, "实际获取数据量总数：%d, 实际总请求数：%d", counter.rownum, counter.count)
	abtest.InfoTestLog(t, "总用时：%d 毫秒 , 平均请求时间：%d 毫秒, 平均查询时间：%d 毫秒", (tq2.Sub(tq1)).Nanoseconds()/1000/1000, counter.requestDuration/counter.count, counter.queryDuration/counter.count)

	time.Sleep(time.Duration(exitSleep) * time.Second)

	abtest.InfoTestLog(t, "Test_connectionpool002..........................第二次")
	wg = sync.WaitGroup{}
	wg.Add(requestCount)
	lock = new(sync.Mutex)
	counter = new(test_connectionpool_counter002)
	tq1 = time.Now()
	for i := 0; i < requestCount; i++ {
		num := i
		go func() {
			defer wg.Done()

			test_connectionpool_query002(cp, lock, counter, sqlstring, (num%10 == 0), t)

		}()
	}

	wg.Wait()

	tq2 = time.Now()

	abtest.InfoTestLog(t, "Test_connectionpool002..........................第二次...... 完成")

	abtest.InfoTestLog(t, "最大数据库连接数(dbMaxPoolSize =%d), 计划测试的并发请求数（requestCount: %d）", dbMaxPoolSize, requestCount)
	abtest.InfoTestLog(t, "实际获取数据量总数：%d, 实际总请求数：%d", counter.rownum, counter.count)
	abtest.InfoTestLog(t, "总用时：%d 毫秒 , 平均请求时间：%d 毫秒, 平均查询时间：%d 毫秒", (tq2.Sub(tq1)).Nanoseconds()/1000/1000, counter.requestDuration/counter.count, counter.queryDuration/counter.count)

	time.Sleep(time.Duration(exitSleep) * time.Second)

	abtest.InfoTestLog(t, "Test_connectionpool002..........................第三次")
	wg = sync.WaitGroup{}
	wg.Add(requestCount)
	lock = new(sync.Mutex)
	counter = new(test_connectionpool_counter002)
	tq1 = time.Now()
	for i := 0; i < requestCount; i++ {
		num := i
		go func() {
			defer wg.Done()
			test_connectionpool_query002(cp, lock, counter, sqlstring, (num%10 == 0), t)
		}()
	}

	wg.Wait()

	tq2 = time.Now()

	abtest.InfoTestLog(t, "Test_connectionpool002..........................第三次...... 完成")

	abtest.InfoTestLog(t, "最大数据库连接数(dbMaxPoolSize =%d), 计划测试的并发请求数（requestCount: %d）", dbMaxPoolSize, requestCount)
	abtest.InfoTestLog(t, "实际获取数据量总数：%d, 实际总请求数：%d", counter.rownum, counter.count)
	abtest.InfoTestLog(t, "总用时：%d 毫秒 , 平均请求时间：%d 毫秒, 平均查询时间：%d 毫秒", (tq2.Sub(tq1)).Nanoseconds()/1000/1000, counter.requestDuration/counter.count, counter.queryDuration/counter.count)

	abtest.InfoTestLog(t, "测试结束！")
	//time.Sleep(time.Duration(exitSleep*2) * time.Second)
}

type test_connectionpool_counter002 struct {
	rownum          int64
	count           int64
	queryDuration   int64
	requestDuration int64
}

func test_connectionpool_query002(cp ConnectionPool, lock *sync.Mutex, counter *test_connectionpool_counter002, sqlstring string, isDoCLose bool, t *testing.T) {
	tq1 := time.Now()
	db := cp.GetConnection()
	if db == nil {
		return
	}
	defer cp.ReleaseConnection(db)
	if isDoCLose {
		defer db.Close()
	}
	tq2 := time.Now()

	tx, errTX := db.Begin()
	if errTX != nil {
		abtest.ErrorTestLog(t, errTX.Error())
	}
	defer tx.Commit()

	stmt, errSTMT := tx.Prepare(sqlstring)
	if errSTMT != nil {
		abtest.ErrorTestLog(t, errSTMT.Error())
	}
	defer stmt.Close()

	rows, errRow := stmt.Query()
	if errRow != nil {
		abtest.ErrorTestLog(t, errRow.Error())
	}
	defer rows.Close()

	num := 0
	for rows.Next() {
		num = num + 1
	}
	tq3 := time.Now()
	d32 := (tq3.Sub(tq2)).Nanoseconds() / 1000 / 1000
	d31 := (tq3.Sub(tq1)).Nanoseconds() / 1000 / 1000
	lock.Lock()
	counter.rownum = counter.rownum + int64(num)
	counter.count++
	counter.requestDuration = counter.requestDuration + d31
	counter.queryDuration = counter.queryDuration + d32
	lock.Unlock()

	//	t.Log(fmt.Sprintf("%v , 行数：%d .   .  获取连接后的查询时间(tq3-tq2): %d  . 该方法的整体时间（含连接等待时间）(tq3-tq1): %d", &db, num, d32, d31))

}
