// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package i18n

import (
	"alphabet/core/toml"
	"alphabet/env"
	"fmt"
)

/**

Locale.Get(id,context)
Locale.contextMap[context][id]
*/

const Locale_Context_Default_value = "__"

type Locale struct {
	Lang  string
	Datas map[string]map[string]string
}

type Locale_Msg_Toml_Config struct {
	Msg []Locale_Msg
}

type Locale_Msg struct {
	Id      string `toml:"msgid"`
	Str     string `toml:"msgstr"`
	Context string `toml:"msgctxt"`
}

func (this *Locale) String() (format string) {
	format = fmt.Sprintf("Locale : \n   Lang  : %s  \n   Datas : %v \n", this.Lang, this.Datas)
	return
}

/*
初始化Locale 的Datas对象

*/
func (this *Locale) Init() {
	this.Datas = make(map[string]map[string]string)
}

/*
装载国际化文件

@param lang  语言名称

@param path  国际化文件地址

@param appname  应用名，作为域名，避免重名现象

*/
func (this *Locale) Load(lang, path string, appname string) {
	var msgTomlConfig Locale_Msg_Toml_Config
	if _, err := toml.DecodeFile(path, &msgTomlConfig); err != nil {
		fmt.Println(err)
	} else {
		for _, msg := range msgTomlConfig.Msg {
			if len(msg.Context) == 0 {
				msg.Context = Locale_Context_Default_value
			}
			msg.Context = appname + "." + msg.Context
			if this.Datas[msg.Context] == nil {
				msgMap := make(map[string]string)
				msgMap[msg.Id] = msg.Str
				this.Datas[msg.Context] = msgMap
			} else {
				msgMap := this.Datas[msg.Context]
				msgMap[msg.Id] = msg.Str
			}
		}

		this.Lang = lang

	}
}

/*
获取国际化内容

@param id  对应配置文件的 msgid 信息

@param context  对应配置文件的 msgctxt 信息

@return str 返回内容，对应配置文件的 msgstr 信息

*/
func (this *Locale) Get(id string, context string) (str string) {
	if context == "" {
		context = Locale_Context_Default_value
	}

	context = env.GetAppInfoFromRuntime(2) + "." + context

	str = this.Datas[context][id]
	return
}

/*
获取国际化内容，与Get方法不同的是，context内容需要预先加入appname名称，例如：｀app1.helloworld｀

@param id  对应配置文件的 msgid 信息，传入参数需要带上appname名称，例如：｀app1.helloworld｀

@param context  对应配置文件的 msgctxt 信息

@return str 返回内容，对应配置文件的 msgstr 信息

*/
func (this *Locale) GetWithoutAppname(id string, context string) (str string) {
	if context == "" {
		context = Locale_Context_Default_value
	}
	str = this.Datas[context][id]
	return
}
