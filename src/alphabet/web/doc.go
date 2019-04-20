// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*
web应用开发，按模块化管理，采用MVC框架实现，提供：web上下文管理、session管理、过滤器、控制器和展现模版等功能。

web应用包创建，是按应用分为不同独立的文件夹，每个文件夹中都可以包含完整的web应用要素，包含，过滤器定义，展现模板定义等。
具体一个webapp的完整目录结构可以参考｀alphabet｀的代码结构描述。

一个应用包的结构是：
   glassProject
    ...
    src
      glass                             // 应用包集合，包含了 app1,app2,app3 等应用。
        app1                            // 一个应用包，实现该应用的所有功能，具有独立的国际化、控制器、展现等。
          config
            db
              pg_test1.toml             // 配置pg_test1库的所有sql语句信息。[详细解释如下]
              ...
          i18n
            en.toml                     // 针对 en 的国际化配置。
            zh.toml                     // 针对 zh 的国际化配置。
            ...
          module
          resource
            css                         // css文件夹，包含：css文件，可以多个。
              app1.css
              ...
            img                         // img文件夹
            js                          // javascript文件夹
          view                          // 展现层文件，类似于jsp或者php文件，采用html/template模板实现，必须以gml为后缀名
            helloworld_index.gml
            helloworld_gift.gml
            ...
          helloworld.go                 // 控制器处理类，必须在应用包app1下定义。实现func func1(c *web.Context){} 方法处理，方法名可自定义。
          restful.go
          ...
          route.toml                    // 该应用包的路由配置，包括注册：控制器、展现层、过滤器(全局)等

其中最关键的要素是：
  glassProject/src/glass/app1/route.toml                  // 用于定义路由配置
  glassProject/src/glass/app1/helloworld.go               // 创建一个控制器
  glassProject/src/glass/app1/view/helloworld_index.gml   // 创建一个展现层模板文件

-------------------------------------------------------------------------------

1.路由器说明

路由文件 route.toml 配置的语法，例如：

  ＃ [[action]] , [[gml]]               可以在每个app应用包中定义，在每个app应用包中定义多个。
  ＃ [[static]]                         可以在每个app应用包中定义，在每个app应用包中只能定义一个。
  ＃ [[filter]] , [[404]] , [[500]]     在整个webapp中只能定义一个。

  # helloworld 的处理

  [[action]]               # 注册一个控制器
  HelloWorld="helloworld"  # ｀HelloWorld｀ 实际上是一个控制器方法，具体语法结构是：func (*web.Context){} ，例如：
                           #                func HelloWorld(context *web.Context) {
                           #                  ...
                           #                }
                           # ｀HelloWorld｀，该方法必须在appname根目录下定义方法，例如存放在：app1/helloworld.go中。
                           # 如果设置 HelloWorld="/app3/helloworld" 那么不需要加appname
                           # 如果设置 HelloWorld="helloworld" 那么需要加appname ，判断依据是否／开始


  [[gml]]                                         ＃ 注册一个模板
  helloworld_index="view/helloworld_index.gml"
  helloworld_gift="view/helloworld_gift.gml"      ＃ ｀helloworld_gift｀ 是一个模板别名。在控制器处理页面跳转时候被调用，例如：
                                                  ＃                 context.Forward("helloworld_gift")
                                                  ＃ ｀view/helloworld_gift.gml｀是该模板文件路径，存储在app1/view/helloworld_gift.gml

  # upload 的处理

  [[action]]
  Upload="uploadfile"

  [[gml]]
  upload_index="view/upload_index.gml"
  upload_success="view/upload_success.gml"


  # json 的处理

  [[action]]
  JsonInfo="jsoninfo/{min}/{max}"

  [[gml]]
  jsoninfo_notfound="view/jsoninfo_notfound.gml"


  # download 的处理

  [[action]]
  Download="download"
  DoDownload="do_download"

  [[gml]]
  download_index="view/download_index.gml"

  # 静态资源加载

  [[static]]
  name="resource/"   # ｀name｀该命名固定。
                     # ｀resource/｀资源文件夹目录，如果是多个，逗号分隔。

  # 过滤器设置

  [[filter]]                 # 整个webapp只能定义一次
  name="ServFilter"          #｀name｀该命名固定。
                             # `ServFilter` 是 filter 过滤器方法，具体语法结构是：func (*web.Context) bool{} ，例如：
                             #              func ServFilter(context *web.Context) bool {
                             #                ...
                             #              }
                             # 在过滤器方法中，如果返回值是 true，表示继续执行 控制器方法，如果为false，就直接跳出，不再执行后续的控制器方法。

  # 404页面设置

  [[404]]                    # 整个webapp只能定义一次
  name="resource/404.html"   # 定义404资源文件

  # 500页面设置

  [[500]]                    # 整个webapp只能定义一次
  name="resource/500.html"   # 定义500资源文件

  # 首页

  [[index]]
  name="/app1/resource/datas.html"     # 可以配置一个.html／.html结尾的静态文件
                                       # 也可以配置一个action动态页面，例如： /app1/helloworld ，这个在action中已定义


  # 启动器，用于启动时执行
  [[starter]]
  StartRun=1            # 格式：[方法名]=[等级]   ， 说明：设置方法执行等级，1为最高级，99为最低级，最高级先执行。
                        # StartRun 方法不带参数，也没有返回值，并且该方法的定义是在当前app的根目录下

