// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/env"
	"alphabet/i18n"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"alphabet/mock"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

//
// 提供国际化方法，用于html/template 模板中。
// 在模板中只需要如下使用方式：
//  {{Locale "许愿" "app1.helloworld" .I18n  }}
//
//  其中 .I18n 在返回html/template 模板对象中就已经封装，优先取Session信息，没有取全局信息
//
func Locale(id string, context string, i18nType string) string {
	if i18nType == "" {
		i18nType = env.Env_Web_I18n
	}

	locale := i18n.GetLocale(i18nType)
	if locale != nil {
		str := locale.GetWithoutAppname(id, context)
		if str != "" {
			return str
		}

	}

	return id
}

/*
提供页面嵌入的方法，用于html/template 模板中。其中html标签不换转义
在模板中只需要如下使用方式：
{{ IncludeHTML "hello/helloworld" "friend=moon"  .Params  .Header  "POST" }}
其中 .Params 是父级页面请求时参数

  @Param urlStr  当前webapp服务的路由地址，可以带参数，例如：/hello/init?p1=v1&p2=v2
  @Param params0 请求参数，字符串形式，例如：p1=v1&p2=v2
  @Param params1 数组形式的参数，一般从上下文中获取， 例如： .Param
  @Param header1 HttpHeader信息，也记录了cookie数据，一般从上下文中获取， 例如： .Header
  @Param method  POST方式或GET方式。

  @Return 返回不转义的HTML数据
*/
func IncludeHTML(urlStr string, params0 string, params1 map[string][]string,
	header1 http.Header, method string, rootAppname string, appname string) template.HTML {
	return template.HTML(include(urlStr, params0, params1, header1, method, rootAppname, appname))
}

/*
提供页面嵌入的方法，用于html/template 模板中。其中html标签直接转义
在模板中只需要如下使用方式：
{{IncludeHTML "../hello/helloworld" "friend=moon&like=music" .Params .Header "POST"  }}
其中 .Params 是父级页面请求时参数

  @Param urlStr  当前webapp服务的路由地址，可以带参数，例如：/hello/init?p1=v1&p2=v2
  @Param params0 请求参数，字符串形式，例如：p1=v1&p2=v2
  @Param params1 数组形式的参数，一般从上下文中获取， 例如： .Param
  @Param header1 HttpHeader信息，也记录了cookie数据，一般从上下文中获取， 例如： .Header
  @Param method  POST方式或GET方式。

  @Return 返回转义后的数据
*/
func IncludeText(urlStr string, params0 string, params1 map[string][]string,
	header1 http.Header, method string, rootAppname string, appname string) string {
	return include(urlStr, params0, params1, header1, method, rootAppname, appname)
}

func include(urlStr string, params0 string, params1 map[string][]string,
	header1 http.Header, method string, rootAppname string, appname string) string {
	u, _ := url.Parse(urlStr)
	if u.Host == "" {
		path := u.Path
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		if !strings.HasPrefix(path, "/"+env.Env_Web_Context[rootAppname]) {
			path = "/" + env.Env_Web_Context[rootAppname] + path
		}
		route := GetMuxerObject().GetRouter().GetRoute(path)
		if route != nil {
			handleServHTTP := route.GetHandler()
			if handleServHTTP != nil {
				code, respBody := mock.MockWebAction(method, path, params0, params1, header1, handleServHTTP.ServeHTTP)
				if code != 404 && code != 500 {

					return respBody
				} else if code == 404 {
					log4go.ErrorLog(message.ERR_WEB0_39054, urlStr, code, path)
					return string(routePage404[rootAppname])
				} else if code == 500 {
					log4go.ErrorLog(message.ERR_WEB0_39054, urlStr, code, path)
					return string(routePage500[rootAppname])
				}
			} else {
				log4go.ErrorLog(message.ERR_WEB0_39055, urlStr, path)
				return string(routePage500[rootAppname])
			}
		} else {
			log4go.WarnLog(message.WAR_WEB0_69004, urlStr)
			return ""
		}

	} else {
		log4go.WarnLog(message.WAR_WEB0_69005, urlStr)
		return string(routePage500[rootAppname])
	}

	log4go.ErrorLog(message.ERR_WEB0_39056, urlStr)

	return string(routePage500[rootAppname])
}

/*
 使数据的HTML标签 不转义

  @param data  原始数据

  @Return  返回原始数据 不转义 HTML标签
*/
func ParseHTML(data string) template.HTML {
	return template.HTML(data)
}

//
// 自定义模板函数，当前定义了 Locale 方法用于处理国际化
//
func InitTemplateFuncs(rootAppname string, appname string) template.FuncMap {
	return template.FuncMap{"Locale": func(id string, context string, i18nType string) string {
		return Locale(id, rootAppname+"."+appname+"."+context, i18nType)
	}, "IncludeHTML": func(urlStr string, params0 string, params1 map[string][]string, header1 http.Header, method string) template.HTML {
		return IncludeHTML(urlStr, params0, params1, header1, method, rootAppname, appname)
	}, "IncludeText": func(urlStr string, params0 string, params1 map[string][]string, header1 http.Header, method string) string {
		return IncludeText(urlStr, params0, params1, header1, method, rootAppname, appname)
	}, "ParseHTML": func(data string) template.HTML {
		return ParseHTML(data)
	}}
}
