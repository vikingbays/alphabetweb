// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

//
// 实现数据库的访问、处理、查询等操作。当前支持 postgresql,mysql,sqlite 分别定义为：postgres,mysql,sqlite3 。
//
// sqler的实现思路是借鉴java的MyBatis的sqlMap的ORM的框架模型思路。
// 首先在dbconfig.toml文件中定义数据库连接，根据数据库名，在每个app下生成对应的sql配置文件。
// 假设在dbconfig.toml中定义｀name="pg_test1" ｀ ，那么就在app下创建｀pg_test1.toml｀ 文件，用于配置sql。
// 根据配置的sql别名，通过调用 dataFinder.xxxx 方法进行数据操作和查询等。
//
// 其中，dbsconfig.toml的定义结构可参考：
//  [[dbs]]                               ## 可定义多个数据库连接 。
//  name="pg_test1"                       ## 数据库连接别名 。
//  driverName="postgres"                 ## 数据库驱动类型，当前支持 postgresql,mysql,sqlite 分别定义为：postgres,mysql,sqlite3 。
//  dataSourceName="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&application_name=alphabet"
//                                        ## 数据库连接 。
//  maxPoolSize=5                         ## 最大连接数 。
//
//  [[dbs]]
//  ...
//
// 说明：如果是Mysql的配置如下：
//  [[dbs]]
//  name="pg_test1"
//  driverName="mysql"
//  dataSourceName="root:root1234@tcp(127.0.0.1:33060)/mysql?timeout=90s"
//  maxPoolSize=5
//
//
// 其中，pg_test1.toml用于定义的定义名称是"pg_test1"的数据库访问sql配置，该配置文件的结构可参考：
//  ## sql配置借鉴于Java的MyBatis的设计思想，采用SqlMap的配置方式，变量定义中，#{} 表示 ? 方式，${} 表示 直接当字符串替换。
//  ## #{} 的变量定义
//  ##    执行SQL：Select * from emp where name = #{employeeName}
//  ##    参数：   employeeName=>Smith
//  ##    解析后执行的SQL：Select * from emp where name = ？
//  ## ${} 的变量定义
//  ##    执行SQL：Select * from emp where name = ${employeeName}
//  ##    参数：   employeeName传入值为：Smith
//  ##    解析后执行的SQL：Select * from emp where name =Smith
//  ##
//
//  [[db]]                        ## 定义一个sql
//  name="deleteUserList"         ## 定义别名，通过 ｀ dataFinder.Exec("pg_test1", "deleteUserList", nil) ｀ 方式执行。
//  sql="""                       ## 具体的sql
//    delete From user1
//      """
//
//  [[db]]
//  name="getUserCount"
//  sql="""
//    select count(1) From user1
//      """
//
//  [[db]]
//  name="getUserList"            ## 定义别名，通过 ｀ dataFinder.QueryList("pg_test1", "getUserList", *paramUser1, *resultUser1) ｀ 方式执行。
//  sql="""
//    select * From user1 where usrid>#{minuSrid} and usrid<${maxusrID} and name like '${name}' and nanjing = #{nanjing}
//      """
//
// 如果需要解决多个数据库兼容性问题，例如：Mysql和Postgresql数据库下的系统表查询sql不一致问题。
// 引入 sql_postgres , sql_mysql , sql_sqlite3 属性支持。
//
//  ## 判断如果本次访问的数据库是mysql那么就使用｀sql_mysql｀定义的语句，其他数据库使用默认的｀sql｀定义的语句。
//
//  [[db]]
//  name="existUser"
//  sql="""
//      select count(1) from pg_class where relname = 'alphabet_db_user1'
//      """
//  sql_mysql="""
//      select count(1) from information_schema.TABLES WHERE table_name ='alphabet_db_user1'
//      """
//
//
//
// 一个数据库操作的代码结构是：
//  dataFinder, err := sqler.GetConnectionForDataFinder("pg_test1")  // 获取连接并开启数据库事务
//  defer dataFinder.ReleaseConnection(1)                             // 在defer中回收数据库连接，重新放回连接池中
//  if err != nil {
//	  log4go.ErrorLog(err)
//  } else {
//    dataFinder.StartTrans()                             // 开启一个事务
//	  ...                                                // 进行数据库操作或查询，调用方法有：DoQueryList，DoQueryMap，DoExec
//    ...
//    dataFinder.CommitTrans()                            // 提交一个事务
//  }
//
// 进行数据库操作或查询的方法有：
//  DoQueryList(...)   查询一组对象存储到List中，该对象可以是struct 也可以是基础类型。
//
//  DoQueryMap(...)    查询一组对象存储到Map中，存储信息是 key-value 模式的键值对 。其中：key 为某字段，value可以是某字段，也可以是该行纪录。
//
//  DoExec(...)        执行所有DML的操作，包括：Create、Drop、Insert、Delete、Update等操作
//
//  GetTx()          获取 sql.Tx 对象，用于执行原生sql 。
//
// 注意：以上数据库查询和操作的方法，如果在执行过程中出现异常，会自动使用panic抛出异常，使数据库事务回滚。
//
//
// -------------------------------------------------------------------------------
//
// [fix.001]
//   Q1: mysql查询时提交(Commit)或回滚(Rollback)操作的时候，会抛出异常信息，例如："packets.go:357: Busy buffer"。
//       产生现象的原因，是进行Commit和Rollback操作时，mysql数据库协议无返回值。因此抛出。
//
//   S1: 针对transaction.go此部分代码做优化，增加如下代码：
//       tx.mc.buf.length = 0        // 修复Busy buffer 的bug
//
//
//
package sqler
