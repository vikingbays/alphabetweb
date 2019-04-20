// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/core/toml"
	"alphabet/env"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

/*
 针对web模块进行初始化工作，主要做如下3个步骤：
  初始化：web配置，主要针对apps/webconfig.toml 配置文件
  初始化 Session
  初始化路由
*/
func Init() {
	InitBase()

	go func() {
		for true {
			time.Sleep(time.Second * 10)
			alphaSessionFactory.FlushAll()
		}
	}()
}

func InitBase() {
	for _, rootAppName := range env.Env_Project_Resource_Apps {
		initWebConfig(rootAppName)
	}
	InitSession(env.Env_Web_Sessionid)
	//InitSession("alpion-id", "/Users/vikingbays/golang/AlphabetwebProject/alphabetweb/sesspath", []byte("something-very-secret"))
	InitRouteStruct()

}

/*
apps/webconfig.toml 配置文件
例如：

  [[servers]]
  protocol="http"                   # 定义协议，采用http方式， 支持协议有：http，https, fcgi_unix ,rpc_unix,
                                    # 其中fcgi_unix是指：fcgi UnixDomainSocket ,rpc_unix是指：rpc方式的 UnixDomainSocket
  addr=""                           # 定义地址, 如果protocol="unix"，那么addr就需要设置为xxxx.sock的文件地址，例如： "/tmp/alphabetsample.sock"
                                    # 支持 ${project} 变量，表示项目根目录
                                    #   addr="${project}/alphabetsample.sock"
  port=9000                         # 定义端口
  timeout=10                        # 超时时长，单位秒
  maxconn=100                       # 最大连接数，暂时不支持

  [[web]]
  context="web2"
  apps="*"
  mode="develop"                    # develop 是开发者模式， product 是生产模式
  i18n="en"                         # 设置缺省的国际化信息，如果定义 i18n="en" ，那么就会查找 /apps/app1～n/i18n/en.toml内容
  sessionid="alphabet-session-id"   # 设置sessionid 信息
  sessionmaxage=3600                # 设置session失效日期，单位是秒
  httpscertfile="/Users/vikingbays/golang/AlphabetwebProject/https_cert/cert.pem"     # 设置https的证书cert
  httpskeyfile="/Users/vikingbays/golang/AlphabetwebProject/https_cert/key.pem"       # 设置https的证书key

  [[uploadfile]]
  memcachesize=8388608              # 文件在读取过程中的缓存大小,单位字节，在一定范围内越大，读写速度越快
  maxsize=33554432                  # 文件大小,单位字节
  storepath="${project}/upload"     # 文件上传的存储地址,其中 ${project} 表示项目根目录，例如：${project}/upload
  splitonapps=true                  # 是否按照应用（apps）分目录

  [[project]]
  cpunum=3                       # 设置CPU使用个数
  appsname="apps"                # 设置包名，应用于 ${project}/src/{appsname}，所有的app应用都在此目录下

*/
type Web_Config_List_Toml struct {
	Servers    []Server_Config_Toml
	Web        []Web_Config_Toml
	Uploadfile []Upload_Config_Toml
	Project    []Project_Config_Toml
	Pprof      []Pprof_Config_Toml
}

/*
Server服务配置
*/
type Server_Config_Toml struct {
	Protocol string // 协议： http 或  https ，存储到全局变量 env.Env_Server_Protocol_List 中
	Addr     string // 地址 ，存储到全局变量 env.Env_Server_Addr_List 中
	Port     int    // 端口 ，存储到全局变量 env.Env_Server_Port_List 中
	Timeout  int    // 超时，单位秒，存储到全局变量 env.Env_Server_Timeout_List 中
	Maxconn  int    // 最大连接数，存储到全局变量 env.Env_Server_Maxconn_List 中
}

/*
Web访问配置
*/
type Web_Config_Toml struct {
	//Protocol           string // 协议： http 或  https ，存储到全局变量 env.Env_Web_Protocol 中
	//Addr               string // 地址 ，存储到全局变量 env.Env_Web_Addr 中
	//Port               int    // 端口 ，存储到全局变量 env.Env_Web_Port 中
	Context            string // web访问的根路径配置，存储到全局变量 env.Env_Web_Context 中，例如 web2 , 那么所有访问的入口都是：http://localhost:9999/web2/ 开始。
	Apps               string // 可加载的appname名称，存储到全局变量 env.Env_Web_Apps 中， * 表示所有apps下的一级文件夹都作为应用 ，如果是多个appname名称，就采用“,”分隔。
	Mode               string // web运行模式，develop 是开发者模式， product 是生产模式
	I18n               string // 设置缺省的国际化信息，如果定义 i18n="en" ，那么就会查找 /apps/app1～n/i18n/en.toml内容
	Sessionid          string // 设置sessionid 信息
	Sessionmaxage      int    // 设置session失效日期，单位是秒
	Sessionobjectcount int    // 设置session可以存储的对象数，超过该对象，就会启动清理 ，一个对象就是一个session会话
	Sessiongenkey      string // 设置session产生key，用于加密
	Sessionstore       string // 设置session存储方式，支持 memory 和 redis 两种
	Sessionstorename   string // 设置session存储库的名称，如果是redis方式，那么就是设置cachesconfig.toml的某个库；如果是memory方式，那么就设置"" 。
	Httpscertfile      string // 设置https的证书cert
	Httpskeyfile       string // 设置https的证书key
}

