// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/core/utils"
	"alphabet/env"
	"alphabet/i18n"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"alphabet/sessions"
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

/*
Web上下文，封装了请求、响应、会话、请求头、参数、页面跳转等信息。

用于控制器的上下文定义：
   func Action1(context *Context){
      ......
   }

用于过滤器的上下文定义：
   func Filter(context *Context) bool{
      ......
   }

*/
type Context struct {
	Request            *http.Request              //请求
	Response           http.ResponseWriter        //响应
	Session            sessions.AlphaSession      //会话
	Header             http.Header                //请求头
	ParamWithMultipart *ParamWithMultipartContext //Multipart的Form表单处理，例如：上传表
	ParamWithSimple    *ParamWithSimpleContext    //普通表单处理
	Return             *ReturnContext             //页面跳转处理，包含：forword，redirect，json，download
	RootAppname        string                     // 当前上下文对应的rootapp根位置，位置在 src/${RootAppname}
}

//
// 获取当前访问的url地址
//
func (context *Context) GetCurrentUrl() string {
	return context.Request.URL.String()
}

//
// 获取当前所在的app名称信息
//
func (context *Context) GetCurrentAppname() string {
	url := context.GetCurrentUrl()
	if len(url) == 0 || url == "/" {
		return ""
	}
	pos := 0
	if env.Env_Web_Context[context.RootAppname] != "" {
		pos = strings.Index(utils.Substr(url, 1, len(url)-1), "/")
	}
	lenPos := strings.Index(utils.Substr(url, pos+2, len(url)-pos-2), "/")
	if lenPos == -1 {
		return utils.Substr(url, pos+2, len(url)-pos-2)
	} else {
		return utils.Substr(url, pos+2, lenPos)
	}
}

//
// 在全局环境变量中设置I18n
//
func (context *Context) SetI18nGlobal(i18nType string) {
	if i18nType != "" {
		env.Env_Web_I18n = i18nType
	}
}

//
// 在Session中设置I18n，优先于SetI18nGlobal设置的信息，优先使用。
//
func (context *Context) SetI18nSession(i18nType string) {
	if i18nType != "" {
		context.Session.Set("I18n", i18nType)
	}
}

//
// 获取I18n访问对象Locale，如果Session中定义I18n信息，就使用Session的，如果没有，就使用全局定义的
//
func (context *Context) GetLocale() *i18n.Locale {
	return i18n.GetLocale(context.getI18nType())

}

//获取当前I18n的语言信息。
func (context *Context) getI18nType() string {
	lang := ""
	if context.Session != nil {
		langtmp := context.Session.Get("I18n")
		if langtmp != "" {
			lang = langtmp
		}
	}

	if lang != "" {
		return lang
	} else {
		return env.Env_Web_I18n
	}

}

/*
 Multipart的Form表单处理上下文
*/
type ParamWithMultipartContext struct {
	Request  *http.Request       //请求
	Response http.ResponseWriter //响应
	Datas    map[string][]string //存储表单数据
	context  *Context            //Web上下文
}

/*
如果需要修改文件上传的参数信息，可以通过该方法设置，在调用GetParams方法前执行。
如果不设置，默认使用系统的参数配置。

@params memcachesize 缓冲区大小

@params maxsize 最大上传文件大小

@params storepath 存储的根目录

@params splitonapps  是否需要按照应用目录分出上传文件
*/
func (p *ParamWithMultipartContext) SetUploadFileOptions(memcachesize int64,
	maxsize int64, storepath string, splitonapps bool) {
	SetUploadFileOptions(p.Request, memcachesize, maxsize, storepath, splitonapps)
}

/*

获取表单参数

@return map[string][]string  表单数据，如果有上传文件，那么会纪录本地文件存储地址

*/
func (p *ParamWithMultipartContext) GetParams() map[string][]string {
	p.Datas = GetMultipartFormParams(p.Request, p.context.RootAppname)
	return p.Datas
}

func (p *ParamWithMultipartContext) getParams_keyLower() map[string][]string {
	p.Datas = getMultipartFormParams_nest(p.Request, true, p.context.RootAppname)
	return p.Datas
}

/*
 普通Form表单处理上下文
*/
type ParamWithSimpleContext struct {
	Request  *http.Request       //请求
	Response http.ResponseWriter //响应
	Datas    map[string][]string //存储表单数据
	context  *Context            //Web上下文
}

/*
 获取表单参数

 @return map[string][]string  表单数据，如果有上传文件，那么会纪录本地文件存储地址

*/
func (p *ParamWithSimpleContext) GetParams() map[string][]string {
	p.Datas = GetSimpleFormParams(p.Request)
	return p.Datas
}

func (p *ParamWithSimpleContext) getParams_keyLower() map[string][]string {
	var returnDatas map[string][]string
	returnDatas, p.Datas = getSimpleFormParams_nest(p.Request, true)
	return returnDatas
}

/*
 页面跳转的上下文，用于处理 转发、重定向、Json、下载等操作
*/
type ReturnContext struct {
	Request          *http.Request       // 请求
	Response         http.ResponseWriter // 响应
	IsContainSession bool                // 是否包含Session信息
	IsContainParams  bool                // 是否包含参数数据
	context          *Context            // Web上下文
}

/*
设置Forword数据是否包含 Session 和表单参数Params信息

@param isContainSession , 如果 true ，表示包含Session数据。

@param isContainParams , 如果 true ，表示包含表单参数Params数据。

*/
func (r *ReturnContext) SetForwardDataType(isContainSession bool, isContainParams bool) {
	r.IsContainSession = isContainSession
	r.IsContainParams = isContainParams
}