-------------------------------------------------------------------------------

2.控制器说明

控制器的实现方法建议存放在appname的目录下，例如：app1/helloworld.go 。
控制器的语法结构是：
  func 方法名(context *web.Context) {
    ...
  }
例如：
  func HelloWorld(context *web.Context) {
  	context.SetI18nSession("en")
  	params := context.ParamWithSimple.GetParams()
  	if len(params) == 0 {
  		context.Return.Forward("helloworld_index", nil)
  	} else {
  		if params["I18n"] != nil {
  			context.SetI18nSession(params["I18n"][0])
  		}
  		if params["giftName"] != nil && params["giftName"][0] != "" {
  			giftInfo := new(GiftInfo)
  			giftInfo.GName = params["giftName"][0]
  			giftInfo.GNumber, _ = strconv.Atoi(params["giftNumber"][0])
  			giftInfo.GFriend = params["friends"]
  			context.Return.Forward("helloworld_gift", giftInfo)
  		} else {
  			context.Return.Forward("helloworld_index", nil)
  		}

  	}
  }

其中，传入参数 web.Context ，可以获取到如下信息：(假设实例对象是：context)
  context.Request                          // 获取http.Request对象
  context.Response                         // 获取http.Request对象
  context.Session                          // 获取Session对象
            Set(...)                               // 设置session信息，key，value方式
            Get(...)                               // 根据key获取session信息
            Delete(...)                            // 根据key删除session信息
            Clear(...)                             // 清除该session所有信息
  context.Header                           // 获取HttpHeader头对象
  context.ParamWithMultipart               // 获取表单对象，Multipart的Form表单处理，例如：上传表
            SetUploadFileOptions(...)              // 设置当前文件上传所使用的配置信息，如果没有设置，就使用系统默认的配置信息。
            GetParams()                            // 获取表单信息，此时文件已上传到上传目录中
  context.ParamWithSimple                  // 获取表单对象，普通表单处理
            GetParams()                            // 获取表单信息
  context.Return                           // 获取页面展现处理对象，页面跳转处理，包含：forword，redirect，json，download
            SetForwardDataType(...)                // 设置Forword数据是否包含 Session 和表单参数Params信息
            Forward(...)                           // 使用Forward方式进行页面跳转
            Redirect(...)                          // 使用Redirect方式进行页面重定向
            Json(...)                              // 返回json数据
            DownloadFile(...)                      // 下载文件操作，通过本地文件路径下载
            DownloadBufferIO(...)                  // 下载文件操作，通过bufio.Reader流下载
            Forward404()                           // Forward 到404页面
            Forward500()                           // Forward 到500页面
            Redirect404()                          // Redirect 到404页面
            Redirect500()                          // Redirect 到500页面
  context.GetCurrentUrl()                  // 获取当前访问的url地址
  context.GetCurrentAppname()              // 获取当前所在的app名称信息
  context.SetI18nGlobal(...)               // 在全局环境变量中设置I18n
  context.SetI18nSession(...)              // 在Session中设置I18n，优先于SetI18nGlobal设置的信息，优先使用。
  context.GetLocale()                      // 获取I18n访问对象Locale，如果Session中定义I18n信息，就使用Session的，如果没有，就使用全局定义的
  context.getI18nType()                    // 获取当前I18n的语言信息。

