// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sessions

import (
	"alphabet/cache"
	"alphabet/core/redis"
	"alphabet/env"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"net/http"
	"time"
)

var sessionIdPrefix = "sessionid_"

type AlphaRedisSessionFactory struct {
	options *AlphaOptions //cookie选项
}

func (sf *AlphaRedisSessionFactory) GetSession(w http.ResponseWriter, r *http.Request) AlphaSession {
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

	alphaSession := sf.createAlphaSession(sessionIdValue)

	return alphaSession
}

func (sf *AlphaRedisSessionFactory) createAlphaSession(sessionId string) *AlphaRedisSession {

	cacheFinder, err := cache.GetCacheFinder(env.Env_Web_Session_Store_Name)
	defer cacheFinder.Close()
	if err != nil {
		log4go.ErrorLog(message.ERR_SESS_39005, err)
	} else {
		sessionKeys, _ := redis.Values(cacheFinder.GetCache().Do("keys", sessionIdPrefix+sessionId))

		sessionKey := ""
		if len(sessionKeys) == 1 {
			sessionKey = string(sessionKeys[0].([]byte))
		}

		if sessionKey == "" {

			cacheFinder.SetMap(sessionIdPrefix+sessionId, sessionIdPrefix+sessionId, "1")
			cacheFinder.Expire(sessionIdPrefix+sessionId, sf.options.MaxAge)
		}
	}

	return &AlphaRedisSession{sessionId, sf.options.MaxAge}
}

func (sf *AlphaRedisSessionFactory) updateCookieSessionId(w http.ResponseWriter, name string, value string) {
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

func (sf *AlphaRedisSessionFactory) getCookieSessionId(r *http.Request, sessionId string) string {
	cookie, err := r.Cookie(sessionId)
	if err != nil {
		return ""
	} else {
		return cookie.Value
	}
}

//
// 定期清除超过有效期的session数据
//
func (sf *AlphaRedisSessionFactory) FlushAll() {

}

type AlphaRedisSession struct {
	sessionId string // 记录sessionid信息，做为唯一标示
	maxAge    int    // 有效期，单位秒
}

func (session *AlphaRedisSession) Get(key string) string {
	cacheFinder, err := cache.GetCacheFinder(env.Env_Web_Session_Store_Name)
	defer cacheFinder.Close()
	if err != nil {
		log4go.ErrorLog(message.ERR_SESS_39006, err)
	} else {
		value := cacheFinder.GetMapField(sessionIdPrefix+session.sessionId, key)
		cacheFinder.Expire(sessionIdPrefix+session.sessionId, session.maxAge)
		return value
	}
	return ""
}

func (session *AlphaRedisSession) Set(key string, value string) {
	cacheFinder, err := cache.GetCacheFinder(env.Env_Web_Session_Store_Name)
	defer cacheFinder.Close()
	if err != nil {
		log4go.ErrorLog(message.ERR_SESS_39007, err)
	} else {
		cacheFinder.SetMap(sessionIdPrefix+session.sessionId, key, value)
		cacheFinder.Expire(sessionIdPrefix+session.sessionId, session.maxAge)
	}
}

func (session *AlphaRedisSession) Maps() map[string]string {
	cacheFinder, err := cache.GetCacheFinder(env.Env_Web_Session_Store_Name)
	defer cacheFinder.Close()
	if err != nil {
		log4go.ErrorLog(message.ERR_SESS_39008, err)
	} else {
		return cacheFinder.GetMap(sessionIdPrefix + session.sessionId)
	}
	return make(map[string]string)
}

func (session *AlphaRedisSession) Clear() {
	cacheFinder, err := cache.GetCacheFinder(env.Env_Web_Session_Store_Name)
	defer cacheFinder.Close()
	if err != nil {
		log4go.ErrorLog(message.ERR_SESS_39009, err)
		//log4go.ErrorLog(err)
	} else {
		cacheFinder.DelMap(sessionIdPrefix + session.sessionId)
	}
}
