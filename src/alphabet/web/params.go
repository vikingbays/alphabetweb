// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/env"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"alphabet/mux"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	UPLOAD_MEMCACHESIZE = "UPLOAD_MEMCACHESIZE"
	UPLOAD_MAXSIZE      = "UPLOAD_MAXSIZE"
	UPLOAD_STOREPATH    = "UPLOAD_STOREPATH"
	UPLOAD_SPLITONAPPS  = "UPLOAD_SPLITONAPPS"
)

/*
获取表单的参数信息，不能处理multipart表单
*/
func GetSimpleFormParams(r *http.Request) map[string][]string {
	paramsUsedMap, _ := getSimpleFormParams_nest(r, false)
	return paramsUsedMap
}

func getSimpleFormParams_nest(r *http.Request, keyToLower bool) (map[string][]string, map[string][]string) {
	paramsUsedMap := make(map[string][]string)
	paramsUsedMap2 := make(map[string][]string)

	urlParamsMap := mux.Vars(r)
	for key, value := range urlParamsMap {
		paramsUsedMap2[key] = []string{value}
		if keyToLower {
			key = strings.ToLower(key)
		}
		paramsUsedMap[key] = []string{value}
	}
	err := r.ParseForm()
	if err != nil {
		log4go.ErrorLog(message.ERR_WEB0_39048, err)
	} else {
		for key, values := range r.Form {
			paramsUsedMap2[key] = values
			if keyToLower {
				key = strings.ToLower(key)
			}
			paramsUsedMap[key] = values
		}
	}

	return paramsUsedMap, paramsUsedMap2
}

/*
获取表单的参数信息，该表单必须定义为：method="post" enctype="multipart/form-data" 。
如果需要修改文件上传的参数信息，可以通过SetUploadFileOptions方法设置
*/
func GetMultipartFormParams(r *http.Request, rootAppname string) map[string][]string {
	return getMultipartFormParams_nest(r, false, rootAppname)
}

