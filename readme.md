# Vikingbays • AlphabetWeb

AlphabetWeb希望提供更加专业的，更加高效的，更加方便的开发架构，

AlphabetWeb是一个服务端开发框架，可用于搭建一个web网站，也可以构建复杂的微服务场景，也可用于中后台服务定制。AlphabetWeb采用Golang语言实现，具有一次开发多平台部署，运行效率更优。AlphabetWeb架构设计之初，不仅提供丰富的组件包，而且从开发者角度，注重开发过程的便捷性。

AlphabetWeb提供的架构图如下：


AlphabetWeb 采用MVC的web架构，提供分布式session管理（AlphabetWeb.Web）。提供多种组件：包括：数据库访问（AlphabetWeb.Sqler）、日志管理（AlphabetWeb.log4go）、国际化支持（AlphabetWeb.i18n）、缓存管理（AlphabetWeb.Cache）。服务访问支持标准的应用协议，包括：http/1.1，https，http/2，fcgi，unixDS等协议（AlphabetWeb.Protocol），同时也可以采用微服务方式访问（AlphabetWeb.MicroService），同时提供对应的客户端包访问（AlphabetWeb.Client）。AlphabetWeb 项目管理提供一套命令工具集、构建约定等（AlphabetWeb.Project），同时提供项目配置项管理（AlphabetWeb.Sysconfig）。AlphabetWeb 可以在Linux、Mac、Windows等平台上开发（AlphabetWeb.RunningOS），可运行在ARM、X86等架构上（AlphabetWeb.Architecture）。

## 一、特性/价值介绍

### 1. AlphabetWeb 提供丰富的套餐组合。
AlphabetWeb提供常用的服务端开发解决方案 ，包含了：1）、MVC 的web架构，2）、ORM数据访问、3）、统一日志管理、4）、集成第三方中间件框架，适配mysql、Postgresql、redis 等。

### 2. AlphabetWeb 支持友好的微服务场景。
AlphabetWeb提供微服务架构支持，微服务的开发和web应用开发无差异化，降低开发门槛。能够更好/更快的支持 web应用功能 向 微服务功能转移。

### 3. AlphabetWeb 采用Golang语言开发，很好支持server端开发。
golang语言先天性优势，更易于server端应用开发。与java相比，提供更好的性能保障，节约资源损耗。

### 4. http/2支持，实现多路复用。
AlphabetWeb提供基于http/2的协议模型，实现多路复用。有效的减少了TCP/IP连接数开销，提升并发访问性能。在Web场景、微服务场景都可适用。

### 5. 统一配置中心，使代码和环境配置分离。
AlphabetWeb提供统一配置中心，为应用系统提供多套配置项适配于不同的场景。例如：开发环境配置文件、测试环境配置文件、线上环境配置文件。使应用配置与代码工程解藕。

### 6. 支持访问调用链，跟踪服务的请求轨迹。
AlphabetWeb提供访问调用链管理，跟踪服务调用轨迹，有效支持多级微服务场景下的请求跟踪。

### 7. 跨平台支持，兼容性较好。
AlphabetWeb提供跨平台支持，支持windows、mac、linux、arm-linux 环境下运行。实现 一处开发，多平台部署。针对源码支持二进制打包，减少部署依赖（不需要运行环境），提升代码安全性保护。

### 8. UnixDS协议，更好的模块化分层。
UnixDS协议（Unix Domain Socket），通过进程协议模拟tcp/ip通信，更好的把一个大服务拆分成多个小服务，实现物理模块化分层，多项目隔离。

### 9. 模块化代码管理，采用目录约定模式，两级模块化管理。
AlphabetWeb的目录约定，第一级：子系统级（ 例如：/src/simple001 ），第二级：功能级（ 例如：/src/simple001/func001 ）。子系统级具有完整性和隔离性，可以单独部署或重新组合。易于复杂项目支撑。通过一个项目挂载多个子系统，使开发者便于在一个IDE中开发，开发时，可以打包成一个服务运行调试，测试/发布时，可以自动拆分成微服务部署架构。帮助开发者在复杂环境下（例如：多层微服务），提升开发的便捷性。


