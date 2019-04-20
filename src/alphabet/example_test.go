// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package alphabet_test

/*
import (
	"alphabet"
)
*/

func Example() {
	//创建代码目录结构如下：假设该目录结构是：  /home/project1/alphabetsample
	/*
	   alphabetsample
	     bin
	     logs
	     pkg
	     src
	       alphabetsample
	         app1
	           config
	             db
	               pg_test1.toml
	           i18n
	             en.toml
	           resource
	             css
	               app1.css
	             js
	               jquery-1.12.1.js
	           view
	             action1_run1_page1.gml
	           action1.go
	           route.toml
	         app2
	           ...
	       dbsconfig.toml
	       logconfig.toml
	       webconfig.toml
	*/

	//设置web配置信息：webconfig.toml
	/*
	   [[web]]
	   protocol="http"                     # 支持协议有：http，https , unix
	   addr=""                            # 如果是protocol="unix"， 该地址可定义为： "/tmp/alphabetsample.sock"
	   port=9000
	   context="web2"
	   apps="*"
	   mode="develop"                      # develop 是开发者模式， product 是生产模式 。 如果设置 develop模式可以查看pprof信息：http://ip:port/debug/pprof
	   i18n="en"                           # 设置缺省的国际化信息，如果定义 i18n="en" ，那么就会查找 /apps/appxxxx/i18n/en.toml内容
	   sessionid="alphabet09-session-id"   # 设置sessionid 信息
	   sessionmaxage=600                    # 设置session失效日期，单位是秒
	   sessionobjectcount=10000             # 设置session可以存储的对象数，超过该对象，就会启动清理 ，一个对象就是一个session会话
	   httpscertfile="/Users/vikingbays/golang/AlphabetwebProject/https_cert/cert.pem"     # 设置https的证书cert
	   httpskeyfile="/Users/vikingbays/golang/AlphabetwebProject/https_cert/key.pem"       # 设置https的证书key

	   [[uploadfile]]
	   memcachesize=8388608                # 文件在读取过程中的缓存大小,单位字节，在一定范围内越大，读写速度越快
	   maxsize=33554432                    # 文件大小,单位字节
	   storepath="${project}/upload"       # 文件上传的存储地址,其中 ${project} 表示项目根目录，例如：${project}/upload
	   splitonapps=true                    # 是否按照应用（apps）分目录

	   [[project]]
	   cpunum=3                            # 设置CPU使用个数
	   appsname="alphabetsample"           # 设置包名，应用于 ${project}/src/{appsname}，所有的app应用都在此目录下

	*/

	//设置数据库配置信息：dbsconfig.toml
	/*
	   [[dbs]]
	   name="pg_test1"
	   driverName="postgres"
	   dataSourceName="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&application_name=alphabet"
	   maxPoolSize=90
	*/

	//设置日志配置信息：logconfig.toml
	/*
	   [[filters]]
	   enabled="true"
	   tag="file"
	   type="file"
	   level="INFO"   # 级别定义 (:?FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR)
	   format="[%D %T] [%L] (%S) %M "
	*/

	//定义一个控制器：app1/action1.go
	/*
	   package app1

	   import (
	   	"alphabet/log4go"
	   	"alphabet/sqler"
	   	"alphabet/web"
	   )

	   // 需要创建user1表
	   // CREATE TABLE user1
	   // ( usrid integer,
	   //   name varchar(100),
	   //   nanjing boolean,
	   //   money numeric(20,2),
	   //   hello varchar(100)
	   // )

	   type User1 struct {
	   	Usrid    int
	   	Name     string
	   	Nanjing  bool
	   	Money    float32
	   	Hello    string
	   	MinUsrid int
	   	MaxUsrid int
	   }

	   func Action1Run1(context *web.Context) {

	   	dataFinder, err := sqler.GetDataFinderAndBeginTrans("pg_test1")
	   	defer dataFinder.EndTransAndClose(1)
	   	var list []interface{}
	   	if err != nil {
	   		log4go.ErrorLog(err)
	   		panic(&sqler.RollbackError{err})
	   	} else {
	   		paramUser1 := new(User1)
	   		paramUser1.MinUsrid = 0
	   		paramUser1.MaxUsrid = 1000
	   		paramUser1.Name = "Ikkk_%"
	   		paramUser1.Nanjing = false
	   		resultUser1 := new(User1)
	   		list, err = dataFinder.QueryList("pg_test1", "getUserList", *paramUser1, *resultUser1)
	   	}

	   	context.Return.Forward("action1_run1_page1", list)
	   }
	*/

	//定义一个视图：app1/view/action1_run1_page1.gml
	/*
		<html>
		<script src="resource/js/jquery-1.12.1.js"></script>
		<link rel="stylesheet" type="text/css" href="resource/css/app1.css">
		<body style='width:100%'>
		<div class='TitleLevel1Class'>
		       HelloWorld....你好！
		  </div>
		  <div class='ContentInfoClass'>
		    <div>
		        {{Locale "许愿" "helloworld" .I18n  }}  , {{Locale "礼物" "helloworld" .I18n  }}
		    </div>
			<div>
				session datas print :
			</div>
			<div>
				{{.Session}}
			</div>
			<div>
				params datas print :
			</div>
			<div>
				{{.Params}}
			</div>
			<div>
				usr datas print :
			</div>
			<div>
				{{.Datas}}
			</div>
			<div>
				{{range $num, $v := .Datas}}
					{{$num}} ,  {{$v}}
				{{end}}
			</div>
			<div>
				all datas print :
			</div>
			<div>
				{{.}}
			</div>
		</body>
		</html>
	*/

	//定义一个数据库配置信息：app1/config/db/pg_test1.toml
	/*
		[[db]]
		name="getUserList"
		sql="""
		  select * From user1 where usrid>#{minuSrid} and usrid<${maxusrID} and name like '${name}' and nanjing = #{nanjing}
		    """
	*/

	//定义国际化文件：app1/i18n/en.toml
	/*
		[[msg]]
		msgctxt="helloworld"   # 消息上下文
		msgid="许愿"            # 消息编码
		msgstr="Hope..."       # 消息内容

		[[msg]]
		msgctxt="helloworld"   # 消息上下文
		msgid="礼物"            # 消息编码
		msgstr="gift"          # 消息内容
	*/

	//定义路由器：app1/route.toml
	/*
		# action1 的处理

		[[action]]
		Action1Run1="action1run1"

		[[gml]]
		action1_run1_page1="view/action1_run1_page1.gml"
	*/

	///////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////

	//运行该web应用，使用abserver命令
	/*
		go run abserver.go -start /home/project1/alphabetsample/
	*/

	//打包成二进制项目
	/*
		go run abserver.go -build /home/project1/alphabetsample/  /home/project1/out/ linux
	*/

}