/*

@Param keyToLower  是否设置 key 转小写
*/
func getMultipartFormParams_nest(r *http.Request, keyToLower bool, rootAppname string) map[string][]string {
	paramsUsedMap := make(map[string][]string)

	urlParamsMap := mux.Vars(r)
	for key, value := range urlParamsMap {
		if keyToLower {
			key = strings.ToLower(key)
		}
		paramsUsedMap[key] = []string{value}
	}

	storeFolder, memcachesize, _ := getUploadFileOptions(r, rootAppname)

	r.ParseMultipartForm(memcachesize)
	if r.MultipartForm != nil {

		for key, values := range r.MultipartForm.Value {
			if keyToLower {
				key = strings.ToLower(key)
			}
			paramsUsedMap[key] = values
		}

		if storeFolder != "" {
			for key, files := range r.MultipartForm.File {
				if keyToLower {
					key = strings.ToLower(key)
				}
				paramsUsedMap[key] = make([]string, 0)
				for _, fileHandler := range files {
					file, err := fileHandler.Open()
					if err != nil {
						log4go.ErrorLog(message.ERR_WEB0_39049, err)
					}
					defer file.Close()

					filename := ""

					if runtime.GOOS == "windows" {
						contentDisposition := fileHandler.Header.Get("Content-Disposition")
						startPos := strings.Index(contentDisposition, "filename=\"")
						if startPos != -1 {
							filenamePart := contentDisposition[(startPos + 10):]
							endPos := strings.Index(filenamePart, "\"")
							if endPos != -1 {
								filenamePart = filenamePart[:endPos]
								startPos = strings.LastIndex(filenamePart, "\\")
								if startPos != -1 {
									filename = filenamePart[(startPos + 1):]
								}
							}
						}
						if filename == "" {
							filename = "temp_" + time.Now().Format("2006-01-02_15.04.05.000000000")
						}
					} else {
						filename = fmt.Sprintf("u%d_%s", time.Now().UnixNano(), fileHandler.Filename)
					}

					f, err := os.OpenFile(storeFolder+"/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
					if err != nil {
						log4go.ErrorLog(message.ERR_WEB0_39050, storeFolder, filename, err)
					}
					_, err2 := io.Copy(f, file)
					if err2 != nil {
						log4go.InfoLog(message.ERR_WEB0_39051, storeFolder, filename, err2)
					}
					defer f.Close()
					paramsUsedMap[key] = append(paramsUsedMap[key], storeFolder+"/"+filename)
				}
			}
		}

	}

	return paramsUsedMap
}

/*
如果需要修改文件上传的参数信息，可以通过该方法设置

@param r http.Request指针对象

@param memcachesize 缓冲区大小

@param maxsize 最大上传文件大小

@param storepath 存储的根目录

@param splitonapps  是否需要按照应用目录分出上传文件
*/
func SetUploadFileOptions(r *http.Request, memcachesize int64, maxsize int64, storepath string, splitonapps bool) {

	r.Header.Set(UPLOAD_MEMCACHESIZE, strconv.FormatInt(memcachesize, 10))
	r.Header.Set(UPLOAD_MAXSIZE, strconv.FormatInt(maxsize, 10))
	r.Header.Set(UPLOAD_STOREPATH, storepath)
	r.Header.Set(UPLOAD_SPLITONAPPS, strconv.FormatBool(splitonapps))

}

//清除上传文件选项，清除的内容包括：缓存大小、文件大小、存储路径、是否按appname分类等。
//
//@param r
func ClearUploadFileOptionsOfRequest(r *http.Request) {
	r.Header.Set(UPLOAD_MEMCACHESIZE, "")
	r.Header.Set(UPLOAD_MAXSIZE, "")
	r.Header.Set(UPLOAD_STOREPATH, "")
	r.Header.Set(UPLOAD_SPLITONAPPS, "")
}

//获取上传文件选项
//
//@param r
//
//@return storefolder  存储路径
//
//@return memcachesize  缓存大小
//
//@return maxsize  文件大小
//
func getUploadFileOptions(r *http.Request, rootAppname string) (storefolder string, memcachesize int64, maxsize int64) {
	memcachesizeStr := r.Header.Get(UPLOAD_MEMCACHESIZE)
	maxsizeStr := r.Header.Get(UPLOAD_MAXSIZE)
	storepathStr := r.Header.Get(UPLOAD_STOREPATH)
	splitonappsStr := r.Header.Get(UPLOAD_SPLITONAPPS)

	ClearUploadFileOptionsOfRequest(r)

	if memcachesizeStr == "" {
		memcachesize = env.Env_Web_Upload_Memcachesize[rootAppname]
	} else {
		i64, err := strconv.ParseInt(memcachesizeStr, 10, 0)
		if err == nil {
			memcachesize = i64
		} else {
			memcachesize = env.Env_Web_Upload_Memcachesize[rootAppname]
		}
	}

	if maxsizeStr == "" {
		maxsize = env.Env_Web_Upload_Maxsize[rootAppname]
	} else {
		i64, err := strconv.ParseInt(maxsizeStr, 10, 0)
		if err == nil {
			maxsize = i64
		} else {
			maxsize = env.Env_Web_Upload_Maxsize[rootAppname]
		}
	}

	var splitonappsBool bool = env.Env_Web_Upload_Splitonapps[rootAppname]
	if splitonappsStr == "true" {
		splitonappsBool = true
	} else if splitonappsStr == "false" {
		splitonappsBool = false
	}

	if storepathStr == "" {
		storepathStr = env.Env_Web_Upload_Storepath[rootAppname]
	}

	if splitonappsBool {
		var appPath string = GetAppNameFromUrl(r, rootAppname)
		storefolder = storepathStr + "/" + appPath
	} else {
		storefolder = storepathStr
	}

	_, err := os.Stat(storefolder)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(storefolder, 0700)
		} else {
			log4go.ErrorLog(message.ERR_WEB0_39052, storefolder, err)
		}
	}

	folder, err := os.Stat(storefolder)

	if !folder.IsDir() {
		log4go.ErrorLog(message.ERR_WEB0_39053, storefolder)
		storefolder = ""
	}
	return
}

/*
根据url获取对应的应用名称

@param r

*/
func GetAppNameFromUrl(r *http.Request, rootAppname string) string {
	var appName string = ""
	url := r.URL.String()
	offsetPos := 0
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		offsetPos = 2
	}
	urlSplits := strings.Split(url, "/")
	if env.Env_Web_Context[rootAppname] != "" {
		if len(urlSplits) >= 3+offsetPos {
			appName = urlSplits[2+offsetPos]
		}
	} else {
		if len(urlSplits) >= 2+offsetPos {
			appName = urlSplits[1+offsetPos]
		}
	}
	if appName != "" { // 判断是否是一个有效的appName
		if strings.Contains(appName, "?") || strings.Contains(appName, "=") {
			appName = ""
		}
	}

	if appName == "" {
		appName = env.GetAppInfoFromRuntime(3)
	}

	return appName
}