## 二、如何安装

### 1.可运行环境
Alphabet可运行在 Microsoft Windows ，Apple macOS ，Linux 等平台上。

### 2.安装golang
在 https://golang.google.cn 中下载golang安装包。根据不同的平台，选择不同的安装包。 下载完成后，解压到制定目录，例如：
```
/Users/vikingbays/golang/go
```

设置golang环境变量配置。

编辑系统环境变量文件：

```
vi ~/.bash_profile
```

在 .bash_profile文件中添加如下信息：

```
export GOROOT = "/Users/vikingbays/golang/go"
export PATH   = "$PATH:$GOROOT/bin"
export GOPATH = "/Users/vikingbays/golang/mygopath"
```

激活配置文件中环境变量：
```
. ~/.bash_profile
```
测试是否配置成功：
```
$ go version
返回：
  go version go1.xx.x darwin/amd64
```
说明：

  GOROOT 用于设置golang的安装路径。

  GOPATH 用于设置第三方应用路径，可以设置多个。

    通过go install/go get等工具获取的第三方包都放入到GOPATH设置的第一个路径中。

    GOPATH设置的路径是一个工程项目，里面主要包含三个目录: bin、pkg、src 。

### 3.安装AlphabetWeb
在github上下载AlphabetWeb安装包，下载后完成后，解压到指定目录，例如：

```
/Users/vikingbays/golang/AlphabetwebProject/alphabetweb
```
编辑系统环境变量文件：
```
vi ~/.bash_profile
```
添加环境变量：
```
export GOPATH = "$GOPATH:/Users/vikingbays/golang/mygopath"
```
激活配置文件中环境变量：
```
. ~/.bash_profile
```

### 4.安装IDE：Atom
在https://atom.io下载Atom安装包。根据不同的平台，选择不同的安装包。

下载并安装完成后启动Atom，打开Atom界面如下： Atom_start

选择Atom->Preferences->Install , 安装插件：go-plus ， file-icons ， atom-beautify 。界面如下： Atom_start

### 5.创建一个AlphabetWeb项目
一个AlphabetWeb项目的构成需要满足如下原则：

按照golang项目目录结构创建一级目录： bin、logs、pkg、upload、src 。

一个AlphabetWeb 可以包含多个应用包集合，他在src的下一级目录。

应用包集合：可以理解为URL地址的一级目录。他的结构是：

```
src
  glass     # 访问地址：  http://ip:port/glass/...
  stone     # 访问地址：  http://ip:port/stone/...
  books     # 访问地址：  http://ip:port/books/...
创建一个项目（glassProject），他的目录结构如下：

  glassProject                        # 项目名称
     bin                                
     logs                             # 日志存储路径
     pkg
     upload
     src                              # 源码存放路径
        glass                         # 应用包集合，包含了 app1,app2,... 等应用。
                                      # 可支持一个src下定义多个应用包集合。
           app1                       # 一个应用，可包含多个功能，具有独立体系的代码。
                                      # 包括：视图、控制器、查询sql、服务等。
           app2
           ...
           dbsconfig.toml             # (可选)数据库配置文件，包括可注册的数据库，例如：app1中注册的pg_test1 。
           logconfig.toml             # (必选)日志配置文件。
           logconfig_http.toml        # (必选)http日志配置文件。记录请求响应时长
           webconfig.toml             # (必选)web配置文件。
           cachesconfig.toml          # (可选)cache配置文件。
           ms_server_config.toml      # (可选)微服务的服务端配置文件 。
           ms_client_config_01.toml   # (可选)微服务的服务端配置文件，可以是多个 _02 _03 _04 ....
```
