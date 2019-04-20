// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package log4gohttp

import (
	"alphabet/core/toml"
	"alphabet/env"
	"alphabet/log4go"
	"fmt"
	"strings"
)

/**
 * 日志级别： FINEST <  FINE  <  DEBUG  <  TRACE  <  INFO  <  WARNING  <  ERROR
 *
 */

type Logger_Config_Toml struct {
	Filters []Filter_Config_Toml
}

type Filter_Config_Toml struct {
	Enabled string
	Tag     string
	Type    string
	Level   string
	Format  string

	Filename       string
	Rotate         string
	RotateMaxSize  string
	RotateMaxLines string
	RotateDaily    string

	Endpoint string
	Protocol string
}

var loggerInstance log4go.Logger
var loggerLevelType log4go.Level = log4go.ERROR

func Init() {
	log4go.LINENO_CALLER_LEVEL = 3
	loggerInstance = InitDefaultLogger()
}

/*
记录Finest级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func FinestLog(arg0 interface{}, args ...interface{}) {
	loggerInstance.Finest(arg0, args...)
}

/*
记录Fine级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func FineLog(arg0 interface{}, args ...interface{}) {
	loggerInstance.Fine(arg0, args...)
}

/*
记录Debug级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func DebugLog(arg0 interface{}, args ...interface{}) {
	loggerInstance.Debug(-1, arg0, args...)
}

/*
记录Trace级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func TraceLog(arg0 interface{}, args ...interface{}) {
	loggerInstance.Trace(-1, arg0, args...)
}

/*
记录Info级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func InfoLog(arg0 interface{}, args ...interface{}) {
	loggerInstance.Info(-1, arg0, args...)
}

/*
记录Warning级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func WarnLog(arg0 interface{}, args ...interface{}) {
	loggerInstance.Warn(-1, arg0, args...)
}

/*
记录Error级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func ErrorLog(arg0 interface{}, args ...interface{}) {
	loggerInstance.Error(-1, arg0, args...)
}

/*
初始化日志组件
*/
func InitDefaultLogger() log4go.Logger {
	appName := env.Env_Project_Resource_Apps[0] // 选择第一个作为日志初始化使用
	logconfigFile := env.Env_Project_Resource + "/" + appName + "/" + env.GetSysconfigPath() + env.Env_Project_Resource_LogHttpConfig

	envProjectRoot := env.Env_Project_Root
	return InitLogger(logconfigFile, envProjectRoot)

}

func InitLogger(pathconfig string, envProjectRoot string) log4go.Logger {
	logger := make(log4go.Logger) //用于替换：logger := NewLogger()
	var loggerConfigToml Logger_Config_Toml
	if _, err := toml.DecodeFile(pathconfig, &loggerConfigToml); err != nil {
		fmt.Println(err)
	} else {
		for _, filter := range loggerConfigToml.Filters {
			SetInitVariable(&filter, envProjectRoot)
			SetLoggerFilterConfig(&filter, &logger)
		}
	}
	return logger
}

func SetInitVariable(filterConfigToml *Filter_Config_Toml, envProjectRoot string) {
	if filterConfigToml.Type == "file" {
		filterConfigToml.Filename = strings.Replace(filterConfigToml.Filename, "${project}", envProjectRoot, -1)
	}
}

func SetLoggerFilterConfig(filterConfigToml *Filter_Config_Toml, logger *log4go.Logger) {
	if filterConfigToml.Enabled == "true" {

		rotatable := (filterConfigToml.Rotate == "true")
		var currLoggerLevelType log4go.Level
		if filterConfigToml.Level == "FINEST" {
			currLoggerLevelType = log4go.FINEST
		} else if filterConfigToml.Level == "FINE" {
			currLoggerLevelType = log4go.FINE
		} else if filterConfigToml.Level == "DEBUG" {
			currLoggerLevelType = log4go.DEBUG
		} else if filterConfigToml.Level == "TRACE" {
			currLoggerLevelType = log4go.TRACE
		} else if filterConfigToml.Level == "INFO" {
			currLoggerLevelType = log4go.INFO
		} else if filterConfigToml.Level == "WARNING" {
			currLoggerLevelType = log4go.WARNING
		} else if filterConfigToml.Level == "ERROR" {
			currLoggerLevelType = log4go.ERROR
		}

		if loggerLevelType > currLoggerLevelType {
			loggerLevelType = currLoggerLevelType
		}

		if filterConfigToml.Type == "file" {

			fileLogWriter := log4go.NewFileLogWriter(filterConfigToml.Filename, rotatable)
			fileLogWriter.SetFormat(filterConfigToml.Format)
			fileLogWriter.SetRotate(rotatable)

			fileLogWriter.SetRotateSize(log4go.StrToNumSuffix(filterConfigToml.RotateMaxSize, 1024))
			fileLogWriter.SetRotateLines(log4go.StrToNumSuffix(filterConfigToml.RotateMaxLines, 1000))
			fileLogWriter.SetRotateDaily((filterConfigToml.RotateDaily == "true"))
			logger.AddFilter(filterConfigToml.Tag, currLoggerLevelType, fileLogWriter)
		} else if filterConfigToml.Type == "console" {
			logger.AddFilter(filterConfigToml.Tag, currLoggerLevelType, log4go.NewConsoleLogWriter())
		} else if filterConfigToml.Type == "socket" {
			logger.AddFilter(filterConfigToml.Tag, currLoggerLevelType,
				log4go.NewSocketLogWriter(filterConfigToml.Protocol, filterConfigToml.Endpoint))
		}

	}

}
