// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/env"
	"net/http"
)

/*
404错误处理类

@param w

@param r
*/
func Page404(w http.ResponseWriter, r *http.Request, rootAppname string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(404)
	w.Write(routePage404[rootAppname])
}

func GlobalPage404Func(rootAppname string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(404)
		w.Write(routePage404[rootAppname])
	}
}

/*
500错误处理类

@param w

@param r
*/
func Page500(w http.ResponseWriter, r *http.Request, rootAppname string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(500)
	w.Write(routePage500[rootAppname])
}

/*
500错误处理类

@param w

@param r
*/
func Page401(w http.ResponseWriter, r *http.Request, rootAppname string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(401)
	w.Write(routePage401[rootAppname])
}

/*
静态文件首页展现

@param w

@param r
*/
func PageIndex(w http.ResponseWriter, r *http.Request, rootAppname string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write(routeIndex[rootAppname])
}

// isRedirect  是否重定向
// isRootIndex  是否是rootIndex
func GlobalPageIndexFunc(rootAppname string, isRedirect bool, isRootIndex bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if isRedirect { // 采用重定向方式
			path := ""
			if isRootIndex {
				path = routeRootIndexWebPath[rootAppname]
			} else {
				path = routeIndexWebPath[rootAppname]
			}
			http.Redirect(w, r, "/"+env.Env_Web_Context[rootAppname]+path, http.StatusFound)
		} else {
			w.WriteHeader(200)
			path := make([]byte, 0, 1)
			if isRootIndex {
				path = routeRootIndex[rootAppname]
			} else {
				path = routeIndex[rootAppname]
			}
			w.Write(path)
		}
	}
}
