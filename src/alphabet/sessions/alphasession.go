// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sessions

import (
	"alphabet/env"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var alphaSessionCounterNum int64 = 0 //计数器
var alphaSessionCounterNum_Lock = new(sync.Mutex)

var GenCookieSessionId_Lock = new(sync.Mutex)

var GenSessionGlobal_Lock = new(sync.Mutex)

func counter() int64 {
	defer func() {
		alphaSessionCounterNum_Lock.Unlock()
	}()
	alphaSessionCounterNum_Lock.Lock()
	alphaSessionCounterNum = alphaSessionCounterNum + 1
	return alphaSessionCounterNum
}

//
// 创建session信息
//    CreateAlphaSessionFactory(&AlphaOptions{ Path:   "/", MaxAge: env.Env_Web_Sessionmaxage,})
//
func CreateAlphaSessionFactory(alphaOptions *AlphaOptions) AlphaSessionFactory {
	if env.Env_Web_Session_Store == "memory" {
		var asf = new(AlphaMemorySessionFactory)
		asf.values = make(map[string]*AlphaMemorySession)
		asf.options = alphaOptions
		asf.maxObjectCount = env.Env_Web_Session_MaxObjectCount

		return asf
	} else {
		var asf = new(AlphaRedisSessionFactory)
		asf.options = alphaOptions
		return asf
	}

}

type AlphaOptions struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

type AlphaSessionFactory interface {
	GetSession(w http.ResponseWriter, r *http.Request) AlphaSession
	FlushAll()
}

type AlphaSession interface {
	Get(key string) string
	Set(key string, value string)
	Maps() map[string]string
	Clear()
}

//
// 生成sessionid的信息
//
func genSessionIdValue() string {
	timeNow := time.Now()
	numRand := 10
	bytesRand := make([]byte, numRand)

	r := rand.New(rand.NewSource(timeNow.UnixNano()))

	for i := 0; i < numRand; i = i + 1 {
		bytesRand[i] = byte(126 - r.Intn(94))
	}

	var sessionIdKey string = env.Env_Web_Sessionid + timeNow.String() + string(timeNow.UnixNano()) + string(bytesRand)
	var sha1Key = env.Env_Web_Sessiongenkey

	key := []byte(sha1Key)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(sessionIdKey))
	return fmt.Sprintf("%x", mac.Sum(nil))
}
