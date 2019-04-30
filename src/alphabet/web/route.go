// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/core/utils"
	"alphabet/core/webutils"
	"alphabet/env"
	"alphabet/log4go"
	"alphabet/log4go/log4gohttp"
	"alphabet/log4go/message"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

var routeActionerList []*RouteActioner
var routeTemplateList map[string][]*RouteTemplate
var routeTemplateMap map[string]*RouteTemplate
var routeStaticsList map[string][]string

var routePage404, routePage500, routePage401, routeIndex, routeRootIndex map[string][]byte // routeRootIndex:根路径 , routeIndex: appname路径

var routePage404File, routePage500File, routeIndexWebPath, routeRootIndexWebPath map[string]string // routeRootIndexWebPath:根路径 , routeIndexWebPath: appname路径

var routeFitlerFunc map[string]func(context *Context) bool // key: apps, 根应用名，对应位置为： src/${key}

// 存储模板信息，按照应用域隔离，假设配置的名称是 abc，应用名是 app1, 根应用名是 sample1，那么实际存储的key是 sample1.app1.abc , 如果应用名不存在，表示是全局模板变量，那么key就是 .abc
var templateDataMap map[string]*template.Template

var muxerObject *Muxer

var DEFAULT_ROUTE_PAGE_404 []byte = []byte("Page Not Found ! 404 ")
var DEFAULT_ROUTE_PAGE_500 []byte = []byte("Page Error ! 500 ")
var DEFAULT_ROUTE_PAGE_401 []byte = []byte("Unauthorized ! 401 ")

/*
初始化路由信息
*/
func InitRouteStruct() {
	routeActionerList = make([]*RouteActioner, 0)
	routeTemplateList = make(map[string][]*RouteTemplate)
	routeTemplateMap = make(map[string]*RouteTemplate)
	routeStaticsList = make(map[string][]string)
	templateDataMap = make(map[string]*template.Template)

	routePage404 = make(map[string][]byte)
	routePage500 = make(map[string][]byte)
	routePage401 = make(map[string][]byte)

	routePage404File = make(map[string]string)
	routePage500File = make(map[string]string)
	routeIndexWebPath = make(map[string]string)
	routeRootIndexWebPath = make(map[string]string)

	routeIndex = make(map[string][]byte)

	routeRootIndex = make(map[string][]byte)

}

