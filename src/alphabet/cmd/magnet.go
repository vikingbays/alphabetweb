// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package main

import (
	"alphabet/core/utils"
	"alphabet/magnet"
	"fmt"
	"os"
)

func main() {
	Run()
}

func Run() {
	lenOfArgs := len(os.Args)
	if lenOfArgs == 3 {
		cmdType := os.Args[1]
		cmdParams := os.Args[2]

		err := validateCmdMagnet(cmdType, cmdParams)
		if err != nil {
			fmt.Println(err)
		} else {
			doCmdMagnet(cmdType, cmdParams)
		}
	} else {
		helperMagnet()
	}
}

func doCmdMagnet(cmdType string, params ...string) (err error) {
	if cmdType == "-start" {

		magnet.Init(params[0])

		magnet.StartMagnetServer()

	} else if cmdType == "-stop" {

	}

	return

}

func validateCmdMagnet(cmdType string, params ...string) (err error) {
	err = nil
	if cmdType == "-start" || cmdType == "-stop" {
		if len(params) == 1 && utils.IsFolder(params[0]) {

		} else if len(params) == 1 {
			err = fmt.Errorf("该命令 [%s] 传入的参数不正确，需要传入的是工程根目录文件夹路径。（该路径无效：%s）", cmdType, params[0])
		} else {
			err = fmt.Errorf("该命令 [%s] 只接收1个参数。（当前传入参数个数为%d，不正确）", cmdType, len(params))

		}
	} else {
		err = fmt.Errorf("指定的命令信息不正确，没有该命令：%s 。", cmdType)
	}
	return
}

func helperMagnet() {
	str := `
参考帮助：

命令：  magnet   [command]

command说明：(所有文件路径的分隔符都必须使用｀／｀)

       -start    string
                   启动magent服务。
                   参数1：web工程根路径，配置文件在 ／config下。

       -stop     string
                   停止magnet服务。
                   参数1：web工程根路径，配置文件在 ／config下。


`
	fmt.Println(".")
	fmt.Printf(str)

}
