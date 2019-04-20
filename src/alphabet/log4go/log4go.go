// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

// [第三方包]实现日志记录
//
// 在当前的alphabet webapp框架中，重新封装了log4go日志。
//
// 日志配置在logconfig.toml中定义。可以配置多个输出。
// 其中，logconfig.toml的定义结构可参考：
//  ##日志输出格式:
//  ##      %T - Time (15:04:05 MST)
//  ##      %t - Time (15:04)
//  ##      %D - Date (2006/01/02)
//  ##      %d - Date (01/02/06)
//  ##      %L - Level (FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR)
//  ##      %S - Source
//  ##      %G - goroutine ID  协程号
//  ##      %U - Unique serial number 请求唯一序列号 ，在web场景下，需要判断header是否有属性 env.Env_Web_Header_Unique_Serial_Number
//  ##      %M - Message
//  ##      It ignores unknown format strings (and removes them)
//  ##      Recommended: "[%D %T] [%L] (%S) %M"
//
//  [[filters]]                          ## 定义一个日志输出格式，可以定义多个。
//  enabled="true"                       ## 设置该日志是否启用
//  tag="stdout"                         ## 设置输出方式，采用控制台输出
//  type="console"                       ## 采用控制台输出
//  level="DEBUG"                        ## 级别定义 (FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR) ，FINEST最低。
//  format="[%D %T] [%L] (%S) %M "       ## 输出格式
//
//  [[filters]]
//  enabled="true"
//  tag="file"                           ## 文件方式
//  type="file"                          ## 文件方式
//  level="INFO"
//  format="[%D %T] [%L] (%S) %M "
//  filename="${project}/logs/apps.log"  ## 可以带变量 ${project}/logs/test.log
//  rotate="false"                       ## 是否采用循环输出
//  rotateMaxSize="0"                    ## \d+[KMG]? Suffixes are in terms of 2**10
//  rotateMaxLines="0"                   ## \d+[KMG]? Suffixes are in terms of thousands
//  rotateDaily="false"
//
// 日志使用非常简单：（传参方式类似于 fmt.Printf 或 fmt.Println）
//  log4go.FinestLog(arg0 interface{}, args ...interface{})    ## 记录Finest级别日志，最低
//
//  log4go.FineLog(arg0 interface{}, args ...interface{})      ## 记录Fine级别日志
//
//  log4go.DebugLog(arg0 interface{}, args ...interface{})     ## 记录Debug级别日志
//
//  log4go.TraceLog(arg0 interface{}, args ...interface{})     ## 记录Trace级别日志
//
//  log4go.InfoLog(arg0 interface{}, args ...interface{})      ## 记录Info级别日志
//
//  log4go.WarnLog(arg0 interface{}, args ...interface{})      ## 记录Warn级别日志
//
//  log4go.ErrorLog(arg0 interface{}, args ...interface{})     ## 记录Error级别日志，最高
//
//
//
//
// Package log4go provides level-based and highly configurable logging.
//
// Enhanced Logging
//
// This is inspired by the logging functionality in Java.  Essentially, you create a Logger
// object and create output filters for it.  You can send whatever you want to the Logger,
// and it will filter that based on your settings and send it to the outputs.  This way, you
// can put as much debug code in your program as you want, and when you're done you can filter
// out the mundane messages so only the important ones show up.
//
// Utility functions are provided to make life easier. Here is some example code to get started:
//
// log := log4go.NewLogger()
// log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())
// log.AddFilter("log",    log4go.FINE,  log4go.NewFileLogWriter("example.log", true))
// log.Info("The time is now: %s", time.LocalTime().Format("15:04:05 MST 2006/01/02"))
//
// The first two lines can be combined with the utility NewDefaultLogger:
//
// log := log4go.NewDefaultLogger(log4go.DEBUG)
// log.AddFilter("log",    log4go.FINE,  log4go.NewFileLogWriter("example.log", true))
// log.Info("The time is now: %s", time.LocalTime().Format("15:04:05 MST 2006/01/02"))
//
// Usage notes:
// - The ConsoleLogWriter does not display the source of the message to standard
//   output, but the FileLogWriter does.
// - The utility functions (Info, Debug, Warn, etc) derive their source from the
//   calling function, and this incurs extra overhead.
//
// Changes from 2.0:
// - The external interface has remained mostly stable, but a lot of the
//   internals have been changed, so if you depended on any of this or created
//   your own LogWriter, then you will probably have to update your code.  In
//   particular, Logger is now a map and ConsoleLogWriter is now a channel
//   behind-the-scenes, and the LogWrite method no longer has return values.
//
// Future work: (please let me know if you think I should work on any of these particularly)
// - Log file rotation
// - Logging configuration files ala log4j
// - Have the ability to remove filters?
// - Have GetInfoChannel, GetDebugChannel, etc return a chan string that allows
//   for another method of logging
// - Add an XML filter type
package log4go