/*
封装控制器处理方法，需要考虑

@param  clazz  定义的控制器方法

@param  clazzType  定义方法的类型， 1: func( *Context)  ，2 ： func( *struct , *Context)

@param  url  请求路径

@param  appname  app的名称，如果是"" ，表示全局

@param  isRedirect  判断是不是需要重定向？

@return  func(http.ResponseWriter, *http.Request)  生成符合httpHandler规范的执行方法
*/
func doAction(clazz interface{}, clazzType int, url string, appname string, rootAppname string, isRedirect bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if env.Switch_CallChain { // 进入请求时，设置 唯一序列号
			usn := r.Header.Get(env.Env_Web_Header_Unique_Serial_Number)
			if usn == "" {
				usn = log4go.GetSerialnumber()
			}
			log4go.GetCallChain().AddCaller(log4go.GetGID(), usn)
		}
		startTime := time.Now()

		defer func() { // 用于捕获panic异常，不影响整个服务运行
			endTime := time.Now()
			if err := recover(); err != nil {
				if ActionMonitorExtendsObject != nil { // 如果实现了监听，就记录相关信息

					ActionMonitorExtendsObject.AddActionMonitor(url, appname, startTime.UnixNano(), endTime.UnixNano(), err.(error))
				}
				log4gohttp.ErrorLog("interval : %6d ms . protocol: %s . client.Addr : %s . server.Host : %s  , server.context: %s ,  server.url : %s  .  Webapp Err [errmsg: %v] .  \n %s", time.Now().Sub(startTime)/1000/1000,
					r.Proto, r.RemoteAddr, r.Host, env.Env_Web_Context[rootAppname], url, err, env.GetInheritCodeInfoAlls("panic.go", true))
				//w.WriteHeader(500)
				Page500(w, r, rootAppname)
			} else {
				if ActionMonitorExtendsObject != nil { // 如果实现了监听，就记录相关信息
					ActionMonitorExtendsObject.AddActionMonitor(url, appname, startTime.UnixNano(), endTime.UnixNano(), nil)
				}
				log4gohttp.InfoLog("interval : %6d ms . protocol: %s . client.Addr : %s . server.Host : %s  ,  server.context: %s ,  server.url : %s  .  ", time.Now().Sub(startTime)/1000/1000, r.Proto, r.RemoteAddr, r.Host, env.Env_Web_Context[rootAppname], url)
			}

			if env.Switch_CallChain {
				log4go.GetCallChain().SendCallchainChannel(log4go.NewCaller(log4go.GetGID(), "", time.Now()))
			}

		}()

		context := GetContext(w, r, rootAppname)
		if !authTicketServer_MS(context, url) { // 判断微服务的票据是否有效？
			log4go.ErrorLog(" auth failed : url=%s , context.RootAppName=%s ", url, context.RootAppname)
			return
		}

		if isRedirect { // 说明当前地址需要重定向
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.Redirect(w, r, "/"+env.Env_Web_Context[context.RootAppname]+url, http.StatusFound)
		}

		// 可以在这儿定义filter方法，渲染该连接
		if routeFitlerFunc != nil {
			if routeFitlerFunc[rootAppname] != nil {
				if !routeFitlerFunc[rootAppname](context) {
					return
				}
			}
		}

		if clazzType == 1 {
			vClazz := reflect.ValueOf(clazz)
			in0 := reflect.ValueOf(context)
			ins := []reflect.Value{in0}

			vClazz.Call(ins)
		} else if clazzType == 2 {
			vClazz := reflect.ValueOf(clazz)

			contentType := r.Header.Get("Content-Type")
			var params map[string][]string
			if strings.Contains(contentType, "multipart/form-data") {
				params = context.ParamWithMultipart.getParams_keyLower()
			} else {
				params = context.ParamWithSimple.getParams_keyLower()
			}

			in0 := reflect.Zero(reflect.TypeOf(clazz).In(0)) // 首先设置 nil

			if len(params) > 0 {
				in0 = reflect.New(reflect.TypeOf(clazz).In(0).Elem()) // 如果有数据就设置
			}

			webutils.DecodeHttpParams(params, in0)

			in1 := reflect.ValueOf(context)

			ins := []reflect.Value{in0, in1}

			vClazz.Call(ins)
		} else if clazzType == 3 {

			vClazz := reflect.ValueOf(clazz)

			mr, _ := context.Request.MultipartReader()

			var np *multipart.Part
			var errPart error

			params := make(map[string][]string)
			for true {
				np, errPart = mr.NextPart()
				if errPart == io.EOF {
					break
				} else if errPart != nil {
					log4go.ErrorLog(errPart)
					break
				} else {
					if np.Header.Get("Content-Type") == "application/octet-stream" {
						break
					} else {
						fieldName0 := np.FormName()
						var fieldValue0 string
						numBytes := 1024
						for true {
							bytesFieldValue0 := make([]byte, numBytes)
							nBytes, errBytes := np.Read(bytesFieldValue0)

							if errBytes == nil || errBytes == io.EOF {
								fieldValue0 = fmt.Sprintf("%s%s", fieldValue0, string(bytesFieldValue0[0:nBytes]))
								if errBytes == io.EOF {
									break
								}
							} else {
								if errBytes.Error() != "unexpected EOF" {
									log4go.ErrorLog(errBytes)
								}
								break
							}
						}
						fieldName0 = strings.ToLower(fieldName0)
						if params[fieldName0] == nil {
							params[fieldName0] = make([]string, 0, 1)
						}
						params[fieldName0] = append(params[fieldName0], fieldValue0)
					}
				}
			}

			in0 := reflect.Zero(reflect.TypeOf(clazz).In(0)) // 首先设置 nil

			if len(params) > 0 {
				in0 = reflect.New(reflect.TypeOf(clazz).In(0).Elem()) // 如果有数据就设置
			}

			webutils.DecodeHttpParams(params, in0)

			sph := &StreamParameterHandler{}
			sph.init = true
			sph.mr = mr
			sph.part = np

			in1 := reflect.ValueOf(sph)

			in2 := reflect.ValueOf(context)

			ins := []reflect.Value{in0, in1, in2}

			vClazz.Call(ins)

		}

	}
}

