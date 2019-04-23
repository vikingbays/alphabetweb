# Vikingbays • AlphabetWeb

AlphabetWeb希望提供更加专业的，更加高效的，更加方便的开发架构，

AlphabetWeb是一个服务端开发框架，可用于搭建一个web网站，也可以构建复杂的微服务场景，也可用于中后台服务定制。AlphabetWeb采用Golang语言实现，具有一次开发多平台部署，运行效率更优。AlphabetWeb架构设计之初，不仅提供丰富的组件包，而且从开发者角度，注重开发过程的便捷性。

AlphabetWeb提供的架构图如下：


AlphabetWeb 采用MVC的web架构，提供分布式session管理（AlphabetWeb.Web）。提供多种组件：包括：数据库访问（AlphabetWeb.Sqler）、日志管理（AlphabetWeb.log4go）、国际化支持（AlphabetWeb.i18n）、缓存管理（AlphabetWeb.Cache）。服务访问支持标准的应用协议，包括：http/1.1，https，http/2，fcgi，unixDS等协议（AlphabetWeb.Protocol），同时也可以采用微服务方式访问（AlphabetWeb.MicroService），同时提供对应的客户端包访问（AlphabetWeb.Client）。AlphabetWeb 项目管理提供一套命令工具集、构建约定等（AlphabetWeb.Project），同时提供项目配置项管理（AlphabetWeb.Sysconfig）。AlphabetWeb 可以在Linux、Mac、Windows等平台上开发（AlphabetWeb.RunningOS），可运行在ARM、X86等架构上（AlphabetWeb.Architecture）。


## 特性/价值介绍

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
