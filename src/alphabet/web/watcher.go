// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/core/fsnotify"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"sync"
	"time"
)

//监控文件变化情况的对象类
type MonitorFileEventWatcher struct {
	Watcher    *fsnotify.Watcher //监控对象
	StoreEvent map[string]string //存储需要监控的文件或文件夹信息，w.StoreEvent[filepath] = filepath
	Lock       *sync.Mutex       //全局锁
}

/*
创建监控文件变化的对象
*/
func NewMonitorFileEventWatcher() *MonitorFileEventWatcher {
	monitorFileEventWatcher := new(MonitorFileEventWatcher)
	monitorFileEventWatcher.StoreEvent = make(map[string]string)
	monitorFileEventWatcher.Lock = new(sync.Mutex)
	watcher0, err := fsnotify.NewWatcher()
	if err != nil {
		log4go.ErrorLog(message.ERR_WEB0_39060, err.Error())
		monitorFileEventWatcher.Watcher = nil
	} else {
		monitorFileEventWatcher.Watcher = watcher0

		go func() {
			for {
				select {
				case ev := <-monitorFileEventWatcher.Watcher.Event:
					monitorFileEventWatcher.putEvent(ev.Name)
				case err := <-monitorFileEventWatcher.Watcher.Error:
					log4go.ErrorLog(message.ERR_WEB0_39061, err)
				}
			}
		}()
	}
	return monitorFileEventWatcher
}

/*
 添加需要跟踪的文件或者文件夹，如果是文件夹只能跟踪下一级文件，子目录无法跟踪。
 所有添加的文件或者文件夹下的文件发生变化后，会触发事件（具体调用DoEvent注册的事件）。

*/
func (w *MonitorFileEventWatcher) putEvent(filepath string) {
	w.StoreEvent[filepath] = filepath
}

/*
当跟踪到有文件发现改变时，会对该文件变化情况进行处理
本次事件处理是异步方式，每1.5秒扫描一次。
*/
func (w *MonitorFileEventWatcher) DoEvent(f func(string)) {
	for true {
		time.Sleep(1500 * time.Millisecond) //1.5秒处理一次
		store := make([]string, 10, 10)
		w.Lock.Lock()
		if len(w.StoreEvent) > 0 {
			for _, v := range w.StoreEvent {
				store = append(store, v)
			}
		}
		w.StoreEvent = nil
		w.StoreEvent = make(map[string]string)
		w.Lock.Unlock()

		for _, v := range store {
			f(v)
		}

	}

}

/*
添加需要监控的文件夹信息，只能监控该文件夹下文件变化情况，子文件夹无法监控。

@param  folder  文件夹路径
*/
func (w *MonitorFileEventWatcher) AddMonitorEventFileWatch(folder string) {
	err := w.Watcher.Watch(folder)
	if err != nil {
		log4go.ErrorLog(message.ERR_WEB0_39062, folder, err.Error())
	}
}