//控制器对象
type RouteActioner struct {
	Url         string      //访问的地址
	Action1     interface{} //设置控制器方法，方法定义是：func(context *Context)
	Action2     interface{} //设置控制器方法，方法定义是：func(paramBean interface{}, context *web.Context)
	Action3     interface{} //设置控制器方法，方法定义是：func(paramBean interface{},sph *web.StreamParameterHandler, context *web.Context)
	Appname     string      //app的名称
	RootAppname string      //apps的名称
}

//设置控制器对象对应的控制器处理方法
//
//@param clazz 控制器处理方法
func (ra *RouteActioner) Use(clazz interface{}) {
	mClazz := reflect.ValueOf(clazz)
	if mClazz.Kind() == reflect.Func {
		if mClazz.Type().NumIn() == 1 {
			if mClazz.Type().In(0).Kind() == reflect.Ptr {
				if mClazz.Type().In(0).Elem().Kind() == reflect.Struct &&
					mClazz.Type().In(0).Elem().Name() == "Context" {
					// 方法只有一个参数： *web.Context
					//ra.Action = clazz.(func(*Context))
					ra.Action1 = clazz
				}
			}
		} else if mClazz.Type().NumIn() == 2 {
			if mClazz.Type().In(1).Kind() == reflect.Ptr {
				if mClazz.Type().In(1).Elem().Kind() == reflect.Struct &&
					mClazz.Type().In(1).Elem().Name() == "Context" {
					if mClazz.Type().In(0).Kind() == reflect.Ptr {
						if mClazz.Type().In(0).Elem().Kind() == reflect.Struct {
							// 方法有两个参数： *ParamBean , *web.Context  ,其中 ：ParamBean可以任意名称，但是必须是指针结构体
							ra.Action2 = clazz
						}
					}
				}
			}
		} else if mClazz.Type().NumIn() == 3 {
			if mClazz.Type().In(2).Kind() == reflect.Ptr {
				if mClazz.Type().In(2).Elem().Kind() == reflect.Struct &&
					mClazz.Type().In(2).Elem().Name() == "Context" {
					if mClazz.Type().In(1).Kind() == reflect.Ptr {
						if mClazz.Type().In(1).Elem().Kind() == reflect.Struct &&
							mClazz.Type().In(1).Elem().Name() == "StreamParameterHandler" {
							if mClazz.Type().In(0).Kind() == reflect.Ptr {
								if mClazz.Type().In(0).Elem().Kind() == reflect.Struct {
									// 方法有两个参数： *ParamBean ,*web.StreamParameterHandler, *web.Context  ,
									// 其中 ：ParamBean可以任意名称，但是必须是指针结构体
									// 其中 ：StreamParameterHandler 是流处理
									ra.Action3 = clazz
								}
							}
						}
					}
				}
			}
		}

	}
	//ra.Action = clazz
}

/*
func (ra *RouteActioner) Use(clazz func(context *Context)) {
	ra.Action = clazz
}
*/

/*
创建控制器对象

@param  url  访问路径

@param  appname  app的名称

@return RouteActioner  控制器对象
*/
func AddRouteActioner(url string, appname string, rootAppName string) *RouteActioner {
	ra := new(RouteActioner)
	ra.Url = url
	ra.Appname = appname
	ra.RootAppname = rootAppName
	routeActionerList = append(routeActionerList, ra)
	return ra
}

/*
添加过滤器

@param  clazz  定义的过滤器方法
*/
func AddRouteFilter(clazz func(context *Context) bool, rootAppName string) {
	if routeFitlerFunc == nil {
		routeFitlerFunc = make(map[string]func(context *Context) bool)
	}
	routeFitlerFunc[rootAppName] = clazz
}

//展现层模板对象
type RouteTemplate struct {
	Alias       string // 模板别名，用于在Action中的Context的Forward中调用
	Path        string // gml文件路径
	Appname     string // app名称，位置在 src/${RootAppname}/${Appname}
	RootAppname string // 根路径，位置在 src/${RootAppname}
}

//设置模板路径
//
//@param  path  gml文件路径
func (rt *RouteTemplate) Tpl(path string) {
	rt.Path = path
}

