// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package i18n

import (
	"alphabet/env"
	"alphabet/log4go"
	"io/ioutil"
	"strings"
)

/**

国际化文件命名：
en.toml   --  表示英文
zh.toml   --  表示中文

文件内容格式：

[[msg]]
msgctxt="msgContext"   # 消息上下文
msgid="myId"           # 消息编码
msgstr="msg content"   # 消息内容

[[msg]]
msgctxt="msgContext"   # 消息上下文
msgid="myId"           # 消息编码
msgstr="msg content"   # 消息内容

[[msg]]
msgctxt="msgContext"   # 消息上下文
msgid="myId"           # 消息编码
msgstr="msg content"   # 消息内容

[[msg]]
msgctxt="msgContext"   # 消息上下文
msgid="myId"           # 消息编码
msgstr="msg content"   # 消息内容

**/

var localeMap = make(map[string]*Locale)

/*
初始化国际化信息
*/
func Init() {
	for appName, modelsInfo := range env.Models {
		for _, model := range modelsInfo {
			i18nFolder := env.Env_Project_Resource + "/" + appName + "/" + model + "/" + env.Env_Project_Resource_Apps_I18n
			i18nFiles, err := ioutil.ReadDir(i18nFolder)
			if err == nil {
				for _, i18nFile := range i18nFiles {
					i18nFileName := i18nFile.Name()
					nameAndSuffix := strings.Split(i18nFileName, ".")
					if len(nameAndSuffix) == 2 && nameAndSuffix[1] == env.Env_Project_Resource_Apps_I18n_File_Suffix { //判断后缀名为toml
						locale := localeMap[nameAndSuffix[0]]
						if locale == nil {
							locale = &Locale{}
							locale.Init()
							localeMap[nameAndSuffix[0]] = locale
						}
						locale.Load(nameAndSuffix[0], i18nFolder+"/"+i18nFileName, model)
					}
				}
			} else {
				if log4go.IsFinestLevel() {
					log4go.FinestLog("I18N Initialization Warning : Folder[%s] is Not Found %s . err: %s . ", model, env.Env_Project_Resource_Apps_I18n, err.Error())
				}
			}
		}
	}
}

/*
获取国际化信息

@param lang  语言名称

@return (*Locale)

*/
func GetLocale(lang string) *Locale {
	locale := localeMap[lang]
	return locale
}