/*
使用Forward方式进行页面跳转

@param aliasOfPath  转发到具体的页面（使用html/template），该页面别名，在route中定义，例如：
  web.AddRouteTemplate("action1_run1_page1").Tpl("/app1/view/action1_run1_page1.gml")
  那么，aliasOfPath="action1_run1_page1"

@param data  转发的数据，在具体页面中加载（使用html/template）
*/
func (r *ReturnContext) Forward(aliasOfPath string, data interface{}) {
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	returnDataStore := ReturnDataStore{}
	if r.IsContainSession {
		returnDataStore.Session = r.context.Session.Maps()
	}
	if r.IsContainParams {
		if r.context.ParamWithSimple.Datas != nil && len(r.context.ParamWithSimple.Datas) > 0 {
			returnDataStore.Params = r.context.ParamWithSimple.Datas
		} else if r.context.ParamWithMultipart.Datas != nil && len(r.context.ParamWithMultipart.Datas) > 0 {
			returnDataStore.Params = r.context.ParamWithMultipart.Datas
		}
	}
	returnDataStore.Header = r.Request.Header
	returnDataStore.Datas = data
	returnDataStore.I18n = r.context.getI18nType()
	returnDataStore.WebRoot = env.Env_Web_Context[r.context.RootAppname]

	appname := GetAppNameFromUrl(r.Request, r.context.RootAppname)
	t := templateDataMap[r.context.RootAppname+"."+appname+"."+aliasOfPath]
	if t == nil {
		t = templateDataMap["."+aliasOfPath]
	}
	if t == nil {
		log4go.ErrorLog(message.ERR_WEB0_39045, aliasOfPath, appname)
		Page404(r.Response, r.Request, r.context.RootAppname)
	} else {
		err := t.Execute(r.Response, returnDataStore)
		if err != nil {
			log4go.ErrorLog(message.ERR_WEB0_39046, aliasOfPath, appname, err)
			Page500(r.Response, r.Request, r.context.RootAppname)
		}
	}

}

/*
使用Redirect方式进行页面重定向

@param url  重定向的页面

*/
func (r *ReturnContext) Redirect(url string) {
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(r.Response, r.Request, "/"+env.Env_Web_Context[r.context.RootAppname]+url, http.StatusFound)
}

//
// Forward 到404页面
//
func (r *ReturnContext) Forward404() {
	Page404(r.Response, r.Request, r.context.RootAppname)
}

//
// Forward 到500页面
//
func (r *ReturnContext) Forward500() {
	Page500(r.Response, r.Request, r.context.RootAppname)
}

//
// Forward 到401页面
//
func (r *ReturnContext) Forward401() {
	Page401(r.Response, r.Request, r.context.RootAppname)
}

//
// Redirect 到404页面
//
func (r *ReturnContext) Redirect404() {
	r.Redirect(routePage404File[r.context.RootAppname])
}

//
// Redirect 到500页面
//
func (r *ReturnContext) Redirect500() {
	r.Redirect(routePage500File[r.context.RootAppname])
}

/*
使用Json方式返回数据

@param data  json数据对象

*/
func (r *ReturnContext) Json(data interface{}) {
	r.Response.Header().Set("Content-Type", "text/json; charset=utf-8") // json header
	r.Response.WriteHeader(http.StatusOK)
	b, err := json.Marshal(data)
	if err != nil {
		log4go.ErrorLog(message.ERR_WEB0_39047, err)
		io.WriteString(r.Response, "")
	} else {
		io.WriteString(r.Response, string(b))
	}
}

/*
使用Text方式返回字符串数据

@param text  字符串数据

*/
func (r *ReturnContext) Text(text string) {
	r.Response.Header().Set("Content-Type", "text/json; charset=utf-8") // json header
	r.Response.WriteHeader(http.StatusOK)
	io.WriteString(r.Response, text)
}

//
// 下载文件操作，通过本地文件路径下载
//
// @param filePath 需要下载的文件路径
//
// @param aliasName 下载文件别名
//
func (r *ReturnContext) DownloadFile(filePath string, aliasName string) {
	fi, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	reader := bufio.NewReader(fi)
	r.StreamDownBufferIO(reader, aliasName)
}

//
// 通过数据流方式传输数据，基于bufio.Reader
//
// @param reader 数据流
//
// @param aliasName 下载文件别名
//
func (r *ReturnContext) StreamDownBufferIO(reader *bufio.Reader, aliasName string) {
	r.Response.Header().Set("Content-Type", "application/octet-stream")
	r.Response.Header().Set("Content-Disposition", "attachment;filename="+aliasName)
	r.Response.WriteHeader(http.StatusOK)

	buf := make([]byte, 8192)

	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		r.Response.Write(buf[:n])
	}
}

//
// 页面返回数据封装
//
type ReturnDataStore struct {
	Session map[string]string   //session数据
	Header  http.Header         //http请求头
	Params  map[string][]string //表单数据
	Datas   interface{}         //返回数据
	I18n    string              //国际化信息
	WebRoot string              //web的根路径,例如： http://ip:port/web1/xxxxxx ， 那么 WebRoot=web1
}

//
//初始化上下文
//
//@param w
//
//@param r
//
//@return *Context  初始化上下文
//
func GetContext(w http.ResponseWriter, r *http.Request, rootAppname string) *Context {
	context := &Context{}
	context.Response = w
	context.Request = r
	context.Header = w.Header()
	context.Session = GetSession(w, r)
	context.ParamWithMultipart = &ParamWithMultipartContext{r, w, nil, context}
	context.ParamWithSimple = &ParamWithSimpleContext{r, w, nil, context}
	context.Return = &ReturnContext{r, w, true, true, context}
	context.RootAppname = rootAppname
	return context
}
