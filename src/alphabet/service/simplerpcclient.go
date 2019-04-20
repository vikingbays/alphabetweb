// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package service

import (
	"alphabet/core/utils"
	"alphabet/env"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/http2"
)

type SimpleRpcClient struct {
	Protocol    string
	Ip          string
	Port        int
	client      *http.Client
	conn        net.Conn
	header      map[string]string
	cookies     []*http.Cookie
	respCookies []*http.Cookie
	prefix      string
}

/*
清空header
*/
func (sc *SimpleRpcClient) ClearHeader() {
	sc.header = nil
}

/*
添加header
*/
func (sc *SimpleRpcClient) AddHeader(k, v string) {
	if sc.header == nil {
		sc.header = make(map[string]string)
	}
	sc.header[k] = v
}

/*
清空cookie信息
*/
func (sc *SimpleRpcClient) ClearCookies() {
	sc.cookies = nil
}

/*
通过kv方式添加cookie信息
*/
func (sc *SimpleRpcClient) AddCookieFromKV(k, v string) {
	if sc.cookies == nil {
		sc.cookies = make([]*http.Cookie, 0, 1)
	}
	cookie := &http.Cookie{
		Name:  k,
		Value: v,
	}
	sc.cookies = append(sc.cookies, cookie)
}

/*
通过http.Cookie对象添加cookie信息
*/
func (sc *SimpleRpcClient) AddCookieFromObject(cookie *http.Cookie) {
	if sc.cookies == nil {
		sc.cookies = make([]*http.Cookie, 0, 1)
	}
	sc.cookies = append(sc.cookies, cookie)
}

func (sc *SimpleRpcClient) GetRespCookies() []*http.Cookie {
	return sc.respCookies
}

/*
创建客户端连接。
@param timeout 设置网络超时时间，如果 timeout<=0，表示不设置超时
*/
func (sc *SimpleRpcClient) Connect() (err error) {

	log4go.DebugLog(message.DEG_MSSV_79005, sc.Protocol, sc.Ip, sc.Port)

	if sc.Protocol == "rpc_tcp" {
		sc.prefix = "http://tcp%s"
		tr := &http.Transport{
			DialContext: func(ctx context.Context, proto, addr string) (conn net.Conn, err error) {
				var d *net.Dialer = &net.Dialer{}
				conn, err = d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", sc.Ip, sc.Port))
				if err != nil {
					log4go.ErrorLog(message.ERR_MSSV_39034, "tcp", fmt.Sprintf("%s:%d", sc.Ip, sc.Port), err.Error())
				}
				if sc.conn != nil {
					sc.conn.Close()
					sc.conn = nil
				}
				sc.conn = conn
				return
			},
		}
		sc.client = &http.Client{Transport: tr}
	} else if sc.Protocol == "rpc_tcp_ssl" {
		sc.prefix = "https://tcp%s"
		/*
			tr = &http.Transport{
				DialContext: func(ctx context.Context, proto, addr string) (conn net.Conn, err error) {
					var d *net.Dialer = &net.Dialer{}
					conn, err = d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", sc.Ip, sc.Port))
					if err != nil {
						log4go.ErrorLog(message.ERR_MSSV_39034, "tcp", fmt.Sprintf("%s:%d", sc.Ip, sc.Port), err.Error())
					}
					if sc.conn != nil {
						sc.conn.Close()
						sc.conn = nil
					}
					sc.conn = conn
					return
				},
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		*/

		tr := &http2.Transport{
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				addr = fmt.Sprintf("%s:%d", sc.Ip, sc.Port)
				conn, err := tls.Dial(network, addr, cfg)
				if err != nil {
					return nil, err
				}
				if err := conn.Handshake(); err != nil {
					return nil, err
				}
				if !cfg.InsecureSkipVerify {
					if err := conn.VerifyHostname(cfg.ServerName); err != nil {
						return nil, err
					}
				}
				state := conn.ConnectionState()
				if p := state.NegotiatedProtocol; p != http2.NextProtoTLS {
					return nil, fmt.Errorf("http2: unexpected ALPN protocol %q; want %q", p, http2.NextProtoTLS)
				}
				if !state.NegotiatedProtocolIsMutual {
					return nil, errors.New("http2: could not negotiate protocol mutually")
				}
				if sc.conn != nil {
					sc.conn.Close()
					sc.conn = nil
				}
				sc.conn = conn
				return conn, nil
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true}, //跳过验证服务端证书
		}
		sc.client = &http.Client{Transport: tr}

	} else if sc.Protocol == "rpc_unix" {
		//sc.prefix = "http://tcp%s"
		sc.prefix = "http://unix%s"
		tr := &http.Transport{
			DialContext: func(ctx context.Context, proto, addr string) (conn net.Conn, err error) {
				var d *net.Dialer = &net.Dialer{}
				conn, err = d.DialContext(ctx, "unix", sc.Ip)
				if err != nil {
					log4go.ErrorLog(message.ERR_MSSV_39034, "unix", sc.Ip, err.Error())
				}
				if sc.conn != nil {
					sc.conn.Close()
					sc.conn = nil
				}
				sc.conn = conn
				return
			},
		}
		sc.client = &http.Client{Transport: tr}
	} else {
		err = errors.New(fmt.Sprintf("protocol = %s , protocol is error . ", sc.Protocol))
	}
	return

}

