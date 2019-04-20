// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sessions

import (
	"alphabet/env"
	"net/http"
	"sort"
	"time"
)

type AlphaMemorySessionFactory struct {
	values         map[string]*AlphaMemorySession // 存储的session对象
	options        *AlphaOptions                  //cookie选项
	maxObjectCount int                            // 记录最大对象数

}

func (sf *AlphaMemorySessionFactory) GetSession(w http.ResponseWriter, r *http.Request) AlphaSession {
	sessionId := env.Env_Web_Sessionid
	//sessionId := "session-id-new-test"

	sessionIdValue := sf.getCookieSessionId(r, sessionId) // 从cookie中查找是否有sessionid ？
	if sessionIdValue == "" {                             //如果在cookie中没有，就创建key
		GenCookieSessionId_Lock.Lock()
		sessionIdValue = genSessionIdValue()
		GenCookieSessionId_Lock.Unlock()
	}

	sf.updateCookieSessionId(w, sessionId, sessionIdValue) // 更新cookie信息
	GenSessionGlobal_Lock.Lock()
	defer GenSessionGlobal_Lock.Unlock()
	alphaSession := sf.values[sessionIdValue]
	if alphaSession == nil {
		alphaSession = sf.createAlphaSession(sessionIdValue) // 如果没有，就创建一个session对象
		sf.values[sessionIdValue] = alphaSession
	}
	return alphaSession
}

func (sf *AlphaMemorySessionFactory) createAlphaSession(sessionId string) *AlphaMemorySession {
	alphaSession := &AlphaMemorySession{sessionId, make(map[string]string), time.Now(), 0}
	return alphaSession
}

func (sf *AlphaMemorySessionFactory) updateCookieSessionId(w http.ResponseWriter, name string, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     sf.options.Path,
		Domain:   sf.options.Domain,
		MaxAge:   sf.options.MaxAge,
		Secure:   sf.options.Secure,
		HttpOnly: sf.options.HttpOnly,
	}
	if sf.options.MaxAge > 0 {
		d := time.Duration(sf.options.MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else if sf.options.MaxAge < 0 {
		// Set it to the past to expire now.
		cookie.Expires = time.Unix(1, 0)
	}
	http.SetCookie(w, cookie)
}

func (sf *AlphaMemorySessionFactory) getCookieSessionId(r *http.Request, sessionId string) string {
	cookie, err := r.Cookie(sessionId)
	if err != nil {
		return ""
	} else {
		return cookie.Value
	}
}

type Int64Sort []int64 // int64排序

func (i64s Int64Sort) Len() int           { return len(i64s) }
func (i64s Int64Sort) Less(i, j int) bool { return i64s[i] < i64s[j] }
func (i64s Int64Sort) Swap(i, j int)      { i64s[i], i64s[j] = i64s[j], i64s[i] }

//
// 定期清除超过有效期的session数据
//
func (sf *AlphaMemorySessionFactory) FlushAll() {
	d := time.Duration(sf.options.MaxAge) * time.Second
	currTime := time.Now()
	GenSessionGlobal_Lock.Lock()
	defer GenSessionGlobal_Lock.Unlock()
	extraLen := len(sf.values) - sf.maxObjectCount
	counterKeyMap := make(map[int64]string)
	var counterArray = make(Int64Sort, 0, 1)
	deleteLen := 0
	for key, sess := range sf.values {
		if sess.getLastVisitTime().Add(d).Before(currTime) {
			delete(sf.values, key)
			deleteLen = deleteLen + 1
		} else {
			if extraLen > 0 {
				counterKeyMap[sess.counterNum] = key
				counterArray = append(counterArray, sess.counterNum)
			}
		}
	}
	extraLen = extraLen - deleteLen // 排除删除项，再查看超出存储的Session
	if extraLen > 0 {
		sort.Sort(counterArray)
		if extraLen < len(counterArray) {
			for i := 0; i < extraLen; i++ {
				key := counterKeyMap[counterArray[i]]
				if key != "" {
					delete(sf.values, key)
				}
			}
		}
	}
}

type AlphaMemorySession struct {
	sessionId     string            // 记录sessionid信息，做为唯一标示
	values        map[string]string // 存储session信息
	lastVisitTime time.Time         // 最后一次访问时间，用于session销毁
	counterNum    int64             // 计数器，数值越大，表示越靠近最近被调用。
}

func (session *AlphaMemorySession) updateStatus() {
	session.lastVisitTime = time.Now()
	session.counterNum = counter()
}

func (session *AlphaMemorySession) getLastVisitTime() time.Time {
	return session.lastVisitTime
}

func (session *AlphaMemorySession) Get(key string) string {
	session.updateStatus()
	return session.values[key]
}

func (session *AlphaMemorySession) Set(key string, value string) {
	session.updateStatus()
	session.values[key] = value
}

func (session *AlphaMemorySession) Maps() map[string]string {
	return session.values
}

func (session *AlphaMemorySession) Clear() {
	session.values = nil
	session.values = make(map[string]string)
}