//创建展现层模板对象
//
//@param rootappname 根名称，位置在 src/${rootappname}
//
//@param appname 应用名
//
//@param alias 模板别名，用于在Action中的Context的Forward中调用
func AddRouteTemplate(rootAppName string, appname string, alias string) *RouteTemplate {
	if appname == "" {
		appname = env.GetAppInfoFromRuntime(2)
	}

	rt := new(RouteTemplate)
	rt.Alias = alias
	rt.Appname = appname
	rt.RootAppname = rootAppName
	if routeTemplateList[rootAppName] == nil {
		routeTemplateList[rootAppName] = make([]*RouteTemplate, 0, 1)
	}
	routeTemplateList[rootAppName] = append(routeTemplateList[rootAppName], rt)
	return rt
}

//创建静态资源
//
//@param respath 静态资源文件夹，该文件夹下面的所有文件都会被注册，含子目录下。
//@param appName 应用根目录，位置在src/${appName}
func AddRouteStatics(respath string, rootAppName string) {
	if routeStaticsList[rootAppName] == nil {
		routeStaticsList[rootAppName] = make([]string, 0, 1)
	}
	routeStaticsList[rootAppName] = append(routeStaticsList[rootAppName], respath)
}

//创建404页面
//
//@param page404 404页面地址，必须是静态资源
func AddRoutePage404(page404 string, rootAppName string) {
	routePage404File[rootAppName] = page404
	path := env.Env_Project_Resource + "/" + rootAppName + routePage404File[rootAppName]
	b, err := ioutil.ReadFile(path)
	if err == nil {
		routePage404[rootAppName] = b
	} else {
		routePage404[rootAppName] = DEFAULT_ROUTE_PAGE_404
		log4go.ErrorLog(message.ERR_WEB0_39057, "404", path, err.Error())
	}

}

//创建500页面
//
//@param page500 500页面地址，必须是静态资源
func AddRoutePage500(page500 string, rootAppName string) {
	routePage500File[rootAppName] = page500
	path := env.Env_Project_Resource + "/" + rootAppName + routePage500File[rootAppName]
	b, err := ioutil.ReadFile(path)
	if err == nil {
		routePage500[rootAppName] = b
	} else {
		routePage500[rootAppName] = DEFAULT_ROUTE_PAGE_500
		log4go.ErrorLog(message.ERR_WEB0_39057, "500", path, err.Error())
	}
}

//创建appname的首页
//
//@param indexWebPath 首页地址，实际是一个web访问地址
func AddRouteIndex(indexWebPath string, rootAppName string) {
	routeIndexWebPath[rootAppName] = indexWebPath
	if strings.HasSuffix(routeIndexWebPath[rootAppName], ".html") || strings.HasSuffix(routeIndexWebPath[rootAppName], ".htm") { // 首页不为空，需要添加首页
		path := env.Env_Project_Resource + "/" + rootAppName + routeIndexWebPath[rootAppName]
		b, err := ioutil.ReadFile(path)
		if err == nil {
			routeIndex[rootAppName] = b
		} else {
			routeIndex[rootAppName] = make([]byte, 0)
			log4go.ErrorLog(message.ERR_WEB0_39057, "index(首页)", path, err.Error())
		}
	}
}

//创建根路径的首页
//
//@param indexWebPath 首页地址，实际是一个web访问地址
func AddRouteRootIndex(indexWebPath string, rootAppName string) {

	routeRootIndexWebPath[rootAppName] = indexWebPath
	if strings.HasSuffix(routeRootIndexWebPath[rootAppName], ".html") || strings.HasSuffix(routeRootIndexWebPath[rootAppName], ".htm") { // 首页不为空，需要添加首页
		path := env.Env_Project_Resource + "/" + rootAppName + routeRootIndexWebPath[rootAppName] //routeIndexWebPath[rootAppName]
		b, err := ioutil.ReadFile(path)
		if err == nil {
			routeRootIndex[rootAppName] = b
		} else {
			routeRootIndex[rootAppName] = make([]byte, 0)
			log4go.ErrorLog(message.ERR_WEB0_39057, "index(首页)", path, err.Error())
		}
	}

}

//
// 获取首页的web访问地址，一般使用forward定向
//
func GetIndexWebPath(rootAppname string) string {
	return routeIndexWebPath[rootAppname]
}

func GetMuxerObject() *Muxer {
	return muxerObject
}

