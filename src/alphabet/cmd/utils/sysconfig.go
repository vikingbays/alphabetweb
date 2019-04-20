// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package utils

import (
	"alphabet/core/utils"
	"alphabet/env"
	"alphabet/log4go"
	"alphabet/service"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func DoSysconfigPathForUrl(sysconfigUrl string) {
	u, err := url.Parse(sysconfigUrl)
	if err != nil {
		fmt.Println("param [-config_url] is error. err: " + err.Error())
		log4go.ErrorLog("param [-config_url] is error. err: ", err.Error())
	}
	port, _ := strconv.Atoi(u.Port())
	sc := &service.SimpleRpcClient{Protocol: "rpc_tcp", Ip: u.Hostname(), Port: port}
	err = sc.Connect()
	if err != nil {
		fmt.Println(err.Error())
		log4go.ErrorLog(err.Error())
	}
	defer sc.Close()
	aliasKey := ""
	reqStr := u.RequestURI()
	reqPos := strings.LastIndex(reqStr, "/")
	if reqPos != -1 {
		if (reqPos + 1) >= len(reqStr) {
			reqPos2 := strings.LastIndex(reqStr[0:reqPos], "/")
			if reqPos2 != -1 {
				reqPos2 = reqPos2 + 1
				aliasKey = reqStr[reqPos2:reqPos]
			}
		} else {
			reqPos = reqPos + 1
			aliasKey = reqStr[reqPos:]
		}
	}
	srcDirs, _ := ioutil.ReadDir(env.Env_Project_Root + "/src")
	for _, dir := range srcDirs {
		dirName := dir.Name()
		if !strings.HasPrefix(dirName, ".") {
			if aliasKey != "" {
				os.RemoveAll(env.Env_Project_Root + "/src/" + dirName + "/sysconfig." + aliasKey)
			}
			srcFile := env.Env_Project_Root + "/src/" + dirName + "/sysconfig.zip"
			err := sc.DoDownload(u.RequestURI()+"/"+dirName, "", srcFile)
			if err != nil {
				fmt.Println("err:::", err.Error())
				log4go.ErrorLog(err.Error())
			} else {
				dest := env.Env_Project_Root + "/src/" + dirName
				deRootFolder := utils.DeCompressZip(srcFile, dest)
				os.RemoveAll(srcFile)
				if strings.HasPrefix(deRootFolder, "sysconfig.") {
					log4go.InfoLog("Update sysconfig from '-config_url' , path is [%s/%s] . ", dest, deRootFolder)
					fmt.Printf("abserver.starting....  config ::::: Update sysconfig from '-config_url' , path is [%s/%s] . \n", dest, deRootFolder)
					env.SetSysconfigPath(deRootFolder[10:])
				}

			}
		}

	}
	//port, _ := strconv.Atoi(u.Port())
	//sc := &service.SimpleRpcClient{Protocol: "rpc_tcp", Ip: u.Hostname(), Port: port}
	//sc.DoDownload(u.RequestURI(), "", env.Env_Project_Root+"/")
}

func DoSysconfigPathForKey(sysconfigKey string) {
	env.SetSysconfigPath(sysconfigKey)
}
