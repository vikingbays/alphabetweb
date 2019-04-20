// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/core/utils"
	"alphabet/env"
	"alphabet/log4go"
	"alphabet/mux"
	"strings"
	//"alphabet/web/longserver"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"net/http/pprof"
	"time"
)

type Muxer struct {
	router            *mux.Router
	Prefix            map[string]string
	CertFile          string
	KeyFile           string
	staticResourceMap map[string]map[string]StaticResourceOfMuxer
}

type StaticResourceOfMuxer struct {
	webPath     string
	resPath     string
	rootAppname string
}

func (m *Muxer) AddDynamicHandler(webPath string, f func(http.ResponseWriter,
	*http.Request), rootAppname string) {
	m.router.HandleFunc(m.Prefix[rootAppname]+webPath, f)
}

/*
设置静态资源文件访问，例如：javascript , css , html , jpg 等

@param webPath 表示web访问路径，例如： /app1/resource/js/jquery.js

@param resPath 表示资源文件路径，例如：/GoProject/webcontainer/src/apps/app1/resource/js/jquery.js
*/
func (m *Muxer) AddStaticResource(webPath string, resPath string, rootAppname string) {
	/*  无法支持多服务  */
	//http.Handle(m.Prefix+webPath, http.StripPrefix(m.Prefix+webPath, FileServer(Dir(resPath))))
	if m.staticResourceMap[rootAppname] == nil {
		m.staticResourceMap[rootAppname] = make(map[string]StaticResourceOfMuxer)
	}
	m.staticResourceMap[rootAppname][webPath] = StaticResourceOfMuxer{webPath, resPath, rootAppname}
}

func (m *Muxer) ClearStaticResourceStore() {
	m.staticResourceMap = nil
}

/*
设置静态资源文件夹访问，例如：javascript , css , html , jpg 等

@param webPath 表示web访问路径，例如： /app1/resource/js/jquery.js

@param resPathInApps 表示资源文件路径，例如：/app1/resource/js/jquery.js

@param rootAppName  表示 应用根路径，他的位置在 src/${rootAppName}
*/
func (m *Muxer) AddStaticResourceInApps(webPath string, resPathInApps string, rootAppName string) {
	if utils.ExistFile(env.Env_Project_Resource + "/" + rootAppName + "/" + resPathInApps) {
		m.AddStaticResource(webPath, env.Env_Project_Resource+"/"+rootAppName+"/"+resPathInApps, rootAppName)
	}
}

/*
  添加pprof信息
*/
func (m *Muxer) addWebpprof(mux *http.ServeMux, rootAppName string) {
	if env.Env_Pprof {

		mux.Handle("/debug/pprof/", http.HandlerFunc(pprofFunc("Index", rootAppName)))
		mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprofFunc("Cmdline", rootAppName)))
		mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprofFunc("Profile", rootAppName)))
		mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprofFunc("Symbol", rootAppName)))
		mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprofFunc("Trace", rootAppName)))

	}
}

func pprofFunc(funcName string, rootAppName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		paramMap := GetSimpleFormParams(r)

		cliIp := strings.Split(r.RemoteAddr, ":")[0]
		ipFlag := false
		if len(env.Env_Pprof_ClientIP) > 0 {
			for _, ip0 := range env.Env_Pprof_ClientIP {
				if ip0 == cliIp {
					ipFlag = true
					break
				}

			}
		} else {
			ipFlag = true
		}

		if !ipFlag { // 如果不符合验证的IP
			Page500(w, r, rootAppName)
			return
		}

		if env.Env_Pprof_Username != "" {
			usr := ""
			pwd := ""
			usrs := paramMap["username"]
			pwds := paramMap["password"]
			if len(usrs) > 0 {
				usr = usrs[0]
			}
			if len(pwds) > 0 {
				pwd = pwds[0]
			}
			if usr == "" {
				c, err := r.Cookie("pprof_username")
				if err == nil {
					usr = c.Value
				}
			}
			if pwd == "" {
				c, err := r.Cookie("pprof_password")
				if err == nil {
					pwd = c.Value
				}
			}
			if env.Env_Pprof_Username != usr || env.Env_Pprof_Password != pwd {
				Page500(w, r, rootAppName)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "pprof_username",
				Value:    usr,
				Path:     "/",
				HttpOnly: false,
			})
			http.SetCookie(w, &http.Cookie{Name: "pprof_password", Value: pwd,
				Path:     "/",
				HttpOnly: false})
		}
		if funcName == "Index" {
			pprof.Index(w, r)
		} else if funcName == "Cmdline" {
			pprof.Cmdline(w, r)
		} else if funcName == "Profile" {
			pprof.Profile(w, r)
		} else if funcName == "Symbol" {
			pprof.Symbol(w, r)
		} else if funcName == "Trace" {
			pprof.Trace(w, r)
		}
	}
}

