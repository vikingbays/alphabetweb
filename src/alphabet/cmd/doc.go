// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*
alphabet的命令管理工具，包括：webapp应用启动、停止等。

具体命令解释如下。
	参考帮助：

	命令：  abserver   [command]

	command说明：

	       -start    string
	                 启动web服务。
	                 参数：web工程根路径，源码存放在 src/apps下。

	       -stop    string
	                 停止web服务。
	                 参数：web工程根路径，源码存放在 src/apps下。

	       -genmain    string
	                 生成web服务启动代码，代码存储到 “${project}/src/apps”下。
	                 参数：web工程根路径，源码存放在 src/apps下。

	       -build    string  string  string
	                 生成web服务启动二进制程序。
	                 参数1：web工程根路径，源码存放在 src/apps下。
                     参数2：输出项目文件夹路径，默认可执行文件在此目录的bin目录下。
                     参数3: 平台定义，如果不传参数，表示就是当前平台，
                            其中：｀windows｀  ：表示windows平台。
                                 ｀mac｀      ：表示mac平台。
                                 ｀linux｀    ：表示linux平台。
                                 ｀linuxarm｀ ：表示linux的arm平台。

 	       -build_for_debug    string  string  string
	                 生成web服务启动二进制程序，支持gdb进行debug（ 在编译时增加参数 -gcflags "-N -l" ）。
	                 参数1：web工程根路径，源码存放在 src/apps下。
                     参数2：输出项目文件夹路径，默认可执行文件在此目录的bin目录下。
                     参数3: 平台定义，如果不传参数，表示就是当前平台，
                            其中：｀windows｀  ：表示windows平台。
                                 ｀mac｀      ：表示mac平台。
                                 ｀linux｀    ：表示linux平台。
                                 ｀linuxarm｀ ：表示linux的arm平台。


	       -help
	                 展现帮助信息。
*/
package cmd
