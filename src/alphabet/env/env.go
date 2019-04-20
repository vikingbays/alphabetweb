// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*
定义webapp的全局环境变量。

*/
package env

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
)

// 设置启动时使用的CPU数量
var Env_Project_Cpunum = 1

// 项目根目录，例如： /home/glassProject
var Env_Project_Root string

// 项目资源目录，假设资源在 src下，那么就指向他
//
// 该目录信息，是完成的文件夹路径，例如： /home/glassProject/src
//
var Env_Project_Resource string

// 在src下的根目录名（非路径）,所有app都在此目录下。
var Env_Project_Resource_Apps []string = []string{"apps"}

// 国际化目录名称（非路径）
var Env_Project_Resource_Apps_I18n string = "i18n"

// 国际化文件后缀名
var Env_Project_Resource_Apps_I18n_File_Suffix string = "toml"

//var Env_Project_Resource_Apps_Model string = "model"

//var Env_Project_Resource_Apps_Action string = "action"

//var Env_Project_Resource_Apps_View string = "view"

// 路由文件名称，存放在 Env_Project_Resource_Apps 目录下。
var Env_Project_Resource_Apps_WebConfig string = "route.toml"

// sql语句配置路径
var Env_Project_Resource_Apps_SqlConfig string = "config/db"

// 日志文件名称，存放在 Env_Project_Resource_Apps 目录下。
var Env_Project_Resource_LogConfig string = "logconfig.toml"

// 日志文件名称，存放在 Env_Project_Resource_Apps 目录下。
var Env_Project_Resource_LogHttpConfig string = "logconfig_http.toml"

// 数据库连接文件名称，存放在 Env_Project_Resource_Apps 目录下。
var Env_Project_Resource_DbConfig string = "dbsconfig.toml"

// redis缓存连接文件名称，存放在 Env_Project_Resource_Apps 目录下。
var Env_Project_Resource_CacheConfig string = "cachesconfig.toml"

// web配置文件名称，存放在 Env_Project_Resource_Apps 目录下。
var Env_Project_Resource_WebConfig string = "webconfig.toml"

// service（manager/rpc）配置文件名称。设置微服务的server端配置文件，只能有一个，后缀名是：toml  ， 例如： ms_server_config.toml
var Env_Project_Resource_MS_ServerConfig_Name string = "ms_server_config.toml"

// 设置微服务的client端配置文件，可以有多个，后缀名是：toml ， 例如： ms_client_config_01.toml , ms_client_config_02.toml
var Env_Project_Resource_MS_ClientConfig_Prefix string = "ms_client_config"

var Env_Project_Resource_Sysconfig_Folder_Prefix string = "sysconfig"

var Env_Project_Resource_Sysconfig_Folder_Key string = ""

var Env_Project_Resource_Sysconfig_Folder string = ""

//指定web server访问的协议，
var Env_Server_Protocol_List []string

//设置web server路径，
var Env_Server_Addr_List []string

//设置web server访问端口
var Env_Server_Port_List []int

//设置web server访问端口
var Env_Server_Timeout_List []int

//设置web server访问端口
var Env_Server_Maxconn_List []int

//指定web访问的协议，在Env_Project_Resource_WebConfig 文件中配置。
//var Env_Web_Protocol string = "http"

//设置web路径，在Env_Project_Resource_WebConfig 文件中配置。
//var Env_Web_Addr string = ""

//设置web访问端口，在Env_Project_Resource_WebConfig 文件中配置。
//var Env_Web_Port int = 9669

//设置web上下文路径，在Env_Project_Resource_WebConfig 文件中配置。
var Env_Web_Context map[string]string

var Env_Web_Context_DEFAULT_VALUE string = "alphabet"

var Env_Web_Apps map[string][]string // key: appName，所在位置 src/${key},   values:  module名称，所在位置 src/${key}/${values}

//设置是否是生成模式，在Env_Project_Resource_WebConfig 文件中配置。
var Env_Web_Mode_IsProduct bool = true

//设置缺省的语言环境，在Env_Project_Resource_WebConfig 文件中配置。
var Env_Web_I18n string = "zh"

