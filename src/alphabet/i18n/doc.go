// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*

实现国际化功能，主要针对每个app下的国际化文件解析 （ 例如：i18n/zh.toml ）。
在整个运行环境中，避免了多个app下命名冲突的问题。
假设 : 在 /app1/i18n/zh.toml 和 /app2/i18n/zh.toml 中都定义了相同的msgctxt和msgid信息，msgstr信息不同，
在app1中的应用只会读取app1国际化的定义，在app2中的应用只会读取app2国际化的定义.虽然获取消息的key（msgctxt和msgid）重名，但是还是各自取自身的信息，互不影响。

当前国际化的使用方式主要有三种：

1. 在webapp的控制器中直接使用，绑定到context上下文中，例如：
  func HelloWorld(context *web.Context) {
    ...
    context.SetI18nSession("en")                         ## 设置使用哪种语言，当前设置是为｀en｀，表示使用该app下的en.toml的国际化信息。
                                                         ## 可以不设置，就使用webapp启动时默认的语言环境
    context.GetLocale().Get("姓名", "helloworld")         ## 获取msgctxt="helloworld",msgid="姓名" 的信息
    ...
  }

2. 在gml模板文件中使用，例如：
  ## 使用Locale 标签定义，第一个参数 msgid，第二个参数是 msgctxt，第三个参数是语言，固定使用变量 ｀.I18n｀ 。
  <div>
    {{Locale "姓名" "helloworld" .I18n  }} : <input type='text' name='Name' value=''>
  </div>

3. 通用使用方式，先设置语言环境，再获取国际化信息，例如：
  locale := i18n.GetLocale("en")                        ## 实例化一个｀en｀的国际化对象
  locale().Get("姓名", "helloworld")                     ## 获取msgctxt="helloworld",msgid="姓名" 的信息

国际化文件命名：(可根据语言环境，不断扩展)
  en.toml   --  表示英文
  zh.toml   --  表示中文

文件内容格式：

  [[msg]]
  msgctxt="msgContext"   # 消息上下文
  msgid="myId1"           # 消息编码
  msgstr="msg content"   # 消息内容

  [[msg]]
  msgctxt="msgContext"   # 消息上下文
  msgid="myId2"           # 消息编码
  msgstr="msg content"   # 消息内容

  [[msg]]
  msgctxt="msgContext"   # 消息上下文
  msgid="myId3"           # 消息编码
  msgstr="msg content"   # 消息内容


*/
package i18n
