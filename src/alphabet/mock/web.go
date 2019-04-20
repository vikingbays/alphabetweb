// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package mock

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

/*
模拟一个HTTP请求。

@param method 设置POST或GET方式
@param urlStr 请求地址
@param param0str  请求参数，字符串形式，例如："s1=1&s2=2"
@param param1Array 请求参数，数组形式
@param header1 http请求头
@param handleServHTTP  实现一个 http.Handler.ServeHTTP 方法，并把方法名作为传参

@Return (int,string)  (httpcode:200/404/500 等，  输出的数据)
*/
func MockWebAction(method string, urlStr string,
	param0str string, param1Array map[string][]string, header1 http.Header,
	handleServHTTP func(w http.ResponseWriter, r *http.Request)) (int, string) {

	mockWebResponseWriter := NewMockWebResponseWriter()

	method = strings.ToUpper(method)

	handleServHTTP(mockWebResponseWriter, NewMockWebRequest(method, urlStr,
		param0str, param1Array, header1))
	return mockWebResponseWriter.Code, mockWebResponseWriter.Body.String()

}

/*
用于http.ResponseWriter的模拟对象测试。

@Return 返回一个MockWebResponseWriter对象，该对象可模拟http.ResponseWriter
*/
func NewMockWebResponseWriter() *MockWebResponseWriter {
	return &MockWebResponseWriter{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Code:      200,
	}
}

/*
创建一个http.Header对象
*/
func NewMockWebHeader() http.Header {
	header := make(http.Header)
	return header
}

/*
 cookie信息是纪录在http.Header 的 "Cookie" 中。
 本方法是添加cookie信息

 @Param header
 @Param cookiesMaps   cookie字符串信息
*/
func AddCookies(header http.Header, cookiesMaps map[string]string) {
	//数据格式如下：
	//Cookie:[alphabet09-session-id=364d2a29787b1eed5b386e6bf51638ad41e42d9a]
	for key, value := range cookiesMaps {
		header.Add("Cookie", fmt.Sprintf("%s=%s", key, value))
	}
}

/*
用于http.Request的模拟对象测试。

@param method 设置POST或GET方式
@param urlStr 请求地址
@param param0str  请求参数，字符串形式，例如："s1=1&s2=2"
@param param1Array 请求参数，数组形式
@param header1 http请求头

@Return 返回一个http.Request对象
*/
func NewMockWebRequest(method string, urlStr string, param0str string,
	param1Array map[string][]string, header1 http.Header) *http.Request {
	param1str := ""
	if param1Array != nil {
		for key, values := range param1Array {
			for _, value := range values {
				if param1str != "" {
					param1str = fmt.Sprintf("%s&%s=%s", param1str, key, value)
				} else {
					param1str = fmt.Sprintf("%s=%s", key, value)
				}
			}

		}
	}

	bodyStr := ""
	if param0str == "" && param1str != "" {
		bodyStr = param1str
	} else if param0str != "" && param1str == "" {
		bodyStr = param0str
	} else if param0str != "" && param1str != "" {
		bodyStr = param0str + "&" + param1str
	}

	var body io.Reader

	if bodyStr != "" {
		body = io.MultiReader(strings.NewReader(bodyStr))
	} else {
		body = nil
	}

	req, err := http.NewRequest(method, urlStr, body)

	if header1 != nil {
		req.Header = header1
	}

	if err == nil {

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		return req
	} else {
		return nil
	}

}

/*
定义MockWebResponseWriter，实现http.ResponseWriter
*/
type MockWebResponseWriter struct {
	Code        int           // the HTTP response code from WriteHeader
	HeaderMap   http.Header   // the HTTP response headers
	Body        *bytes.Buffer // if non-nil, the bytes.Buffer to append written data to
	Flushed     bool
	wroteHeader bool
}

// DefaultRemoteAddr is the default remote address to return in RemoteAddr if
// an explicit DefaultRemoteAddr isn't set on ResponseRecorder.
const DefaultRemoteAddr = "0.0.0.0"

// Header returns the response headers.
func (rw *MockWebResponseWriter) Header() http.Header {
	m := rw.HeaderMap
	if m == nil {
		m = make(http.Header)
		rw.HeaderMap = m
	}
	return m
}

// Write always succeeds and writes to rw.Body, if not nil.
func (rw *MockWebResponseWriter) Write(buf []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(200)
	}
	if rw.Body != nil {
		rw.Body.Write(buf)
	}
	return len(buf), nil
}

// WriteHeader sets rw.Code.
func (rw *MockWebResponseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.Code = code
	}
	rw.wroteHeader = true
}

// Flush sets rw.Flushed to true.
func (rw *MockWebResponseWriter) Flush() {
	if !rw.wroteHeader {
		rw.WriteHeader(200)
	}
	rw.Flushed = true
}