//设置sessonid信息，在Env_Project_Resource_WebConfig 文件中配置。
var Env_Web_Sessionid string = "alphabet-session-id"

//设置用于产生session的key
var Env_Web_Sessiongenkey string = "ABCD..WERT..3e@hjj$%kkk,,1~a"

var Env_Web_Sessionmaxage int = 3600

//存储最大的session对象
var Env_Web_Session_MaxObjectCount int = 10000

//设置session存储方式，支持内存方式(memory) 和 redis方式（redis）
var Env_Web_Session_Store string = "memory"

//设置session存储库的名称，如果是redis方式，那么就是设置cachesconfig.toml的某个库；如果是memory方式，那么就设置""。如果没有设置，那么默认是memory方式。
var Env_Web_Session_Store_Name string = ""

var Env_Web_HttpsCertFile string = ""

var Env_Web_HttpsKeyFile string = ""

// 判断是否运行Web模式。如果运行，会对部分程序进行特殊处理，例如：DataFinder中的sql别名都会加入appname名称作为区分。
var Env_Web_Running bool = false

// 文件上传的缓存参数，单位是字节码，例如 8388608 表示 8M，在Env_Project_Resource_WebConfig 文件中配置。
var Env_Web_Upload_Memcachesize map[string]int64 // = 8388608 //8M
var Env_Web_Upload_Memcachesize_DEFAULT_VALUE int64 = 8388608

// 文件上传的大小，单位是字节码，例如 268435456 表示 256M，在Env_Project_Resource_WebConfig 文件中配置。
var Env_Web_Upload_Maxsize map[string]int64 // = 268435456 //256M
var Env_Web_Upload_Maxsize_DEFAULT_VALUE int64 = 268435456

// 上传文件默认存放路径，在Env_Project_Resource_WebConfig 文件中配置。
var Env_Web_Upload_Storepath map[string]string
var Env_Web_Upload_Storepath_DEFAULT_VALUE string = "${project}/upload"

// 上传文件是否按照app名称分多个目录，在Env_Project_Resource_WebConfig 文件中配置。
var Env_Web_Upload_Splitonapps map[string]bool // = true
var Env_Web_Upload_Splitonapps_DEFAULT_VALUE bool = true

// 时间格式
var FORMAT_TIME_YMDHMS_NS string = "2006-01-02 15:04:05.000000000"

var Models map[string][]string

// 定义web header 的唯一序列号
var Env_Web_Header_Unique_Serial_Number string = "AB_UNIQUE_SERIAL_NUMBER"

// microService 最大客户端连接数。
var Env_MS_MaxConn int = 100

//pprof是否可用
var Env_Pprof bool = false

// pprof 可访问客户端ip
var Env_Pprof_ClientIP []string = make([]string, 0, 1)

// pprof 可访问用户名
var Env_Pprof_Username string = ""

// pprof 可访问密码
var Env_Pprof_Password string = ""

/*
env初始化工作，用于初始化全局变量环境。

@param path webapp项目路径，例如：/home/glassProject
*/
func Init(path string) {
	InitWithoutLoad(path)
	load()
}

func InitWithoutLoad(path string) {
	//path := os.Getenv("GOPATH")
	path = strings.Replace(path, "\\", "/", -1)
	if strings.HasSuffix(path, "/") {
		path = path[0:(len(path) - 1)]
	}
	Env_Project_Root = path
	Env_Project_Resource = path + "/src"

	Env_Web_Upload_Memcachesize = make(map[string]int64)
	Env_Web_Upload_Maxsize = make(map[string]int64)
	Env_Web_Upload_Storepath = make(map[string]string)
	Env_Web_Upload_Splitonapps = make(map[string]bool)

}

func load() {
	Models = make(map[string][]string)
	for _, rootAppName := range Env_Project_Resource_Apps {
		appDirs, err := ioutil.ReadDir(Env_Project_Resource + "/" + rootAppName)
		if err == nil {
			Models[rootAppName] = make([]string, 0)
			var idxApps int = 0
			for _, appdir := range appDirs {
				if appdir.IsDir() {
					Models[rootAppName] = append(Models[rootAppName], appdir.Name())
					idxApps++
				}
			}
		} else {
			fmt.Println(err.Error())
		}
	}
}