/*
启动web服务
*/
func RegisterRouteAndStartServer() {
	prefix := make(map[string]string)
	for k, v := range env.Env_Web_Context {
		prefix[k] = v
	}
	muxerObject = NewMuxer(prefix, env.Env_Web_HttpsCertFile, env.Env_Web_HttpsKeyFile)
	for _, routeActioner := range routeActionerList {
		if routeActioner.Action1 != nil {
			muxerObject.AddDynamicHandler(routeActioner.Url, doAction(routeActioner.Action1, 1, routeActioner.Url, routeActioner.Appname, routeActioner.RootAppname, false), routeActioner.RootAppname)
			if routeIndexWebPath[routeActioner.RootAppname] == routeActioner.Url { // 首页不为空，需要添加首页
				muxerObject.AddDynamicHandler("/", doAction(routeActioner.Action1, 1, routeActioner.Url, "", routeActioner.RootAppname, true), routeActioner.RootAppname)
				muxerObject.AddDynamicHandler("", doAction(routeActioner.Action1, 1, routeActioner.Url, "", routeActioner.RootAppname, true), routeActioner.RootAppname)
			}

			if routeRootIndexWebPath[routeActioner.RootAppname] == routeActioner.Url { // 首页不为空，需要添加首页
				muxerObject.AddDynamicHandler("/", doAction(routeActioner.Action1, 1, routeActioner.Url, "", routeActioner.RootAppname, true), "")
				muxerObject.AddDynamicHandler("", doAction(routeActioner.Action1, 1, routeActioner.Url, "", routeActioner.RootAppname, true), "")
			}

		} else if routeActioner.Action2 != nil {
			muxerObject.AddDynamicHandler(routeActioner.Url, doAction(routeActioner.Action2, 2, routeActioner.Url, routeActioner.Appname, routeActioner.RootAppname, false), routeActioner.RootAppname)
			if routeIndexWebPath[routeActioner.RootAppname] == routeActioner.Url { // 首页不为空，需要添加首页
				muxerObject.AddDynamicHandler("/", doAction(routeActioner.Action2, 2, routeActioner.Url, "", routeActioner.RootAppname, true), routeActioner.RootAppname)
				muxerObject.AddDynamicHandler("", doAction(routeActioner.Action2, 2, routeActioner.Url, "", routeActioner.RootAppname, true), routeActioner.RootAppname)
			}

			if routeRootIndexWebPath[routeActioner.RootAppname] == routeActioner.Url { // 首页不为空，需要添加首页
				muxerObject.AddDynamicHandler("/", doAction(routeActioner.Action2, 2, routeActioner.Url, "", routeActioner.RootAppname, true), "")
				muxerObject.AddDynamicHandler("", doAction(routeActioner.Action2, 2, routeActioner.Url, "", routeActioner.RootAppname, true), "")
			}

		} else if routeActioner.Action3 != nil {
			muxerObject.AddDynamicHandler(routeActioner.Url, doAction(routeActioner.Action3, 3, routeActioner.Url, routeActioner.Appname, routeActioner.RootAppname, false), routeActioner.RootAppname)
		}
	}

	loadRouteTemplate()

	for appName, webPaths := range routeStaticsList {
		for _, webPath := range webPaths {
			muxerObject.AddStaticResourceInApps(webPath, webPath, appName)
		}
	}

	for _, rootAppName := range env.Env_Project_Resource_Apps {
		if strings.HasSuffix(routeIndexWebPath[rootAppName], ".html") || strings.HasSuffix(routeIndexWebPath[rootAppName], ".htm") { // 首页不为空，需要添加首页
			muxerObject.AddDynamicHandler("/", GlobalPageIndexFunc(rootAppName, true, false), rootAppName)
			muxerObject.AddDynamicHandler("", GlobalPageIndexFunc(rootAppName, true, false), rootAppName)
		}

		if strings.HasSuffix(routeRootIndexWebPath[rootAppName], ".html") || strings.HasSuffix(routeRootIndexWebPath[rootAppName], ".htm") { // 首页不为空，需要添加首页
			muxerObject.AddDynamicHandler("/", GlobalPageIndexFunc(rootAppName, true, true), "")
			muxerObject.AddDynamicHandler("", GlobalPageIndexFunc(rootAppName, true, true), "")
		}

	}

	if !env.Env_Web_Mode_IsProduct { // 如果是开发模式
		reloadRouteTemplate()
	}

	var wg sync.WaitGroup

	for i, _ := range env.Env_Server_Protocol_List {
		wg.Add(1)
		protocol := env.Env_Server_Protocol_List[i]
		addr := env.Env_Server_Addr_List[i]
		port := env.Env_Server_Port_List[i]
		timeout := env.Env_Server_Timeout_List[i]

		if protocol == "http" {
			go func() {
				defer wg.Add(-1)
				muxerObject.StartHttpServer(addr, port, timeout, protocol)
			}()
		} else if protocol == "https" {
			go func() {
				defer wg.Add(-1)
				muxerObject.StartHttpsServer(addr, port, timeout, protocol)
			}()
		} else if protocol == "fcgi_unix" {
			go func() {
				defer wg.Add(-1)
				os.Remove(addr)
				muxerObject.StartFcgiSocketServer("unix", addr, port, timeout, protocol)
			}()
		} else if protocol == "rpc_unix" {
			go func() {
				defer wg.Add(-1)
				os.Remove(addr)
				muxerObject.StartRpcSocketServer("unix", addr, port, timeout, false, protocol)
			}()
		} else if protocol == "rpc_tcp" {
			go func() {
				defer wg.Add(-1)
				os.Remove(addr)
				muxerObject.StartRpcSocketServer("tcp", addr, port, timeout, false, protocol)
			}()
		} else if protocol == "rpc_tcp_ssl" {
			go func() {
				defer wg.Add(-1)
				os.Remove(addr)
				muxerObject.StartHttpsServer(addr, port, timeout, protocol)
			}()
		} else {
			fmt.Printf("Protocol : %s  is error . \n", protocol)
			log4go.ErrorLog("Protocol : %s  is error . \n", protocol)
			wg.Add(-1)
		}

	}

	for _, rootAppName := range env.Env_Project_Resource_Apps {
		log4go.InfoLog("WebContext is initialized .   FolderRoot is  [ %16s ] . WebRoot is : [ %12s ] ", rootAppName, env.Env_Web_Context[rootAppName])
	}

	go func() {
		time.Sleep(600 * time.Second)
		muxerObject.ClearStaticResourceStore() // 清除缓存信息。
	}()

	wg.Wait()

}