import (
	"alphabet/env"
	"alphabet/log4go/message"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// Version information
const (
	L4G_VERSION = "log4go-v3.0.1"
	L4G_MAJOR   = 3
	L4G_MINOR   = 0
	L4G_BUILD   = 1
)

/****** Constants ******/

// These are the integer logging levels used by the logger
type Level int

const (
	FINEST Level = iota
	FINE
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
)

var LINENO_CALLER_LEVEL int = 2

// Logging level strings
var (
	levelStrings = [...]string{"FNST", "FINE", "DEBG", "TRAC", "INFO", "WARN", "EROR", "CRIT"}
)

func (l Level) String() string {
	if l < 0 || int(l) > len(levelStrings) {
		return "UNKNOWN"
	}
	return levelStrings[int(l)]
}

/****** Variables ******/
var (
	// LogBufferLength specifies how many log messages a particular log4go
	// logger can buffer at a time before writing them.
	LogBufferLength = 32
)

/****** LogRecord ******/

// A LogRecord contains all of the pertinent information for each message
type LogRecord struct {
	Level     Level     // The log level
	Created   time.Time // The time at which the log message was created (nanoseconds)
	Source    string    // The message source
	Message   string    // The log message
	CallerObj *Caller   // 调用链信息
}

/****** LogWriter ******/

// This is an interface for anything that should be able to write logs
type LogWriter interface {
	// This will be called to log a LogRecord message.
	LogWrite(rec *LogRecord)

	// This should clean up anything lingering about the LogWriter, as it is called before
	// the LogWriter is removed.  LogWrite should not be called after Close.
	Close()
}

/****** Logger ******/

// A Filter represents the log level below which no log records are written to
// the associated LogWriter.
type Filter struct {
	Level Level
	LogWriter
}

// A Logger represents a collection of Filters through which log messages are
// written.
type Logger map[string]*Filter

// Create a new logger.
//
// DEPRECATED: Use make(Logger) instead.
func NewLogger() Logger {
	os.Stderr.WriteString("warning: use of deprecated NewLogger\n")
	return make(Logger)
}

// Create a new logger with a "stdout" filter configured to send log messages at
// or above lvl to standard output.
//
// DEPRECATED: use NewDefaultLogger instead.
func NewConsoleLogger(lvl Level) Logger {
	os.Stderr.WriteString("warning: use of deprecated NewConsoleLogger\n")
	return Logger{
		"stdout": &Filter{lvl, NewConsoleLogWriter()},
	}
}

// Create a new logger with a "stdout" filter configured to send log messages at
// or above lvl to standard output.
func NewDefaultLogger(lvl Level) Logger {
	return Logger{
		"stdout": &Filter{lvl, NewConsoleLogWriter()},
	}
}

// Closes all log writers in preparation for exiting the program or a
// reconfiguration of logging.  Calling this is not really imperative, unless
// you want to guarantee that all log messages are written.  Close removes
// all filters (and thus all LogWriters) from the logger.
func (log Logger) Close() {
	// Close all open loggers
	for name, filt := range log {
		filt.Close()
		delete(log, name)
	}
}

// Add a new LogWriter to the Logger which will only log messages at lvl or
// higher.  This function should not be called from multiple goroutines.
// Returns the logger for chaining.
func (log Logger) AddFilter(name string, lvl Level, writer LogWriter) Logger {
	log[name] = &Filter{lvl, writer}
	return log
}

/******* Logging *******/
// Send a formatted log message internally
// linenoLevel 如果是 -1 ，表示不设置，使用缺省
func (log Logger) intLogf(linenoLevel int, lvl Level, format string, args ...interface{}) {
	if linenoLevel == -1 {
		linenoLevel = LINENO_CALLER_LEVEL
	}
	skip := true
	// Determine if any logging will be done
	for _, filt := range log {
		if lvl >= filt.Level {
			skip = false
			break
		}
	}
	if skip {
		return
	}
	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(linenoLevel)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	// Make the log record

	var rec *LogRecord

	if env.Switch_CallChain {
		rec = &LogRecord{
			Level:     lvl,
			Created:   time.Now(),
			Source:    src,
			Message:   msg,
			CallerObj: GetCallChain().GetCurrentCaller(),
		}
	} else {
		rec = &LogRecord{
			Level:   lvl,
			Created: time.Now(),
			Source:  src,
			Message: msg,
		}
	}

	// Dispatch the logs
	for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
}

// Send a closure log message internally
func (log Logger) intLogc(lvl Level, closure func() string) {
	skip := true

	// Determine if any logging will be done
	for _, filt := range log {
		if lvl >= filt.Level {
			skip = false
			break
		}
	}
	if skip {
		return
	}

	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	// Make the log record

	var rec *LogRecord

	if env.Switch_CallChain {
		rec = &LogRecord{
			Level:     lvl,
			Created:   time.Now(),
			Source:    src,
			Message:   closure(),
			CallerObj: GetCallChain().GetCurrentCaller(),
		}
	} else {
		rec = &LogRecord{
			Level:   lvl,
			Created: time.Now(),
			Source:  src,
			Message: closure(),
		}
	}

	// Dispatch the logs
	for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
}

// Send a log message with manual level, source, and message.
func (log Logger) Log(lvl Level, source, message string) {
	skip := true
	// Determine if any logging will be done
	for _, filt := range log {
		if lvl >= filt.Level {
			skip = false
			break
		}
	}
	if skip {
		return
	}

	// Make the log record

	var rec *LogRecord

	if env.Switch_CallChain {
		rec = &LogRecord{
			Level:     lvl,
			Created:   time.Now(),
			Source:    source,
			Message:   message,
			CallerObj: GetCallChain().GetCurrentCaller(),
		}
	} else {
		rec = &LogRecord{
			Level:   lvl,
			Created: time.Now(),
			Source:  source,
			Message: message,
		}
	}

	// Dispatch the logs
	for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
}

// Logf logs a formatted log message at the given log level, using the caller as
// its source.
func (log Logger) Logf(lvl Level, format string, args ...interface{}) {
	log.intLogf(-1, lvl, format, args...)
}

// Logc logs a string returned by the closure at the given log level, using the caller as
// its source.  If no log message would be written, the closure is never called.
func (log Logger) Logc(lvl Level, closure func() string) {
	log.intLogc(lvl, closure)
}

// Finest logs a message at the finest log level.
// See Debug for an explanation of the arguments.
func (log Logger) Finest(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINEST
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(-1, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(lvl, first)
	case message.MessageType:
		log.intLogf(-1, lvl, first.String(), args...)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(-1, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Fine logs a message at the fine log level.
// See Debug for an explanation of the arguments.
func (log Logger) Fine(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(-1, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(lvl, first)
	case message.MessageType:
		log.intLogf(-1, lvl, first.String(), args...)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(-1, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Debug is a utility method for debug log messages.
// The behavior of Debug depends on the first argument:
// - arg0 is a string
//   When given a string as the first argument, this behaves like Logf but with
//   the DEBUG log level: the first argument is interpreted as a format for the
//   latter arguments.
// - arg0 is a func()string
//   When given a closure of type func()string, this logs the string returned by
//   the closure iff it will be logged.  The closure runs at most one time.
// - arg0 is interface{}
//   When given anything else, the log message will be each of the arguments
//   formatted with %v and separated by spaces (ala Sprint).
func (log Logger) Debug(linenoLevel int, arg0 interface{}, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(linenoLevel, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(lvl, first)
	case message.MessageType:
		log.intLogf(linenoLevel, lvl, first.String(), args...)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(linenoLevel, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Trace logs a message at the trace log level.
// See Debug for an explanation of the arguments.
func (log Logger) Trace(linenoLevel int, arg0 interface{}, args ...interface{}) {
	const (
		lvl = TRACE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(linenoLevel, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(lvl, first)
	case message.MessageType:
		log.intLogf(linenoLevel, lvl, first.String(), args...)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(linenoLevel, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Info logs a message at the info log level.
// See Debug for an explanation of the arguments.
func (log Logger) Info(linenoLevel int, arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)

	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(linenoLevel, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		log.intLogc(lvl, first)
	case message.MessageType:
		log.intLogf(linenoLevel, lvl, first.String(), args...)
	default:
		// Build a format string so that it will be similar to Sprint
		fmt.Println(arg0)
		log.intLogf(linenoLevel, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Warn logs a message at the warning log level and returns the formatted error.
// At the warning level and higher, there is no performance benefit if the
// message is not actually logged, because all formats are processed and all
// closures are executed to format the error message.
// See Debug for further explanation of the arguments.
func (log Logger) Warn(linenoLevel int, arg0 interface{}, args ...interface{}) error {
	const (
		lvl = WARNING
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		msg = first()
	case message.MessageType:
		msg = fmt.Sprintf(first.String(), args...)
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	log.intLogf(linenoLevel, lvl, msg)
	return errors.New(msg)
}

// Error logs a message at the error log level and returns the formatted error,
// See Warn for an explanation of the performance and Debug for an explanation
// of the parameters.
func (log Logger) Error(linenoLevel int, arg0 interface{}, args ...interface{}) error {
	const (
		lvl = ERROR
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		msg = first()
	case message.MessageType:
		msg = fmt.Sprintf(first.String(), args...)
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	log.intLogf(linenoLevel, lvl, msg)
	return errors.New(msg)
}

// Critical logs a message at the critical log level and returns the formatted error,
// See Warn for an explanation of the performance and Debug for an explanation
// of the parameters.
func (log Logger) Critical(linenoLevel int, arg0 interface{}, args ...interface{}) error {
	const (
		lvl = CRITICAL
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		msg = first()
	case message.MessageType:
		msg = fmt.Sprintf(first.String(), args...)
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	log.intLogf(linenoLevel, lvl, msg)
	return errors.New(msg)
}