/*
通过运行时状态获取当前所在的app应用名

@param depth 表示调用深度数

@return 返回appname名称
*/
func GetAppInfoFromRuntime(depth int) string {
	_, file, _, _ := runtime.Caller(depth)
	file = filepath.ToSlash(file)

	for _, rootAppName := range Env_Project_Resource_Apps {
		appNameLastIndexStart := strings.LastIndex(file, "/src/"+rootAppName+"/")
		//fmt.Println("GetAppInfoFromRuntime:::::::  file = ", file, " __  ", Env_Project_Resource, " __ ", Env_Project_Resource_Apps, "____  ", appNameLastIndexStart)

		if appNameLastIndexStart != -1 {
			len0 := len("/src/" + rootAppName + "/")
			appNameLastIndexStart = appNameLastIndexStart + len0

			appNameLastIndexEnd := strings.Index(file[appNameLastIndexStart:], "/")
			if appNameLastIndexEnd != -1 {
				appNameLastIndexEnd = appNameLastIndexEnd + appNameLastIndexStart
				return file[appNameLastIndexStart:appNameLastIndexEnd]
			}
			break
		}
	}

	/* //打包有问题
	appNameIndexStart := strings.Index(file, Env_Project_Resource+"/"+Env_Project_Resource_Apps+"/")
	fmt.Println("GetAppInfoFromRuntime:::::::  file = ", file, " __  ", Env_Project_Resource, " __ ", Env_Project_Resource_Apps, "____  ", appNameIndexStart)
	if appNameIndexStart != -1 {
		len0 := len(Env_Project_Resource + "/" + Env_Project_Resource_Apps + "/")
		appNameIndexStart = appNameIndexStart + len0

		appNameIndexEnd := strings.Index(file[appNameIndexStart:], "/")
		if appNameIndexEnd != -1 {
			appNameIndexEnd = appNameIndexEnd + appNameIndexStart
			return file[appNameIndexStart:appNameIndexEnd]
		}

	}
	*/
	return ""
}

/*
以当前代码递归，查找被调用关系。

@param depth 表示调用深度数

@return 返回方法名和行号，结构是"[方法名]:[行号]"
*/
func GetCodeInfo(depth int) string {
	pc, _, lineno, ok := runtime.Caller(depth)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	return src
}

/*
获取当前所有调用关系

@param depth 表示调用深度数

@return 返回方法名和行号，结构是"[方法名]:[行号]"
*/
func GetInheritCodeInfoAlls(fileName string, after bool) string {
	pcs := make([]uintptr, 100)
	depths := runtime.Callers(1, pcs)
	src := "InheritCallers , this structure is [lineno] : [method]  ------ filepath:  [filename] "
	showFlag := true
	if fileName != "" {
		showFlag = false
	}
	for i := 0; i < depths; i++ {
		obj := runtime.FuncForPC(pcs[i])
		thisFilename, lineno := obj.FileLine(pcs[i])

		if !showFlag {
			if strings.HasSuffix(thisFilename, fileName) {
				showFlag = true
				if after {
					continue
				}
			} else {
				continue
			}
		}
		if showFlag {
		}
		src = fmt.Sprintf("%s \n [%5d] : [%s] \n           ------ filepath:  %s", src, lineno, obj.Name(), thisFilename)
	}
	return src
}

/**
 * 设置sysconfig配置目录，入参 key，表示的目录是 sysconfig.[key]。 例如：“base”，找到对应的配置目录，例如：“sysconfig.base” 。
 */
func SetSysconfigPath(key string) {
	Env_Project_Resource_Sysconfig_Folder_Key = key
	if Env_Project_Resource_Sysconfig_Folder_Key != "" {
		Env_Project_Resource_Sysconfig_Folder = Env_Project_Resource_Sysconfig_Folder_Prefix + "." + Env_Project_Resource_Sysconfig_Folder_Key + "/"
	}
}

func GetSysconfigPath() string {
	return Env_Project_Resource_Sysconfig_Folder
}