/*
 * 关闭客户端连接。
 */
func (sc *SimpleRpcClient) Ping() bool {
	if sc.conn == nil {
		log4go.ErrorLog(message.ERR_MSSV_39035)
		return false
	}
	laddr := sc.conn.LocalAddr()
	if laddr == nil {
		return false
	} else {
		return true
	}
}

/*
 * 关闭客户端连接。
 */
func (sc *SimpleRpcClient) Close() {
	sc.ClearHeader()
	sc.ClearCookies()
	if sc.conn != nil {
		sc.conn.Close()
	}
}

/*
 * rpc请求，采用post方式，请求结果返回为json数据，并序列化到 respJson 数据结构上。
 * @param reqUrl  请求url，不需要包含协议部分的信息，直接配置相对路径，例如： /web2/restful/jsoninfo/0/1
 * @param reqBody  请求体，例如：请求参数信息， a=1&b=2&c=3
 * @param respJson 返回的json数据内容。
 * @return err 如果有错误，就抛出。
 */
func (sc *SimpleRpcClient) DoJson(reqUrl string, reqBody string, respJson interface{}) (err error) {
	contentType := "application/x-www-form-urlencoded"
	log4go.InfoLog(message.INF_MSSV_09019, "json", contentType, "POST", fmt.Sprintf(sc.prefix, reqUrl)+"___"+fmt.Sprintf("%s:%d", sc.Ip, sc.Port), reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf(sc.prefix, reqUrl), strings.NewReader(reqBody))
	if err != nil {
		log4go.ErrorLog(err)
		return err
	}

	//req.Header.Set(key, value)
	req.Header.Set("Content-Type", contentType)
	sc.addUSN(req)
	if sc.header != nil {
		for k, v := range sc.header {
			req.Header.Set(k, v)
		}
	}

	if sc.cookies != nil {
		for _, c := range sc.cookies {
			req.AddCookie(c)
		}
	}

	resp, err1 := sc.client.Do(req)
	//	resp, err1 := sc.client.Post(fmt.Sprintf("http://tcp%s", reqUrl), "application/x-www-form-urlencoded", strings.NewReader(reqBody))
	if err1 == nil {
		sc.respCookies = resp.Cookies()
		if resp.StatusCode == 401 {
			err99 := errors.New(message.ERR_MSSV_39036.String())
			log4go.ErrorLog(err99.Error())
			return err99
		}
		defer resp.Body.Close()
		body, err2 := ioutil.ReadAll(resp.Body) //读取服务器返回的信息
		if err2 != nil {
			log4go.ErrorLog(err2)
			err = err2
		} else {
			err3 := json.Unmarshal(body, respJson)
			if err3 == nil {

			} else {
				log4go.ErrorLog(err3)
				err = err3
			}
		}
	} else {
		log4go.ErrorLog(err1)
		err = err1
	}
	return
}