/*
装载模板
*/
func loadRouteTemplate() {
	for rootappname, routeTemplates := range routeTemplateList {
		for _, routeTemplate := range routeTemplates {
			path := env.Env_Project_Resource + "/" + rootappname + routeTemplate.Path
			routeTemplateMap[path] = routeTemplate

			bytes, err := ioutil.ReadFile(path)
			if err == nil {
				t, err1 := template.New(path).Funcs(InitTemplateFuncs(routeTemplate.RootAppname, routeTemplate.Appname)).Parse(string(bytes))
				if err1 != nil {
					log4go.ErrorLog(message.ERR_WEB0_39059, path, err1)
				} else {
					templateDataMap[routeTemplate.RootAppname+"."+routeTemplate.Appname+"."+routeTemplate.Alias] = t
				}
			} else {
				log4go.ErrorLog(message.ERR_WEB0_39058, path, err.Error())
			}
		}

	}
}

/*
监控template文件变化情况并重新加载
*/
func reloadRouteTemplate() {
	w := NewMonitorFileEventWatcher()
	folderMap := make(map[string]string)
	for k, _ := range routeTemplateMap {
		folder := utils.GetParentDir(k)
		folderMap[folder] = folder
	}
	for k, _ := range folderMap {
		w.AddMonitorEventFileWatch(k)
	}

	go w.DoEvent(func(file string) {
		routeTemplate := routeTemplateMap[file]
		if routeTemplate != nil {
			path := env.Env_Project_Resource + "/" + routeTemplate.RootAppname + routeTemplate.Path

			bytes, err := ioutil.ReadFile(path)
			if err == nil {
				t, err1 := template.New(path).Funcs(InitTemplateFuncs(routeTemplate.RootAppname, routeTemplate.Appname)).Parse(string(bytes))
				if err1 != nil {
					log4go.ErrorLog(message.ERR_WEB0_39059, path, err1)
				} else {
					templateDataMap[routeTemplate.RootAppname+"."+routeTemplate.Appname+"."+routeTemplate.Alias] = t
				}
			} else {
				log4go.ErrorLog(message.ERR_WEB0_39058, path, err.Error())
			}

		}
	})

}
