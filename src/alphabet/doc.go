// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

//
// alphabet是一个Web应用框架，集成了：web访问、数据库查询、日志处理等能力。alphabet框架是采用容器的设计思想，每个应用都是相对独立封装。
//
// 阅读内容包括：
//
// 1.代码结构介绍
//
// 2.运行机制
//
// 3.web应用开发介绍
//
// 4.Error F&Q
//
//-------------------------------------------------------------------------------
//
// 1.代码结构介绍
//
// 创建一个webapp的代码结构是：
//  glassProject
//    bin
//    logs
//    pkg
//    upload
//    src
//      glass                             // 应用包集合，包含了 app1,app2,app3 等应用。约定在src下只有一个文件夹。
//        app1                            // 一个应用包，实现该应用的所有功能，具有独立的国际化、控制器、展现等。
//          config
//            db
//              pg_test1.toml             // 配置pg_test1库的所有sql语句信息。[详细解释如下]
//              ...
//          i18n
//            en.toml                     // 针对 en 的国际化配置。
//            zh.toml                     // 针对 zh 的国际化配置。
//            ...
//          module
//          resource
//            css                         // css文件夹，包含：css文件，可以多个。
//              app1.css
//              ...
//            img                         // img文件夹
//            js                          // javascript文件夹
//          view                          // 展现层文件，类似于jsp或者php文件，采用html/template模板实现，必须以gml为后缀名
//            helloworld_index.gml
//            helloworld_gift.gml
//            ...
//          helloworld.go                 // 控制器处理类，必须在应用包app1下定义。实现func func1(c *web.Context){} 方法处理，方法名可自定义。
//          restful.go
//          ...
//          route.toml                    // 该应用包的路由配置，包括注册：控制器、展现层、过滤器(全局)等
//        app2
//        app3
//        ...
//        dbsconfig.toml                  // 数据库配置文件，包括可注册的数据库，例如：app1中注册的pg_test1 。[详细解释如下]
//        logconfig.toml                  // 日志配置文件。 [详细解释如下]
//        logconfig_http.toml             // http日志配置文件。 [详细解释如下]
//        webconfig.toml                  // web配置文件。 [详细解释如下]
//        cachesconfig.toml               // (可选)cache配置文件。 ［详细解释如下］
//        ms_server_config.toml           // (可选)微服务的服务端配置文件 . 具体参考 /alphabet/service
//        ms_client_config_01.toml        // (可选)微服务的服务端配置文件，可以是多个 _02 _03 _04 ....  具体参考 /alphabet/service
//
// 其中，dbsconfig.toml的定义结构可参考：
//  [[dbs]]                               ## 可定义多个数据库连接 。
//  name="pg_test1"                       ## 数据库连接别名 , 同时也确定了sql配置文件是： app[x]/config/db/pg_test1.toml 。
//  driverName="postgres"                 ## 数据库驱动类型，当前支持 postgresql,mysql,sqlite 分别定义为：postgres,mysql,sqlite3 。
//  dataSourceName="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&application_name=alphabet"
//                                        ## 数据库连接 。
//  maxPoolSize=5                         ## 最大连接数 。
//
//  [[dbs]]
//  ...
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
//
//
// 其中，logconfig.toml / logconfig_http.toml 的定义结构可参考：
//  ##日志输出格式:
//  ##      %T - Time (15:04:05 MST)
//  ##      %t - Time (15:04)
//  ##      %D - Date (2006/01/02)
//  ##      %d - Date (01/02/06)
//  ##      %L - Level (FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR)
//  ##      %S - Source
//  ##      %G - goroutine ID  协程号
//  ##      %U - Unique serial number 请求唯一序列号 ，在web场景下，需要判断header是否有属性 env.Env_Web_Header_Unique_Serial_Number
//  ##      %M - Message
//  ##      It ignores unknown format strings (and removes them)
//  ##      Recommended: "[%D %T] [%L] (%S) %M"
//
//  [[filters]]                          ## 定义一个日志输出格式，可以定义多个。
//  enabled="true"                       ## 设置该日志是否启用
//  tag="stdout"                         ## 设置输出方式，采用控制台输出
//  type="console"                       ## 采用控制台输出
//  level="DEBUG"                        ## 级别定义 (FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR) ，FINEST最低。
//  format="[%D %T] [%L] (%S) %M "       ## 输出格式
//
//  [[filters]]
//  enabled="true"
//  tag="file"                           ## 文件方式
//  type="file"                          ## 文件方式
//  level="INFO"
//  format="[%D %T] [%L] (%S) %M "
//  filename="${project}/logs/apps.log"  ## 可以带变量 ${project}/logs/test.log
//  rotate="false"                       ## 是否采用循环输出
//  rotateMaxSize="0"                    ## \d+[KMG]? Suffixes are in terms of 2**10
//  rotateMaxLines="0"                   ## \d+[KMG]? Suffixes are in terms of thousands
//  rotateDaily="false"
//
// 其中，webconfig.toml的定义结构可参考：
//  [[servers]]                          ## 支持多个server服务,可以一起启动
//  protocol="http"                      ## 定义协议，采用http方式， 支持协议有：http，https, fcgi_unix ,rpc_unix,rpc_tcp,rpc_tcp_ssl
//                                       ## 其中fcgi_unix是指：fcgi UnixDomainSocket ,rpc_unix是指：rpc方式的 UnixDomainSocket
//  addr=""                              ## 定义地址, 如果protocol="unix"，那么addr就需要设置为xxxx.sock的文件地址，例如： "/tmp/alphabetsample.sock"
//                                       ## 支持 ${project} 变量，表示项目根目录
//                                       ##  addr="${project}/alphabetsample.sock"
//  port=9000                            ## 定义端口
//  timeout=10                           ## 超时时长，单位秒
//  maxconn=100                          ## 最大连接数，暂时不支持
//
//  [[servers]]
//  protocol="https"
//  addr=""
//  port=9443
//  timeout=10                             # 超时时长，单位秒
//  maxconn=100                            # 最大连接数，暂时不支持
//
//  [[servers]]
//  protocol="fcgi_unix"
//  addr="${project}/alphabetsample.sock"  # 特别注意，该地址长度不要超过100字节
//  timeout=10                             # 超时时长，单位秒
//  maxconn=100                            # 最大连接数，暂时不支持
//
//  #[[servers]]
//  #protocol="rpc_unix"
//  #addr="${project}/alphabetsample_rpc.sock"   # 特别注意，该地址长度不要超过100字节
//  #timeout=10                                  # 超时时长，单位秒
//  #maxconn=100                                 # 最大连接数，暂时不支持
//
//  [[servers]]
//  protocol="rpc_tcp"
//  addr=""
//  port=7777
//  #timeout=5                             # 超时时长，单位秒
//  maxconn=100                            # 最大连接数，暂时不支持
//
//  [[servers]]
//  protocol="rpc_tcp_ssl"
//  addr=""
//  port=7766
//  #timeout=5                             # 超时时长，单位秒
//  maxconn=100                            # 最大连接数，暂时不支持
//
//  [[web]]
//  context="web2"                       ## 定义web访问根路径
//  apps="*"                             ## 是否过滤应用包
//  mode="develop"                       ## develop 是开发者模式， product 是生产模式，
//  i18n="en"                            ## 设置缺省的国际化信息，如果定义 i18n="en" ，那么就会查找 /apps/app**/i18n/en.toml内容
//  sessionid="alphabet09-session-id"    ## 设置sessionid 信息
//  sessionmaxage=3600                   ## 设置session失效日期，单位是秒
//  sessionobjectcount=10000             ## 设置session可以存储的对象数，超过该对象，就会启动清理 ，一个对象就是一个session会话
//  sessiongenkey="dA~$%3@2*sAw  (:sQQ"  ## 设置session产生key，用于加密
//  sessionstore="memory"                ## session存储方式，支持 memory 和 redis 两种
//  sessionstorename="cache2"            ## 设置session存储库的名称，
//                                       ## 如果是redis方式，那么就是设置cachesconfig.toml的某个库；
//                                       ## 如果是memory方式，那么就设置""。
//                                       ## 如果没有设置，那么默认是memory方式。
//  httpscertfile="/home/sample1/sampleproject/https_cert/cert.pem"     # 设置https的证书cert，支持 ${project} 变量，表示项目根目录
//  httpskeyfile="/home/sample1/sampleproject/https_cert/key.pem"       # 设置https的证书key，支持 ${project} 变量，表示项目根目录
//
//  [[uploadfile]]
//  memcachesize=8388608                 ## 文件在读取过程中的缓存大小,单位字节，在一定范围内越大，读写速度越快
//  maxsize=33554432                     ## 文件大小,单位字节
//  storepath="${project}/upload"        ## 文件上传的存储地址,其中 ${project} 表示项目根目录，例如：${project}/upload
//  splitonapps=true                     ## 是否按照应用（apps）分目录
//
//  [[project]]
//  cpunum=3                             ## 设置CPU使用个数
//  appsname="apps"                      ## 设置应用包集合名称，应用于 ${project}/src/{appsname}，所有的app应用都在此目录下
//
//  [[pprof]]                            ## 如果设置，就支持pprof，监控程序运行健康度
//                                       ## http://ip:port/debug/pprof/?username=user1&password=123456
//  clientIP=["*"]                       ## 客户端可访问的ip地址，可以是数组,如果配置 * 表示所有ip
//  username="user1"                     ## 用户名 ，如果不设置，表示不需要用户名/密码, 传递参数 username=user1&password=123456
//  password="123456"                    ## 密码 ，如果不设置，表示不需要用户名/密码。
//
// 其中，cachesconfig.toml的定义结构可参考：
//
//  [[caches]]                                                               ## 配置一个缓存库
//  name="cache1"                                                            ## 别名
//  dataSourceName="redis://user:abcdef123456@127.0.0.1:6379/1"              ## 采用tcp方式的连接
//  maxPoolSize=10                                                           ## 连接池
//
//  [[caches]]
//  name="cache2"
//  dataSourceName="redis://user:abcdef123456@/tmp/redisserv/redis.sock/2"   ## 采用unix方式的连接
//  maxPoolSize=10
//
//
// 如果采用https模式，需要如下设置：
//  protocol="https"
//
//  httpscertfile="/home/sample1/sampleproject/https_cert/cert.pem"     # 设置https的证书cert
//
//  httpskeyfile="/home/sample1/sampleproject/https_cert/key.pem"       # 设置https的证书key
//
// 如果测试使用，cert.pem和key.pem文件通过openssl生成。
//  openssl genrsa -out key.pem 2048
//
//  openssl req -new -x509 -key key.pem -out cert.pem -days 3650
//
// 如果采用fcgi_unix模式，需要如下设置：
//  protocol="fcgi_unix"
//
//  addr="/home/sample1/sampleproject/alphabetsample.sock"       ## 注意，该文件不会被删除，每次重启需要手工删除。
//
// 如果采用fcgi_unix模式，那么前置服务应该是fastcgi模式，假设采用nginx配置，具体配置如下：
//
//  worker_processes  1;
//
//  events {
//      worker_connections  1024;
//  }
//
//  http {
//      include       mime.types;
//      default_type  application/octet-stream;
//
//      sendfile        on;
//      keepalive_timeout  65;
//      client_max_body_size 200m;
//
//      server {
//          listen       8880;
//          server_name  localhost;
//
//          location ~ \.*$ {
//                  include         fastcgi.conf;
//                  fastcgi_pass    unix:/tmp/alphabetsample.sock;
//          }
//      }
//  }
//
// 如果rpc模式，请参考 alphabet/service包中的说明。
//
// 路由配置规则，请参考 alphabet/web包中的说明。
//
//
//-------------------------------------------------------------------------------
//
// 2.运行机制
//
// alphabet提供给glassProject项目启动和停止的命令，能够方便的进行项目启停。
// 使用命令启动的时候，实际上是根据所有app定义的route.toml配置信息生成主程序，存放到临时目录中 glassProject/main.go上，
// 执行过程中会生成pid文件，记录当前的进程号信息，web应用停止命令就是根据该进程号停止。
// 在项目运行过程中，所用的.gml的模板都可以实现动态修改。
//
// 注意：停止命令当前只支持linux系统。
//
// 启停命令管理，是在 alphabet/cmd/abserver.go 中实现。
//
// 启动命令说明
//  go run cmd/abserver.go -start /home/glassProejct/     ## `/home/glassProejct/` 就是当前webapp目录，具体结构参考上一章节
//
// 停止命令说明
//  go run cmd/abserver.go -stop /home/glassProejct/     ## `/home/glassProejct/` 就是当前webapp目录，具体结构参考上一章节
//
// 生成主程序命令说明
//  go run cmd/abserver.go -genmain /home/glassProejct/   ## `/home/glassProejct/` 就是当前webapp目录，具体结构参考上一章节
//
// 生成二进制运行程序命令说明
//  go run cmd/abserver.go -build /home/glassProejct/ /home/glassout/ windows   ## `/home/glassProejct/` 就是当前webapp目录，具体结构参考上一章节
//                                                                              ## `/home/glassout/` 输出的可执行工程路径
//                                                                              ## `windows` 生成windows平台执行程序
//
// 帮助说明
//  go run cmd/abserver.go -help
//
//    参考帮助：
//
//    命令：  abserver   [command1]  [command2]  [command3]  [command4]
//
//    command说明：(所有文件路径的分隔符都必须使用｀／｀)
//
//      (1)、[command1] 参数一：（只能一个有效）
//
//           -start    string
//                       启动web服务。
//                       参数1：web工程根路径，源码存放在 src/apps下。
//
//           -stop     string
//                       停止web服务。
//                       参数1：web工程根路径，源码存放在 src/apps下。
//
//           -genmain  string
//                       生成web服务启动代码，代码存储到 “${project}/src/apps”下。
//                       参数1：web工程根路径，源码存放在 src/apps下。
//
//           -build    string  string  string
//                       生成web服务启动二进制程序。
//                       参数1：web工程根路径，源码存放在 src/apps下。
//                       参数2：输出项目文件夹路径，默认可执行文件在此目录的bin目录下。
//                       参数3：平台定义，如果不传参数，表示就是当前平台。
//                               其中：｀windows｀  ：表示windows平台。
//                                     ｀mac｀      ：表示mac平台。
//                                     ｀linux｀    ：表示linux平台。
//                                     ｀linuxarm｀ ：表示linux的arm平台。
//                                     ｀all｀      ：表示生成所有平台。
//
//           -build_for_debug    string  string  string
//                       生成web服务启动二进制程序，支持gdb进行debug（ 在编译时增加参数 -gcflags "-N -l" ）。
//                       参数1：web工程根路径，源码存放在 src/apps下。
//                       参数2：输出项目文件夹路径，默认可执行文件在此目录的bin目录下。
//                       参数3: 平台定义，
//                               其中：｀windows｀  ：表示windows平台。
//                                     ｀mac｀      ：表示mac平台。
//                                     ｀linux｀    ：表示linux平台。
//                                     ｀linuxarm｀ ：表示linux的arm平台。
//
//           -help
//                       展现帮助信息。
//
//      (2)、[command2] 参数二：（只能一个有效）
//           -config_key    string
//                       设置关键字，例如：“base”，找到对应的配置目录，例如：“sysconfig.base” 。
//                       该参数值赋给env.Env_Project_Resource_Sysconfig_Folder_Key 。
//                       程序启动时，会加载此目录，例如：src/[apps名称]/sysconfig.base 。
//                       参数1: 设置这个key值。
//           -config_url    string
//                       设置一个url，获取远程配置项信息，实现服务配置集中管理。
//                       这个url地址会匹配如下目录：
//                               [url]/logconfig.toml
//                               [url]/logconfig_http.toml
//                               [url]/dbsconfig.toml
//                               [url]/cachesconfig.toml
//                               [url]/webconfig.toml
//                               [url]/ms_server_config.toml
//                               [url]/ms_client_config_01.toml
//                               [url]/...
//                       参数1: url地址。
//
//      (3)、[command3] 参数三：（可以多个有效）
//           -env_xxx    string
//                       设置环境变量，可以定义多个（使用方法os.Setenv）。其中 xxx 就是当前定义的key值。
//                       获取环境变量数据，采用 os.Getenv
//                       例如： -env_key1  value1    -env_key2  value2 ，那么就是：key1:value1,key2:value2
//
//      (4)、[command4] 参数四：（只能一个有效）
//           -appsname    string
//                       设置有效的apps名称，如果不设置表示全部有效，如果设置多个就多个有效，多个apps以逗号分割。
//                       例如： -appsname    oct_web,oct_service
//
//-------------------------------------------------------------------------------
//
// 3.web性能检测 pprof
//
// 当前alphabet 已经集成了pprof的性能数据收集模块
//
// 激活pprof模块，需要在 webconfig.toml 中配置：
//   [[pprof]]
//   clientIP=["*"]      # 客户端可访问的ip地址，可以是数组,
//                       # 例如：http://localhost:9000/debug/pprof/?username=user1&password=123456
//   username="user1"    # 用户名 , 传递参数 username=user1&password=123456
//   password="123456"   # 密码
//
// 通过web访问pprof信息
//   http://localhost:9000/debug/pprof/?username=user1&password=123456
//
// 基于pprof工具自助分析：   go tool pprof  xxxxxx 。注意：传入的参数 需要带上转义符 。
//
// 例如：查看堆栈信息 ，（使用top，调用链信息）
//   $ go tool pprof http://localhost:9000/debug/pprof/heap?username=user1\&password=123456
//
//     Fetching profile from http://localhost:9000/debug/pprof/heap?password=123456&username=user1
//     Saved profile in /Users/downloads/pprof/pprof.localhost:9000.alloc_objects.alloc_space.inuse_objects.inuse_space.003.pb.gz
//     Entering interactive mode (type "help" for commands)
//     (pprof) top
//     2585.33kB of 2585.33kB total (  100%)
//     Dropped 54 nodes (cum <= 12.93kB)
//     Showing top 10 nodes out of 57 (cum >= 512.31kB)
//           flat  flat%   sum%        cum   cum%
//      1048.94kB 40.57% 40.57%  1048.94kB 40.57%  runtime.hashGrow
//       512.31kB 19.82% 60.39%   512.31kB 19.82%  alphabet/core/pq.DialOpen
//       512.05kB 19.81% 80.20%   512.05kB 19.81%  html/template.htmlReplacer
//       512.02kB 19.80%   100%   512.02kB 19.80%  html/template.ensurePipelineContains
//              0     0%   100%   512.31kB 19.82%  alphabet/core/pq.(*drv).Open
//              0     0%   100%   512.31kB 19.82%  alphabet/core/pq.Open
//              0     0%   100%   512.31kB 19.82%  alphabet/core/utils.(*AbstractPool).Get
//              0     0%   100%  2048.95kB 79.25%  alphabet/mux.(*Router).ServeHTTP
//              0     0%   100%   512.31kB 19.82%  alphabet/sqler.(*ConnectionPool).GetConnection
//              0     0%   100%   512.31kB 19.82%  alphabet/sqler.(*DBObjectFactory).Valid
//      (pprof) png
//      Generating report in profile001.png
//
// 例如：查看cpu信息
//    $ go tool pprof http://localhost:9000/debug/pprof/profile?username=user1\&password=123456
//      Fetching profile from http://localhost:9000/debug/pprof/profile?password=123456&username=user1
//      Please wait... (30s)
//      Saved profile in /Users/vikingbays/pprof/pprof.localhost:9000.samples.cpu.002.pb.gz
//      Entering interactive mode (type "help" for commands)
//      (pprof) top
//      830ms of 990ms total (83.84%)
//      Showing top 10 nodes out of 217 (cum >= 10ms)
//            flat  flat%   sum%        cum   cum%
//           290ms 29.29% 29.29%      290ms 29.29%  nanotime
//           160ms 16.16% 45.45%      160ms 16.16%  syscall.Syscall
//           130ms 13.13% 58.59%      130ms 13.13%  runtime.mach_semaphore_signal
//            90ms  9.09% 67.68%       90ms  9.09%  runtime.mach_semaphore_wait
//            60ms  6.06% 73.74%       60ms  6.06%  runtime.usleep
//            30ms  3.03% 76.77%      130ms 13.13%  net/http.(*response).write
//            20ms  2.02% 78.79%       20ms  2.02%  reflect.(*structType).FieldByName
//            20ms  2.02% 80.81%       20ms  2.02%  runtime.greyobject
//            20ms  2.02% 82.83%       20ms  2.02%  runtime.kevent
//            10ms  1.01% 83.84%       10ms  1.01%  alphabet/log4go.(*CallChain).AddCaller
//      (pprof) png
//      Generating report in profile002.png
//
//-------------------------------------------------------------------------------
//
// 4.web应用开发介绍
//
// 具体请参考alphabet/web包说明。
//
//-------------------------------------------------------------------------------
//
// 5.微服务介绍
//
// 具体请参考alphabet/service包说明。
//
//-------------------------------------------------------------------------------
//
// 6.Error F&Q
//
// http: multiple response.WriteHeader calls
//  可能重复调用该方法：
//   w.WriteHeader(xxx)      //w 是 http.ResponseWriter 对象
//
//
//
//
package alphabet
