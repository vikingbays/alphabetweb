// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package log4go

import (
	"alphabet/core/toml"
	"alphabet/env"
	"fmt"
	"os"
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

var loggerInstance Logger
var loggerLevelType Level = ERROR

func Init() {
	LINENO_CALLER_LEVEL = 3
	//替代：loggerInstance = InitDefaultLogger()
	InitLoggerInstance(InitDefaultLogger())
}

func InitLoggerInstance(logger Logger) {
	loggerInstance = logger
}

/**
 * 判断当前日志级别是否为：FINEST
 *
 */
func IsFinestLevel() bool {
	if loggerLevelType == FINEST {
		return true
	} else {
		return false
	}
}

func IsFineLevel() bool {
	if loggerLevelType == FINE {
		return true
	} else {
		return false
	}
}

func IsDebugLevel() bool {
	if loggerLevelType == DEBUG {
		return true
	} else {
		return false
	}
}

func IsTraceLevel() bool {
	if loggerLevelType == TRACE {
		return true
	} else {
		return false
	}
}

func IsInfoLevel() bool {
	if loggerLevelType == INFO {
		return true
	} else {
		return false
	}
}

func IsWarningLevel() bool {
	if loggerLevelType == WARNING {
		return true
	} else {
		return false
	}
}

func IsErrorLevel() bool {
	if loggerLevelType == ERROR {
		return true
	} else {
		return false
	}
}

/*
记录Finest级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func FinestLog(arg0 interface{}, args ...interface{}) {
	if loggerInstance == nil {
		return
	}
	loggerInstance.Finest(arg0, args...)
}

/*
记录Fine级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func FineLog(arg0 interface{}, args ...interface{}) {
	if loggerInstance == nil {
		return
	}
	loggerInstance.Fine(arg0, args...)
}

/*
记录Debug级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func DebugLog(arg0 interface{}, args ...interface{}) {
	if loggerInstance == nil {
		return
	}
	loggerInstance.Debug(-1, arg0, args...)
}

/*
记录Trace级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func TraceLog(arg0 interface{}, args ...interface{}) {
	if loggerInstance == nil {
		return
	}
	loggerInstance.Trace(-1, arg0, args...)
}

/*
记录Info级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func InfoLog(arg0 interface{}, args ...interface{}) {
	if loggerInstance == nil {
		return
	}
	loggerInstance.Info(-1, arg0, args...)
}

/*
记录Warning级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func WarnLog(arg0 interface{}, args ...interface{}) {
	if loggerInstance == nil {
		return
	}
	loggerInstance.Warn(-1, arg0, args...)
}

/*
记录Error级别日志。
记录的方式可以类似于 fmt.Printf 或 fmt.Println

@param arg0

@param args ...

*/
func ErrorLog(arg0 interface{}, args ...interface{}) {
	if loggerInstance == nil {
		return
	}
	loggerInstance.Error(-1, arg0, args...)
}

/**
* 判断日志类是否初始化成功。
 */
func IsInitCompleted() bool {
	if loggerInstance == nil {
		return false
	} else {
		return true
	}
}

func GetL() Logger {
	return loggerInstance
}

/*
初始化日志组件
*/
func InitDefaultLogger() Logger {
	appName := env.Env_Project_Resource_Apps[0] // 选择第一个作为日志初始化使用
	logconfigFile := env.Env_Project_Resource + "/" + appName + "/" + env.GetSysconfigPath() + env.Env_Project_Resource_LogConfig

	envProjectRoot := env.Env_Project_Root
	return InitLogger(logconfigFile, envProjectRoot)

}

func InitLogger(pathconfig string, envProjectRoot string) Logger {
	logger := make(Logger) //用于替换：logger := NewLogger()
	var loggerConfigToml Logger_Config_Toml
	if _, err := toml.DecodeFile(pathconfig, &loggerConfigToml); err != nil {
		fmt.Println(err)
		return nil
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

func SetLoggerFilterConfig(filterConfigToml *Filter_Config_Toml, logger *Logger) {
	if filterConfigToml.Enabled == "true" {

		rotatable := (filterConfigToml.Rotate == "true")
		var currLoggerLevelType Level
		if filterConfigToml.Level == "FINEST" {
			currLoggerLevelType = FINEST
		} else if filterConfigToml.Level == "FINE" {
			currLoggerLevelType = FINE
		} else if filterConfigToml.Level == "DEBUG" {
			currLoggerLevelType = DEBUG
		} else if filterConfigToml.Level == "TRACE" {
			currLoggerLevelType = TRACE
		} else if filterConfigToml.Level == "INFO" {
			currLoggerLevelType = INFO
		} else if filterConfigToml.Level == "WARNING" {
			currLoggerLevelType = WARNING
		} else if filterConfigToml.Level == "ERROR" {
			currLoggerLevelType = ERROR
		}

		if loggerLevelType > currLoggerLevelType {
			loggerLevelType = currLoggerLevelType
		}

		if filterConfigToml.Type == "file" {

			fileLogWriter := NewFileLogWriter(filterConfigToml.Filename, rotatable)
			fileLogWriter.SetFormat(filterConfigToml.Format)
			fileLogWriter.SetRotate(rotatable)

			fileLogWriter.SetRotateSize(StrToNumSuffix(filterConfigToml.RotateMaxSize, 1024))
			fileLogWriter.SetRotateLines(StrToNumSuffix(filterConfigToml.RotateMaxLines, 1000))
			fileLogWriter.SetRotateDaily((filterConfigToml.RotateDaily == "true"))
			logger.AddFilter(filterConfigToml.Tag, currLoggerLevelType, fileLogWriter)
		} else if filterConfigToml.Type == "console" {
			consoleLogWriter := NewConsoleLogWriter()
			consoleLogWriter.SetFormat(filterConfigToml.Format)
			logger.AddFilter(filterConfigToml.Tag, currLoggerLevelType, consoleLogWriter)

		} else if filterConfigToml.Type == "socket" {
			socketLogWriter := NewSocketLogWriter(filterConfigToml.Protocol, filterConfigToml.Endpoint)
			//socketLogWriter.SetFormat(filterConfigToml.Format)
			logger.AddFilter(filterConfigToml.Tag, currLoggerLevelType,
				socketLogWriter)
		}

	}

}

func initTestLog() {
	projectPath := os.Getenv("GOLANG_PROJECT_TESTUNIT")
	env.InitWithoutLoad(projectPath)

	logconfigFile := env.Env_Project_Resource + "/alphabet/core/abtest/logconfig_test.toml"
	LINENO_CALLER_LEVEL = 3
	InitLoggerInstance(InitLogger(logconfigFile, env.Env_Project_Root))
}
