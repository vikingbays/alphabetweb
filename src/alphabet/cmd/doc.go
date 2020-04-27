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
                     参数3: 平台定义，设置生成运行平台的代码，例如：windows/arm , darwin/amd64 等。
                           如果是多个用分号分割，例如：  windows/arm:darwin/amd64
                           具体平台信息可以通过命令查找：go tool dist list
                           补充：如果是android平台，需要设置 CC、CXX、CGO_ENABLED、GO111MODULE参数，在执行命令前定义(用AB_前缀)，例如：
                           GOOS=android  GOARCH=arm64    \
                           AB_CC=/xxxx/android-ndk-r19c/toolchains/llvm/prebuilt/darwin-x86_64/bin/aarch64-linux-android21-clang        \
                           AB_CXX=/xxxx/android-ndk-r19c/toolchains/llvm/prebuilt/darwin-x86_64/bin/aarch64-linux-android21-clang++     \
                           AB_CGO_ENABLED=1 AB_GO111MODULE=off     \
                           go run /xxxx/alphabetweb/src/alphabet/cmd/abserver.go  -build "..." "..." "android/arm64"

 	       -build_for_debug    string  string  string
	                 生成web服务启动二进制程序，支持gdb进行debug（ 在编译时增加参数 -gcflags "-N -l" ）。
	                 参数1：web工程根路径，源码存放在 src/apps下。
                     参数2：输出项目文件夹路径，默认可执行文件在此目录的bin目录下。
                     参数3: 参考 -build


	       -help
	                 展现帮助信息。
*/
package cmd