-------------------------------------------------------------------------------

3.展现层模板说明

展现层模板是基于html/template模板实现的。
根据html/template 的数据模型，都是按照树状结构存储数据对象的，根对象访问以｀.｀开始。
在模板中定义的有效根对象有：
  .Session     # Session数据
  .Header      # http请求头
  .Params      # 表单数据
  .Datas       # 业务数据 ，在业务处理过程中产生的结果数据
  .I18n        # 国际化信息
具体在｀type ReturnDataStore｀结构体中封装。
  ## 获取.Datas的属性Id数据
  <div> {{ .Datas.Id }} </div>

模版还提供丰富的表达式：

3.1)、if关键字，用于判断是否满足条件
  ## 判断对象是否存在
  {{if .Datas.queryFlag }}
  ...
  {{end}}
  ## 判断对象是否为true ，相等条件是用 eq
  {{ if eq .Datas.createFlag true}}
  ...
  {{end}}

3.2)、range关键字，用于遍历list对象
  ## 针对.Datas的userList数据列表遍历
  {{range .Datas.userList}}
      <tr>
        <td>{{.Usrid}}</td>
        <td>{{.Name}}</td>
        <td>{{.Nanjing}}</td>
        <td>{{.Money}}</td>
      </tr>
  {{end}}


展现层模板提供自定义方法的封装，例如：国际化方法的Locale。

3.A1)、Locale关键字，用于国际化封装，可以在gml模板文件中使用，例如：
  ## 使用Locale 标签定义，第一个参数 msgid，第二个参数是 msgctxt，第三个参数是语言，固定使用变量 ｀.I18n｀ 。
  <div>
    {{Locale "姓名" "helloworld" .I18n  }} : <input type='text' name='Name' value=''>
  </div>
具体可参考国际化alphabet/i18n包的帮助。

3.A2)、IncludeHTML/IncludeText 关键字，用于引入其他请求页。

当前只支持 webapp应用内部请求页引入，跨服务的请求暂不支持。
  ## 在该div区域载入/web2/hello/helloworld请求页，传递参数"friend=moon" 和 .Params ，传递http请求头.Header，采用POST方式
  <div>
     {{ IncludeHTML "/web2/hello/helloworld" "friend=moon"  .Params  .Header  "POST" }}
  </div>

3.A3)、ParseHTML，使数据的HTML标签不转义。
  ## 输出红色字体的 “Hello football!”
  {{ ParseHTML "<div style='color:red'>Hello football!</div>" }}
-------------------------------------------------------------------------------

4.日志访问的说明

直接使用log4go包的日志方法即可，日志使用非常简单，首先是引入alphabet.log4go包，然后根据日志级别直接调用如下方法：
  log4go.FinestLog(arg0 interface{}, args ...interface{})    ## 记录Finest级别日志，最低

  log4go.FineLog(arg0 interface{}, args ...interface{})      ## 记录Fine级别日志

  log4go.DebugLog(arg0 interface{}, args ...interface{})     ## 记录Debug级别日志

  log4go.TraceLog(arg0 interface{}, args ...interface{})     ## 记录Trace级别日志

  log4go.InfoLog(arg0 interface{}, args ...interface{})      ## 记录Info级别日志

  log4go.WarnLog(arg0 interface{}, args ...interface{})      ## 记录Warn级别日志

  log4go.ErrorLog(arg0 interface{}, args ...interface{})     ## 记录Error级别日志，最高

具体可参考alphabet/log4go包

-------------------------------------------------------------------------------

5.数据库访问的说明

具体可参考alphabet/sqler包

*/
package web