/*
 * rpc请求，采用post方式，查询字符串数据。
 * @param reqUrl  请求url，不需要包含协议部分的信息，直接配置相对路径，例如： /web2/restful/jsoninfo/0/1
 * @param reqBody  请求体，例如：请求参数信息， a=1&b=2&c=3
 * @return returnDatas 返回数据内容。
 * @return err 如果有错误，就抛出。
 */
func (sc *SimpleRpcClient) DoPostDatas(reqUrl string, reqBody string) (returnDatas string, err error) {
	returnDatas = ""

	contentType := "application/x-www-form-urlencoded"
	log4go.InfoLog(message.INF_MSSV_09019, "postDatas", contentType, "POST", fmt.Sprintf(sc.prefix, reqUrl), reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf(sc.prefix, reqUrl), strings.NewReader(reqBody))
	if err != nil {
		return returnDatas, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	sc.addUSN(req)
	if sc.header != nil {
		for k, v := range sc.header {
			req.Header.Set(k, v)
		}
	}

	if sc.cookies != nil {
		for _, c := range sc.cookies {
			req.AddCookie(c)
		}
	}

	resp, err1 := sc.client.Do(req)
	//	resp, err1 := sc.client.Post(fmt.Sprintf("http://tcp%s", reqUrl), "application/x-www-form-urlencoded", strings.NewReader(reqBody))
	if err1 == nil {
		sc.respCookies = resp.Cookies()
		if resp.StatusCode == 401 {
			err99 := errors.New(message.ERR_MSSV_39036.String())
			log4go.ErrorLog(err99.Error())
			return "", err99
		}
		defer resp.Body.Close()
		body, err2 := ioutil.ReadAll(resp.Body) //读取服务器返回的信息
		if err2 != nil {
			err = err2
		} else {
			returnDatas = string(body)
		}
	} else {
		err = err1
	}
	return
}

/*
 * rpc请求，采用get方式，查询字符串数据。
 * @param reqUrl  请求url，不需要包含协议部分的信息，直接配置相对路径，例如： /web2/restful/jsoninfo/0/1
 * @return returnDatas 返回数据内容。
 * @return err 如果有错误，就抛出。
 */
func (sc *SimpleRpcClient) DoGetDatas(reqUrl string) (returnDatas string, err error) {
	returnDatas = ""

	log4go.InfoLog(message.INF_MSSV_09019, "getDatas", "", "GET", fmt.Sprintf(sc.prefix, reqUrl), "")

	req, err := http.NewRequest("GET", fmt.Sprintf(sc.prefix, reqUrl), nil)
	if err != nil {
		return returnDatas, err
	}
	sc.addUSN(req)
	if sc.header != nil {
		for k, v := range sc.header {
			req.Header.Set(k, v)
		}
	}

	if sc.cookies != nil {
		for _, c := range sc.cookies {
			req.AddCookie(c)
		}
	}

	resp, err1 := sc.client.Do(req)

	//resp, err1 := sc.client.Get(fmt.Sprintf("http://tcp%s", reqUrl))
	if err1 == nil {
		sc.respCookies = resp.Cookies()
		if resp.StatusCode == 401 {
			err99 := errors.New(message.ERR_MSSV_39036.String())
			log4go.ErrorLog(err99.Error())
			return "", err99
		}
		defer resp.Body.Close()
		body, err2 := ioutil.ReadAll(resp.Body) //读取服务器返回的信息
		if err2 != nil {
			err = err2
		} else {
			returnDatas = string(body)
		}
	} else {
		err = err1
	}
	return
}

/*
提供数据流下载成文件

@param reqUrl  请求url，不需要包含协议部分的信息，直接配置相对路径，例如： /web2/restful/jsoninfo/0/1
@param reqBody  请求体，例如：请求参数信息， a=1&b=2&c=3
@param localFilePath  需要保存到的本地文件路径。
*/
func (sc *SimpleRpcClient) DoDownload(reqUrl string, reqBody string, localFilePath string) (err error) {

	contentType := "application/x-www-form-urlencoded"
	log4go.InfoLog(message.INF_MSSV_09019, "download", contentType, "POST", fmt.Sprintf(sc.prefix, reqUrl), reqBody)
	req, err := http.NewRequest("POST", fmt.Sprintf(sc.prefix, reqUrl), strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	sc.addUSN(req)
	if sc.header != nil {
		for k, v := range sc.header {
			req.Header.Set(k, v)
		}
	}

	if sc.cookies != nil {
		for _, c := range sc.cookies {
			req.AddCookie(c)
		}
	}

	resp, err1 := sc.client.Do(req)

	//resp, err1 := sc.client.Post(fmt.Sprintf("http://tcp%s", reqUrl), "application/x-www-form-urlencoded", strings.NewReader(reqBody))
	if err1 == nil {
		sc.respCookies = resp.Cookies()
		if resp.StatusCode == 401 {
			err99 := errors.New(message.ERR_MSSV_39036.String())
			log4go.ErrorLog(err99.Error())
			return err99
		}
		defer resp.Body.Close()
		f, err2 := os.Create(localFilePath)
		if err2 != nil {
			err = err2
		}
		defer f.Close()
		io.Copy(f, resp.Body)

	} else {
		err = err1
	}
	return
}

/*
  发起文件上传操作，整个调用模式：
	errUpload := sc.PreUpload("/web2/upload/uploadfile").
		AddFile("path", "/sample1/postgresql-42.1.4.jar", "pgsql20009.jar").
		AddParam("alias", "pgsql2.jar").
		AddParam("author", "jack2").
		AddParam("name", "n1").
		AddParam("name", "n2").DoUpload()
  只有调用DoUpload 方法才实际进行文件上传操作。

  @params reqUrl  请求url地址，相对地址。

*/
func (sc *SimpleRpcClient) PreUpload(reqUrl string) *UploadOperator {
	uo := &UploadOperator{reqUrl: reqUrl}
	uo.fileParamKey = make([]string, 0, 1)
	uo.fileFullFilePath = make([]string, 0, 1)
	uo.fileRenameValue = make([]string, 0, 1)
	uo.fieldParamKey = make([]string, 0, 1)
	uo.fieldParamValue = make([]string, 0, 1)
	uo.streamKeyNames = make([]string, 0, 1)
	uo.streamReaders = make([]io.Reader, 0, 1)
	uo.fileLength = 0
	uo.fieldLength = 0
	uo.err = nil
	uo.simpleRpcClient = sc
	return uo
}

/**

 数据流下载
 请求和响应接近原生，采用流/字节方式交互

@params reqUrl  请求url地址，相对地址。
@params reqBody 请求实体
@return  (io.ReadCloser, error) 返回数据 和错误信息，io.ReadCloser的处理，可以使用 ioutil.ReadAll(reader) 转换成 []byte 。
*/
func (sc *SimpleRpcClient) DoStreamDown(reqUrl string, reqBody io.Reader) (io.ReadCloser, error) {

	contentType := "application/x-www-form-urlencoded"
	log4go.InfoLog(message.INF_MSSV_09019, "stream", contentType, "POST", fmt.Sprintf(sc.prefix, reqUrl), "")

	req, err0 := http.NewRequest("POST", fmt.Sprintf(sc.prefix, reqUrl), reqBody)
	if err0 != nil {
		return nil, err0
	}
	req.Header.Set("Content-Type", contentType)
	sc.addUSN(req)
	if sc.header != nil {
		for k, v := range sc.header {
			req.Header.Set(k, v)
		}
	}

	if sc.cookies != nil {
		for _, c := range sc.cookies {
			req.AddCookie(c)
		}
	}

	resp, err := sc.client.Do(req)
	//	resp, err := sc.client.Post(fmt.Sprintf("http://tcp%s", reqUrl), "application/x-www-form-urlencoded", reqBody)

	if err == nil {
		sc.respCookies = resp.Cookies()
		if resp.StatusCode == 401 {
			err99 := errors.New(message.ERR_MSSV_39036.String())
			log4go.ErrorLog(err99.Error())
			return nil, err99
		} else if resp.StatusCode == 404 {
			log4go.ErrorLog(message.ERR_MSSV_39037)
			return nil, errors.New(message.ERR_MSSV_39037.String())
		} else if resp.StatusCode == 500 {
			log4go.ErrorLog(message.ERR_MSSV_39038)
			return nil, errors.New(message.ERR_MSSV_39038.String())
		} else if resp.Header.Get("Content-Type") != "application/octet-stream" {
			log4go.ErrorLog(message.ERR_MSSV_39039, resp.Header.Get("Content-Type"))
			return nil, errors.New(fmt.Sprintf(message.ERR_MSSV_39039.String(), resp.Header.Get("Content-Type")))
		}
		return resp.Body, nil
	} else {
		return nil, err
	}
}

func (sc *SimpleRpcClient) DoHead(reqUrl string) (string, error) {

	return sc.head_nest(reqUrl, false)
}

func (sc *SimpleRpcClient) head_nest(reqUrl string, ignore bool) (string, error) {

	req, err0 := http.NewRequest("HEAD", fmt.Sprintf(sc.prefix, reqUrl), nil)
	if err0 != nil {
		return "", err0
	}
	sc.addUSN(req)
	if sc.header != nil {
		for k, v := range sc.header {
			req.Header.Set(k, v)
		}
	}

	if sc.cookies != nil {
		for _, c := range sc.cookies {
			req.AddCookie(c)
		}
	}

	resp, err := sc.client.Do(req)

	//resp, err := sc.client.Head(fmt.Sprintf("http://tcp%s", reqUrl))
	if err == nil {
		sc.respCookies = resp.Cookies()
		if !ignore {
			if resp.StatusCode == 401 {
				err99 := errors.New(message.ERR_MSSV_39036.String())
				log4go.ErrorLog(err99.Error())
				return "", err99
			} else if resp.StatusCode == 404 {
				log4go.ErrorLog(message.ERR_MSSV_39037)
				return "", errors.New(message.ERR_MSSV_39037.String())
			} else if resp.StatusCode == 500 {
				log4go.ErrorLog(message.ERR_MSSV_39038)
				return "", errors.New(message.ERR_MSSV_39038.String())
			}
		}
		defer resp.Body.Close()
		body, err2 := ioutil.ReadAll(resp.Body) //读取服务器返回的信息
		if err2 != nil {
			err99 := errors.New(fmt.Sprintf(message.ERR_MSSV_39040.String(), err2))
			return "", err99
		} else {
			return string(body), nil
		}

	} else {
		return "", err
	}
}

/*
在请求头添加全局唯一序列号信息
*/
func (sc *SimpleRpcClient) addUSN(req *http.Request) {
	if env.Switch_CallChain { // 设置 唯一序列号
		if log4go.GetCallChain() != nil {
			caller := log4go.GetCallChain().GetCurrentCaller()
			if caller != nil {
				usn := caller.GetSerialnumber()
				if usn != "" {
					req.Header.Set(env.Env_Web_Header_Unique_Serial_Number, usn)
				}
			}
		}

	}

}

/*
处理文件上传操作
*/
type UploadOperator struct {
	reqUrl string //请求地址，相对地址

	streamReaders  []io.Reader // 基于流存储模式
	streamKeyNames []string

	fileParamKey     []string
	fileFullFilePath []string
	fileRenameValue  []string
	fieldParamKey    []string
	fieldParamValue  []string
	fileLength       int
	fieldLength      int
	err              error
	simpleRpcClient  *SimpleRpcClient
}

/**
添加上传的文件，可以同时添加多个上传文件。实际不进行上传操作，只有出发RunUpload，才完成上传。

@param  paramKey  form表单的名称
@param  fullFilePath  上传文件的本地全路径
@param  renameValue   上传后的文件名，如果是nil或者“”，那么就用  fullFilePath的文件名
*/
func (uo *UploadOperator) AddFile(paramKey, fullFilePath, renameValue string) *UploadOperator {
	if paramKey != "" && fullFilePath != "" {
		if utils.ExistFile(fullFilePath) { // 判断文件存在
			if renameValue == "" { // 如果没有指定保存的文件名，就用现有文件名
				renameValue = utils.BaseFileName(fullFilePath)
			}
			uo.fileParamKey = append(uo.fileParamKey, paramKey)
			uo.fileFullFilePath = append(uo.fileFullFilePath, fullFilePath)
			uo.fileRenameValue = append(uo.fileRenameValue, renameValue)
			uo.fileLength = uo.fileLength + 1
		} else {
			uo.writeErr(fmt.Sprintf(message.ERR_MSSV_39041.String(), uo.fileLength, fullFilePath))
		}
	} else {
		uo.writeErr(fmt.Sprintf(message.ERR_MSSV_39042.String(), uo.fileLength, paramKey, fullFilePath))
	}

	return uo
}

func (uo *UploadOperator) AddStream(streamKeyName string, streamReader io.Reader) *UploadOperator {

	uo.streamKeyNames = append(uo.streamKeyNames, streamKeyName)
	uo.streamReaders = append(uo.streamReaders, streamReader)
	return uo
}

/**
添加表单数据，可以同时添加多个。实际不进行上传操作，只有出发RunUpload，才完成上传。

@param  paramKey  form表单的名称
@param  paramValue  form表单的值
*/
func (uo *UploadOperator) AddParam(paramKey, paramValue string) *UploadOperator {
	if paramKey != "" && paramValue != "" {
		uo.fieldParamKey = append(uo.fieldParamKey, paramKey)
		uo.fieldParamValue = append(uo.fieldParamValue, paramValue)
		uo.fieldLength = uo.fieldLength + 1
	} else {

		uo.writeErr(fmt.Sprintf(message.ERR_MSSV_39043.String(), uo.fieldLength, paramKey, paramValue))
	}
	return uo
}

/**
进行上传操作，实际发起上传请求的方法。
*/
func (uo *UploadOperator) DoUpload() (string, error) {
	if uo.err != nil {
		return "", uo.err
	}
	/*  此处处理方式，是直接把数据读取到内存中，然后再上传，存在内存大量占用的隐患。
	bodyBuf := &bytes.Buffer{}

	bodyWriter := multipart.NewWriter(bodyBuf)

	contentType := bodyWriter.FormDataContentType()

	for i := 0; i < uo.fileLength; i = i + 1 {
		fileWriter, err := bodyWriter.CreateFormFile(uo.fileParamKey[i], uo.fileRenameValue[i])
		if err != nil {
			uo.writeErr(err.Error())
			return "", uo.err
		}
		fileHander, err := os.Open(uo.fileFullFilePath[i])
		if err != nil {
			uo.writeErr(err.Error())
			return "", uo.err
		}
		defer fileHander.Close()
		io.Copy(fileWriter, fileHander)
	}

	for i := 0; i < uo.fieldLength; i = i + 1 {
		bodyWriter.WriteField(uo.fieldParamKey[i], uo.fieldParamValue[i])
	}

	bodyWriter.Close()

	log4go.InfoLog(message.INF_MSSV_09019, "upload", "multipart/form-data", "POST", fmt.Sprintf(uo.simpleRpcClient.prefix, uo.reqUrl), "")
	log4go.InfoLog(bodyBuf.Bytes())
	req, err0 := http.NewRequest("POST", fmt.Sprintf(uo.simpleRpcClient.prefix, uo.reqUrl), bodyBuf)
	if err0 != nil {
		return "", err0
	}

	req.Header.Set("Content-Type", contentType)
	uo.simpleRpcClient.addUSN(req)
	if uo.simpleRpcClient.header != nil {
		for k, v := range uo.simpleRpcClient.header {
			req.Header.Set(k, v)
		}
	}

	resp, err := uo.simpleRpcClient.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		uo.writeErr(err.Error())
		return "", uo.err
	}

	*/

	httpUploadBodyReader := new(HttpUploadBodyReader)
	httpReaderErr := httpUploadBodyReader.Init(uo.streamKeyNames, uo.streamReaders,
		uo.fileParamKey, uo.fileFullFilePath, uo.fileRenameValue,
		uo.fieldParamKey, uo.fieldParamValue)
	if httpReaderErr != nil {
		log4go.InfoLog("httpReaderErr:", httpReaderErr)
	}

	defer httpUploadBodyReader.Close()

	contentType := httpUploadBodyReader.w.FormDataContentType()

	log4go.InfoLog(message.INF_MSSV_09019, "upload", "multipart/form-data", "POST", fmt.Sprintf(uo.simpleRpcClient.prefix, uo.reqUrl), "")
	req, err0 := http.NewRequest("POST", fmt.Sprintf(uo.simpleRpcClient.prefix, uo.reqUrl), httpUploadBodyReader)

	if err0 != nil {
		return "", err0
	}

	req.Header.Set("Content-Type", contentType)
	uo.simpleRpcClient.addUSN(req)
	if uo.simpleRpcClient.header != nil {
		for k, v := range uo.simpleRpcClient.header {
			req.Header.Set(k, v)
		}
	}

	if uo.simpleRpcClient.cookies != nil {
		for _, c := range uo.simpleRpcClient.cookies {
			req.AddCookie(c)
		}
	}

	resp, err := uo.simpleRpcClient.client.Do(req)
	if err != nil {
		uo.writeErr(err.Error())
		return "", uo.err
	}
	defer resp.Body.Close()

	uo.simpleRpcClient.respCookies = resp.Cookies()

	//resp_body, _ := ioutil.ReadAll(resp.Body)  // 暂时不考虑返回
	status := resp.StatusCode

	resp_body, _ := ioutil.ReadAll(resp.Body)

	if status >= 400 {
		uo.writeErr(fmt.Sprintf(message.ERR_MSSV_39044.String(), status, string(resp_body)))
	}

	return string(resp_body), uo.err
}

/**
记录错误信息
*/
func (uo *UploadOperator) writeErr(errInfo string) {
	if uo.err == nil {
		uo.err = errors.New(errInfo)
	} else {
		uo.err = errors.New(fmt.Sprintf(" %s \n %s ", uo.err.Error(), errInfo))
	}

}

type HttpUploadBodyReader struct {
	w *RpcUploadWriter

	streamReaders  []io.Reader // 基于流存储模式
	streamKeyNames []string

	fileParamKey     []string // 基于文件存储模式
	fileFullFilePath []string
	fileRenameValue  []string

	lastFileHander *os.File // 关闭文件句柄

	fieldParamKey   []string
	fieldParamValue []string

	nextReader      io.Reader
	pointerOfReader int
	pointerOfFile   int
}

func (hubr *HttpUploadBodyReader) Init(streamKeyNames []string, streamReaders []io.Reader,
	fileParamKey []string, fileFullFilePath []string, fileRenameValue []string,
	fieldParamKey []string, fieldParamValue []string) error {
	hubr.w = NewRpcUploadWriter()

	hubr.streamKeyNames = streamKeyNames
	hubr.streamReaders = streamReaders

	hubr.fileParamKey = fileParamKey
	hubr.fileFullFilePath = fileFullFilePath
	hubr.fileRenameValue = fileRenameValue
	hubr.fieldParamKey = fieldParamKey
	hubr.fieldParamValue = fieldParamValue
	hubr.nextReader = nil
	hubr.lastFileHander = nil
	hubr.pointerOfFile = -1
	hubr.pointerOfReader = -1

	return nil
}

func (hubr *HttpUploadBodyReader) Close() {
	if hubr.lastFileHander != nil {
		hubr.lastFileHander.Close()
		hubr.lastFileHander = nil
	}
	if hubr.nextReader != nil {
		hubr.nextReader = nil
	}
	hubr.w = nil
}

func (hubr *HttpUploadBodyReader) getFieldReader() error {
	formFieldBuf := &bytes.Buffer{}
	for i := 0; i < len(hubr.fieldParamKey); i = i + 1 {
		b1, b2, err := hubr.w.WriteField(hubr.fieldParamKey[i], hubr.fieldParamValue[i])
		if err == nil {
			formFieldBuf.Write(b1)
			formFieldBuf.Write(b2)
		} else {
			return err
		}
	}
	hubr.nextReader = formFieldBuf
	return nil
}

func (hubr *HttpUploadBodyReader) getNextReader() error {
	lengthOfFile := len(hubr.fileFullFilePath)
	lengthOfReader := len(hubr.streamKeyNames)
	if hubr.pointerOfFile == -1 && hubr.pointerOfReader == -1 { // 如果没有进行文件或者数据流处理，表示还没有初始化
		err := hubr.getFieldReader()
		if err != nil {
			return err
		} else {
			hubr.pointerOfFile = 0
			hubr.pointerOfReader = 0
			return nil
		}
	}

	if lengthOfFile > 0 { // 表示有文件需要传输
		if hubr.lastFileHander != nil {
			hubr.lastFileHander.Close()
			hubr.lastFileHander = nil
		}
		if hubr.pointerOfFile < lengthOfFile*2 {
			inumOfPath := hubr.pointerOfFile / 2
			if hubr.pointerOfFile%2 == 0 { // 设置文件头标示
				endBuf := &bytes.Buffer{}
				b1, err := hubr.w.CreateFormFileStart(hubr.fileParamKey[inumOfPath], hubr.fileRenameValue[inumOfPath])
				if err != nil {
					return err
				} else {
					endBuf.Write(b1)
					hubr.nextReader = endBuf
				}
			} else { // 设置文件句柄
				fileHander, err := os.Open(hubr.fileFullFilePath[inumOfPath])
				if err != nil {
					return err
				} else {
					hubr.nextReader = fileHander

					hubr.lastFileHander = fileHander
				}
			}
			hubr.pointerOfFile = hubr.pointerOfFile + 1
			return nil
		} else if hubr.pointerOfFile == lengthOfFile*2 {
			endBuf := &bytes.Buffer{}
			b1, err := hubr.w.End()
			if err != nil {
				return err
			} else {
				endBuf.Write(b1)
				hubr.nextReader = endBuf
			}
			hubr.pointerOfFile = hubr.pointerOfFile + 1
			return nil
		} else {
			hubr.nextReader = nil
			if hubr.lastFileHander != nil {
				hubr.lastFileHander.Close()
				hubr.lastFileHander = nil
			}
		}
	}

	if lengthOfReader > 0 { // 表示有数据流需要传输
		if hubr.pointerOfReader < lengthOfReader*2 {
			inumOfPath := hubr.pointerOfReader / 2
			if hubr.pointerOfReader%2 == 0 { // 设置文件头标示
				endBuf := &bytes.Buffer{}
				b1, err := hubr.w.CreateFormFileStart(hubr.streamKeyNames[inumOfPath], hubr.streamKeyNames[inumOfPath])
				if err != nil {
					return err
				} else {
					endBuf.Write(b1)
					hubr.nextReader = endBuf
				}
			} else { // 设置文件句柄
				hubr.nextReader = hubr.streamReaders[inumOfPath]
			}
			hubr.pointerOfReader = hubr.pointerOfReader + 1
			return nil
		} else if hubr.pointerOfReader == lengthOfReader*2 {
			endBuf := &bytes.Buffer{}
			b1, err := hubr.w.End()
			if err != nil {
				return err
			} else {
				endBuf.Write(b1)
				hubr.nextReader = endBuf
			}
			hubr.pointerOfReader = hubr.pointerOfReader + 1
			return nil
		} else {
			hubr.nextReader = nil
		}
	}

	return nil
}

func (hubr *HttpUploadBodyReader) isEndReader() bool {
	lengthOfFile := len(hubr.fileFullFilePath)
	if hubr.pointerOfFile > lengthOfFile*2 {
		return true
	} else {
		return false
	}
}

func (hubr *HttpUploadBodyReader) Read(p []byte) (int, error) {
	if hubr.nextReader == nil && !hubr.isEndReader() {
		err := hubr.getNextReader()
		if err != nil {
			return -1, err
		}
	}
	len := len(p)
	pos := 0
	if hubr.nextReader != nil {
		for {
			n, err := hubr.nextReader.Read(p[pos:])

			if err == io.EOF {
				pos = pos + n
				if pos == 0 {
					return 0, io.EOF // 数据查询结束，err信息标示必须是：io.EOF
				} else {
					return pos, nil
				}

			} else if err != nil {
				log4go.ErrorLog(err)
				return -1, err
			} else {
				if n+pos < len {
					pos = n + pos
					hubr.getNextReader()
					if hubr.nextReader == nil {
						break
					}
				} else {
					pos = n + pos
					break
				}
			}
		}
	} else {
		return 0, io.EOF // 数据查询结束，err信息标示必须是：io.EOF
	}
	return pos, nil
}