/*
启动http的web服务，必须在 AddXxxxxXxxxx 方法后执行
*/
func (m *Muxer) StartHttpServer(ipAddr string, ipPort int, timeout int, protocol string) bool {

	flag_StartServer := true

	mux := http.NewServeMux()
	mux.Handle("/", m.router)

	for _, staticResourceSubMap := range m.staticResourceMap {
		for webPath, staticResource := range staticResourceSubMap {
			mux.Handle(m.Prefix[staticResource.rootAppname]+webPath,
				http.StripPrefix(m.Prefix[staticResource.rootAppname]+webPath, FileServer(Dir(staticResource.resPath), staticResource.rootAppname)))
		}
	}

	firstRootAppname := ""
	for _, rootAppName := range env.Env_Project_Resource_Apps {
		if rootAppName == "" {
			firstRootAppname = rootAppName
		}
	}

	m.addWebpprof(mux, firstRootAppname)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", ipPort),
		Handler: mux,
	}

	if timeout > 0 {
		server.ReadTimeout = time.Duration(timeout) * time.Second
		server.WriteTimeout = time.Duration(timeout) * time.Second
	}

	m.printServerStartingLog(protocol, ipAddr, ipPort)
	err := server.ListenAndServe() // 如果启动成功，将被阻塞
	if err != nil {
		flag_StartServer = false
		m.printServerFailedddLog(protocol, ipAddr, ipPort, err.Error())
	}
	return flag_StartServer

}

func stripPrefix(prefix string, h http.Handler, rootAppName string) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			r.URL.Path = p
			h.ServeHTTP(w, r)
		} else {
			Page404(w, r, rootAppName)
		}
	})
}

//启动https的web服务，必须在 AddXxxxxXxxxx 方法后执行
func (m *Muxer) StartHttpsServer(ipAddr string, ipPort int, timeout int, protocol string) bool {

	flag_StartServer := true

	mux := http.NewServeMux()
	mux.Handle("/", m.router)

	for _, staticResourceSubMap := range m.staticResourceMap {
		for webPath, staticResource := range staticResourceSubMap {
			mux.Handle(m.Prefix[staticResource.rootAppname]+webPath,
				http.StripPrefix(m.Prefix[staticResource.rootAppname]+webPath, FileServer(Dir(staticResource.resPath), staticResource.rootAppname)))
		}
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", ipPort),
		Handler: mux,
	}
	if timeout > 0 {
		server.ReadTimeout = time.Duration(timeout) * time.Second
		server.WriteTimeout = time.Duration(timeout) * time.Second
	}
	m.printServerStartingLog(protocol, ipAddr, ipPort)
	err := server.ListenAndServeTLS(m.CertFile, m.KeyFile)
	if err != nil {
		flag_StartServer = false
		m.printServerFailedddLog(protocol, ipAddr, ipPort, err.Error())
	}
	return flag_StartServer
}

//启动unix socket的web服务，必须在 AddXxxxxXxxxx 方法后执行
func (m *Muxer) StartFcgiSocketServer(protoType string, ipAddr string, ipPort int, timeout int, protocol string) bool {

	flag_StartServer := true

	addr, err := net.ResolveUnixAddr(protoType, ipAddr)
	if err != nil {
		flag_StartServer = false
		log4go.ErrorLog("Cannot resolve unix addr: %s", err.Error())
		return flag_StartServer
	}

	listener, err := net.ListenUnix(protoType, addr)
	if err != nil {
		flag_StartServer = false
		log4go.ErrorLog("Cannot listen to unix domain socket: %s", err.Error())
		return flag_StartServer
	}

	mux := http.NewServeMux()
	mux.Handle("/", m.router)

	for _, staticResourceSubMap := range m.staticResourceMap {
		for webPath, staticResource := range staticResourceSubMap {
			mux.Handle(m.Prefix[staticResource.rootAppname]+webPath,
				http.StripPrefix(m.Prefix[staticResource.rootAppname]+webPath, FileServer(Dir(staticResource.resPath), staticResource.rootAppname)))
		}
	}

	m.printServerStartingLog(protocol, ipAddr, ipPort)
	err2 := fcgi.Serve(listener, mux)
	if err2 != nil {
		flag_StartServer = false
		m.printServerFailedddLog(protocol, ipAddr, ipPort, err2.Error())
	}

	return flag_StartServer
}