/*
文件上传 相关参数配置
*/
type Upload_Config_Toml struct {
	Memcachesize int64  //文件在读取过程中的缓存大小,单位字节，在一定范围内越大，读写速度越快 ，存储到全局变量 env.Env_Web_Upload_Memcachesize 中
	Maxsize      int64  //文件大小,单位字节 ，存储到全局变量 env.Env_Web_Upload_Maxsize 中
	Storepath    string //文件上传的存储地址,其中 ${project} 表示项目根目录，例如：${project}/upload ，存储到全局变量 env.Env_Web_Upload_Storepath 中
	Splitonapps  bool   //是否按照应用（apps）分目录 ，存储到全局变量 env.Env_Web_Upload_Splitonapps 中
}

/*
项目 相关参数配置
*/
type Project_Config_Toml struct {
	Cpunum int //设置CPU使用个数
}

/*
pprof 相关参数配置
*/
type Pprof_Config_Toml struct {
	ClientIP []string //客户端可访问的ip地址，可以是数组
	Username string   //用户名
	Password string   //密码
}

/*
 初始化 webapp配置信息，针对apps/webconfig.toml 配置文件

 @param rootAppName  所在位置 src/${rootAppName}
*/
func initWebConfig(rootAppName string) {
	webconfigFile := env.Env_Project_Resource + "/" + rootAppName + "/" + env.GetSysconfigPath() + env.Env_Project_Resource_WebConfig
	var webConfigListToml Web_Config_List_Toml
	if _, err := toml.DecodeFile(webconfigFile, &webConfigListToml); err != nil {
		fmt.Println(err)
	} else {
		for _, project := range webConfigListToml.Project {
			env.Env_Project_Cpunum = project.Cpunum

		}

		if env.Env_Web_Apps == nil {
			env.Env_Web_Apps = make(map[string][]string)
		}

		env.Env_Server_Protocol_List = make([]string, len(webConfigListToml.Servers))
		env.Env_Server_Addr_List = make([]string, len(webConfigListToml.Servers))
		env.Env_Server_Port_List = make([]int, len(webConfigListToml.Servers))
		env.Env_Server_Timeout_List = make([]int, len(webConfigListToml.Servers))
		env.Env_Server_Maxconn_List = make([]int, len(webConfigListToml.Servers))

		for i, server0 := range webConfigListToml.Servers {
			env.Env_Server_Protocol_List[i] = server0.Protocol
			env.Env_Server_Addr_List[i] = strings.Replace(server0.Addr, "${project}", env.Env_Project_Root, -1)
			env.Env_Server_Port_List[i] = server0.Port
			env.Env_Server_Timeout_List[i] = server0.Timeout
			env.Env_Server_Maxconn_List[i] = server0.Maxconn
		}

		if env.Env_Web_Context == nil {
			env.Env_Web_Context = make(map[string]string)
		}

		for _, web := range webConfigListToml.Web {
			//env.Env_Web_Protocol = web.Protocol
			//env.Env_Web_Addr = strings.Replace(web.Addr, "${project}", env.Env_Project_Root, -1)
			//env.Env_Web_Port = web.Port
			env.Env_Web_Context[rootAppName] = web.Context
			env.Env_Web_I18n = web.I18n
			if web.Mode == "develop" {
				env.Env_Web_Mode_IsProduct = false
			} else {
				env.Env_Web_Mode_IsProduct = true
			}
			env.Env_Web_Sessionid = web.Sessionid
			env.Env_Web_Sessionmaxage = web.Sessionmaxage
			env.Env_Web_Session_MaxObjectCount = web.Sessionobjectcount
			env.Env_Web_Session_Store = web.Sessionstore
			env.Env_Web_Session_Store_Name = web.Sessionstorename
			if web.Sessionstorename == "" {
				env.Env_Web_Session_Store = "memory"
			}

			if web.Sessiongenkey != "" {
				env.Env_Web_Sessiongenkey = web.Sessiongenkey
			}

			//env.Env_Web_HttpsCertFile = web.Httpscertfile
			//env.Env_Web_HttpsKeyFile = web.Httpskeyfile

			env.Env_Web_HttpsCertFile = strings.Replace(web.Httpscertfile, "${project}", env.Env_Project_Root, -1)
			env.Env_Web_HttpsKeyFile = strings.Replace(web.Httpskeyfile, "${project}", env.Env_Project_Root, -1)

			env.Env_Web_Running = true

			if web.Apps == "" || web.Apps == "*" {
				appsFolderPath := env.Env_Project_Resource + "/" + rootAppName
				appsFolders, err := ioutil.ReadDir(appsFolderPath)
				if err == nil {
					for _, appFolder := range appsFolders {
						appFolderName := appFolder.Name()
						if appFolder.IsDir() {
							if env.Env_Web_Apps[rootAppName] == nil {
								env.Env_Web_Apps[rootAppName] = make([]string, 0, 1)
							}
							env.Env_Web_Apps[rootAppName] = append(env.Env_Web_Apps[rootAppName], appFolderName)
						}
					}
				} else {
					fmt.Println("This Web is not contain any apps . ", err.Error())
				}

			} else {
				apps := strings.Split(web.Apps, ",")
				for _, app := range apps {
					if env.Env_Web_Apps[rootAppName] == nil {
						env.Env_Web_Apps[rootAppName] = make([]string, 0, 1)
					}
					env.Env_Web_Apps[rootAppName] = append(env.Env_Web_Apps[rootAppName], app)
				}
			}

		}

		if webConfigListToml.Uploadfile != nil && len(webConfigListToml.Uploadfile) > 0 {
			if webConfigListToml.Uploadfile[0].Memcachesize == 0 {
				webConfigListToml.Uploadfile[0].Memcachesize = env.Env_Web_Upload_Memcachesize_DEFAULT_VALUE
			}
			if webConfigListToml.Uploadfile[0].Maxsize == 0 {
				webConfigListToml.Uploadfile[0].Maxsize = env.Env_Web_Upload_Maxsize_DEFAULT_VALUE
			}

			if webConfigListToml.Uploadfile[0].Storepath == "" {
				webConfigListToml.Uploadfile[0].Storepath = strings.Replace(env.Env_Web_Upload_Storepath_DEFAULT_VALUE, "${project}", env.Env_Project_Root, -1)
			} else {
				webConfigListToml.Uploadfile[0].Storepath = strings.Replace(webConfigListToml.Uploadfile[0].Storepath, "${project}", env.Env_Project_Root, -1)
			}
		} else {
			webConfigListToml.Uploadfile = make([]Upload_Config_Toml, 0, 1)
			webConfigListToml.Uploadfile = append(webConfigListToml.Uploadfile,
				Upload_Config_Toml{Memcachesize: env.Env_Web_Upload_Memcachesize_DEFAULT_VALUE,
					Maxsize:     env.Env_Web_Upload_Maxsize_DEFAULT_VALUE,
					Storepath:   strings.Replace(env.Env_Web_Upload_Storepath_DEFAULT_VALUE, "${project}", env.Env_Project_Root, -1),
					Splitonapps: true})
		}

		{

			env.Env_Web_Upload_Memcachesize[rootAppName] = webConfigListToml.Uploadfile[0].Memcachesize
			env.Env_Web_Upload_Maxsize[rootAppName] = webConfigListToml.Uploadfile[0].Maxsize
			env.Env_Web_Upload_Storepath[rootAppName] = webConfigListToml.Uploadfile[0].Storepath
			env.Env_Web_Upload_Splitonapps[rootAppName] = webConfigListToml.Uploadfile[0].Splitonapps

		}

		if len(webConfigListToml.Pprof) > 0 {
			ips := webConfigListToml.Pprof[0].ClientIP
			if ips != nil {
				flag := false
				for _, ip := range ips {
					if ip == "*" {
						flag = true
						break
					}
				}
				if !flag {
					env.Env_Pprof_ClientIP = ips
				}
			}
			env.Env_Pprof_Username = webConfigListToml.Pprof[0].Username
			env.Env_Pprof_Password = webConfigListToml.Pprof[0].Password
			if env.Env_Pprof_Password == "" || env.Env_Pprof_Username == "" {
				env.Env_Pprof_Username = ""
				env.Env_Pprof_Password = ""
			}
			env.Env_Pprof = true
		}

		//	fmt.Printf("alphabet/web/web.go[line 314]  rootAppName is  [%s] . Env_Web_Context is : [%s] \n", rootAppName, env.Env_Web_Context[rootAppName])

	}

}