// 长链接访问
func (m *Muxer) StartRpcSocketServer(protoType string, ipAddr string, ipPort int, timeout int, sslFlag bool, protocol string) bool {
	flag_StartServer := true

	mux := http.NewServeMux()
	mux.Handle("/", m.router)

	for _, staticResourceSubMap := range m.staticResourceMap {
		for webPath, staticResource := range staticResourceSubMap {
			mux.Handle(m.Prefix[staticResource.rootAppname]+webPath,
				http.StripPrefix(m.Prefix[staticResource.rootAppname]+webPath, FileServer(Dir(staticResource.resPath), staticResource.rootAppname)))
		}
	}

	server := http.Server{
		Handler: mux,
	}
	if timeout > 0 {
		server.ReadTimeout = time.Duration(timeout) * time.Second
		server.WriteTimeout = time.Duration(timeout) * time.Second
	}

	addr := ipAddr
	if protoType == "tcp" {
		addr = fmt.Sprintf(":%d", ipPort)
	}
	if sslFlag {
		crt, err0 := tls.LoadX509KeyPair(m.CertFile, m.KeyFile)
		if err0 != nil {
			log4go.ErrorLog(err0.Error())
		}
		tlsConfig := &tls.Config{}
		tlsConfig.Certificates = []tls.Certificate{crt}
		tlsConfig.Time = time.Now
		tlsConfig.Rand = rand.Reader
		listener, err := tls.Listen(protoType, addr, tlsConfig)
		if err != nil {
			flag_StartServer = false
			log4go.ErrorLog(err.Error())
			return flag_StartServer
		}
		m.printServerStartingLog(protocol, ipAddr, ipPort)
		err2 := server.Serve(listener)
		if err2 != nil {
			flag_StartServer = false
			m.printServerFailedddLog(protocol, ipAddr, ipPort, err2.Error())
			return flag_StartServer
		}
	} else {
		listener, err := net.Listen(protoType, addr)
		if err != nil {
			flag_StartServer = false
			log4go.ErrorLog(err.Error())
			return flag_StartServer
		}
		m.printServerStartingLog(protocol, ipAddr, ipPort)
		err2 := server.Serve(listener)
		if err2 != nil {
			flag_StartServer = false
			m.printServerFailedddLog(protocol, ipAddr, ipPort, err2.Error())
		}
	}
	return flag_StartServer
}

func (m *Muxer) printServerStartingLog(protocol string, addr string, port int) {
	log4go.InfoLog("Server is initialized .   Protocol is  [ %12s ] . Addr is : [ %18s ] . Port is : [ %5d ] . ", protocol, addr, port)
	fmt.Printf("Server is initialized .   Protocol is  [ %11s ] . Addr is : [ %11s ] . Port is : [ %5d ] . \n", protocol, addr, port)
}

func (m *Muxer) printServerFailedddLog(protocol string, addr string, port int, errInfo string) {
	log4go.ErrorLog("!!!Server is failed  when starting .   Protocol is  [ %12s ] . Addr is : [ %18s ] . Port is : [ %5d ] . error: %s", protocol, addr, port, errInfo)
}

func (m *Muxer) GetRouter() *mux.Router {
	return m.router
}

/*
初始化路由器，用于启动http服务

@param prefix : 设置webapp前缀

@return Muxer  创建路由器对象
*/
func NewMuxer(prefix map[string]string, cert string, key string) *Muxer {
	firstRootAppname := ""
	for key, p := range prefix {
		if firstRootAppname == "" {
			firstRootAppname = key
		}
		if p != "" {
			prefix[key] = "/" + p
		}
	}
	var muxer *Muxer = new(Muxer)
	muxer.router = mux.NewRouter()
	muxer.Prefix = prefix
	muxer.CertFile = cert
	muxer.KeyFile = key
	muxer.router.NotFoundHandler = http.HandlerFunc(GlobalPage404Func(firstRootAppname)) //注册全局404页面
	muxer.staticResourceMap = make(map[string]map[string]StaticResourceOfMuxer)
	return muxer
}
