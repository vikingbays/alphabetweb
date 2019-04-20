// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package service

import (
	"alphabet/core/webutils"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func Init() {

	flexServicePoolGroup.NewFlexServicePoolGroup()

	initServAndClient()

}

/*
【微服务】根据groupId从微服务连接池中获取一个连接
 @param groupId  groupId唯一标示号
 @return *FlexRpcClient  返回rpc客户端对象
*/
func GetServiceConnection_MS(groupId string) *FlexRpcClientWrapper {
	return flexServicePoolGroup.GetConnection(groupId)
}

/*
【微服务】释放一个rpc客户端对象，重新放入到微服务连接池中
 @param fpc  rpc客户端对象
*/
func ReleaseServiceConnection_MS(fpc *FlexRpcClientWrapper) {
	flexServicePoolGroup.ReleaseConnection(fpc)
}

/*
说明【微服务】，采用post方式，请求数据并返回查询结果（字符串形式）。
  @param groupId 根据 groupId 和 rpcId 获取一个rpc连接
  @param rpcId
  @param reqParams 请求体，支持的格式有两种：
                      1） 字符串方式（string） 例如：请求参数信息， a=1&b=2&c=3
                      2） 结构体方式（struct），可以是指针 例如：请求参数信息， &Bean1
  @return string 返回查询结果
  @return err 如果有错误，就抛出。
*/
func AskCommonPostDatas_MS(groupId string, rpcId string, reqParams interface{}) (string, error) {
	frc, reqUrl, ticket, ticketForLast, err := getNestServiceConnection_MS(groupId, rpcId)
	if err == nil {
		var sc *SimpleRpcClient = nil
		defer func() {
			if sc != nil {
				sc = nil
			}
			ReleaseServiceConnection_MS(frc)
		}()
		sc = frc.CloneRpcClient()
		sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_NAME, ticket)
		sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_LAST_NAME, ticketForLast)

		reqParamsStr := ""
		if reqParams != nil {
			if reflect.TypeOf(reqParams).Kind() == reflect.String {
				reqParamsStr = reqParams.(string)
			} else if reflect.TypeOf(reqParams).Kind() == reflect.Struct {
				reqParamsStr = webutils.EncodeHttpParams(reflect.ValueOf(reqParams))
			} else if reflect.TypeOf(reqParams).Kind() == reflect.Ptr {
				if reflect.TypeOf(reqParams).Elem().Kind() == reflect.Struct {
					reqParamsStr = webutils.EncodeHttpParams(reflect.ValueOf(reqParams))
				}
			}
		}

		returnDatas, err2 := sc.DoPostDatas(reqUrl, reqParamsStr)
		return returnDatas, err2
	} else {
		log4go.ErrorLog(message.ERR_MSSV_39029, err.Error())
	}
	return "", err
}

/*
说明【微服务】，采用get方式，请求数据并返回查询结果（字符串形式）。
  @params groupId  根据 groupId 和 rpcId 获取一个rpc连接
  @params rpcId
  @param reqParams 请求体，支持的格式有两种：
		1） 字符串方式（string） 例如：请求参数信息， a=1&b=2&c=3
		2） 结构体方式（struct），可以是指针 例如：请求参数信息， &Bean1
  @return string 返回查询结果
  @return err 如果有错误，就抛出。
*/
func AskCommonGetDatas_MS(groupId string, rpcId string, reqParams interface{}) (string, error) {
	frc, reqUrl, ticket, ticketForLast, err := getNestServiceConnection_MS(groupId, rpcId)
	if err == nil {
		var sc *SimpleRpcClient = nil
		defer func() {
			if sc != nil {
				sc = nil
			}
			ReleaseServiceConnection_MS(frc)
		}()
		sc = frc.CloneRpcClient()
		sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_NAME, ticket)
		sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_LAST_NAME, ticketForLast)

		reqParamsStr := ""
		if reqParams != nil {
			if reflect.TypeOf(reqParams).Kind() == reflect.String {
				reqParamsStr = reqParams.(string)
			} else if reflect.TypeOf(reqParams).Kind() == reflect.Struct {
				reqParamsStr = webutils.EncodeHttpParams(reflect.ValueOf(reqParams))
			} else if reflect.TypeOf(reqParams).Kind() == reflect.Ptr {
				if reflect.TypeOf(reqParams).Elem().Kind() == reflect.Struct {
					reqParamsStr = webutils.EncodeHttpParams(reflect.ValueOf(reqParams))
				}
			}
		}

		returnDatas, err2 := sc.DoGetDatas(fmt.Sprintf("%s?%s", reqUrl, reqParamsStr))
		return returnDatas, err2
	} else {
		log4go.ErrorLog(message.ERR_MSSV_39029, err.Error())
	}
	return "", err
}

/*
说明：【微服务】，采用post方式，请求结果返回为json数据，并序列化到 respJson 数据结构上。
 @params groupId  根据 groupId 和 rpcId 获取一个rpc连接
 @params rpcId
 @param reqParams  请求体，支持的格式有两种：
                          1） 字符串方式（string） 例如：请求参数信息， a=1&b=2&c=3
                          2） 结构体方式（struct），可以是指针 例如：请求参数信息， &Bean1
 @param respJson 返回的json数据内容。
 @return err 如果有错误，就抛出。
例如：
 reqParams 是字符串方式
     service.AskJson_MS("group_octopus01", "db_delete", fmt.Sprintf("delete=%d", 1), &dbInfo)
 reqParams 是结构体方式
     paramDB_Req := api.ParamDB_Req{Delete: paramBean.Delete}
     service.AskJson_MS("group_octopus01", "db_delete", &paramDB_Req, &dbInfo)
*/
func AskJson_MS(groupId string, rpcId string, reqParams interface{}, returnRespJson interface{}) error {
	frc, reqUrl, ticket, ticketForLast, err := getNestServiceConnection_MS(groupId, rpcId)
	if err == nil {
		var sc *SimpleRpcClient = nil
		defer func() {
			if sc != nil {
				sc = nil
			}
			ReleaseServiceConnection_MS(frc)
		}()
		sc = frc.CloneRpcClient()

		sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_NAME, ticket)
		sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_LAST_NAME, ticketForLast)

		reqParamsStr := ""
		if reqParams != nil {
			if reflect.TypeOf(reqParams).Kind() == reflect.String {
				reqParamsStr = reqParams.(string)
			} else if reflect.TypeOf(reqParams).Kind() == reflect.Struct {
				reqParamsStr = webutils.EncodeHttpParams(reflect.ValueOf(reqParams))
			} else if reflect.TypeOf(reqParams).Kind() == reflect.Ptr {
				if reflect.TypeOf(reqParams).Elem().Kind() == reflect.Struct {
					reqParamsStr = webutils.EncodeHttpParams(reflect.ValueOf(reqParams))
				}
			}
		}
		err2 := sc.DoJson(reqUrl, reqParamsStr, returnRespJson)
		return err2
	} else {
		log4go.ErrorLog(message.ERR_MSSV_39029, err.Error())
	}
	return err
}

/*
说明：【微服务】，提供数据流下载成文件
 @params groupId  根据 groupId 和 rpcId 获取一个rpc连接
 @params rpcId
 @param reqParams  请求体，支持的格式有两种：
                          1） 字符串方式（string） 例如：请求参数信息， a=1&b=2&c=3
                          2） 结构体方式（struct），可以是指针 例如：请求参数信息， &Bean1
 @param localFilePath  需要保存到的本地文件路径。
*/
func AskDownload_MS(groupId string, rpcId string, reqParams interface{}, localFilePath string) (err error) {
	frc, reqUrl, ticket, ticketForLast, err := getNestServiceConnection_MS(groupId, rpcId)
	if err == nil {
		var sc *SimpleRpcClient = nil
		defer func() {
			if sc != nil {
				sc = nil
			}
			ReleaseServiceConnection_MS(frc)
		}()
		sc = frc.CloneRpcClient()
		sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_NAME, ticket)
		sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_LAST_NAME, ticketForLast)

		reqParamsStr := ""
		if reqParams != nil {
			if reflect.TypeOf(reqParams).Kind() == reflect.String {
				reqParamsStr = reqParams.(string)
			} else if reflect.TypeOf(reqParams).Kind() == reflect.Struct {
				reqParamsStr = webutils.EncodeHttpParams(reflect.ValueOf(reqParams))
			} else if reflect.TypeOf(reqParams).Kind() == reflect.Ptr {
				if reflect.TypeOf(reqParams).Elem().Kind() == reflect.Struct {
					reqParamsStr = webutils.EncodeHttpParams(reflect.ValueOf(reqParams))
				}
			}
		}
		err2 := sc.DoDownload(reqUrl, reqParamsStr, localFilePath)
		return err2
	} else {
		log4go.ErrorLog(message.ERR_MSSV_39029, err.Error())
	}
	return err
}

/*
说明：【微服务】，下载数据流，在某些场景下AskDownload_MS是他的简化版本（直接下载到本地文件），而 AskStreamDown_MS 返回的是是一个数据流，适合大数据传输。

 @params groupId  根据 groupId 和 rpcId 获取一个rpc连接
 @params rpcId
 @params reqParams 请求实体，支持的格式有两种：
                          1） 字符串方式（string） 例如：请求参数信息， a=1&b=2&c=3 ,会转换成 strings.NewReader("a=1&b=2&c=3")
                          2） 结构体方式（struct），可以是指针 例如：请求参数信息， &Bean1,会转换成 strings.NewReader("....")
													3)  流数据（io.Reader）
 @return  (io.ReadCloser, error) 返回数据 和错误信息，io.ReadCloser的处理，可以使用 ioutil.ReadAll(reader) 转换成 []byte 。
例如：
 reqParams参数是字符串方式：
     reqParams:=fmt.Sprintf("filepath=%s&aliasname=%s", paramDownload.Filepath, paramDownload.Aliasname)
     reader, err := service.AskStream_MS("group_octopus01", "download_do_download", reqParams)
 reqParams参数是结构体方式：
     paramDownload_Req := api.ParamDownload_Req{Filepath: paramDownload.Filepath, Aliasname: paramDownload.Aliasname}
     reader, err := service.AskStream_MS("group_octopus01", "download_do_download", paramDownload_Req)
 reqParams参数是流数据（io.Reader）
     var reqBody io.Reader = strings.NewReader(fmt.Sprintf("filepath=%s&aliasname=%s", paramDownload.Filepath, paramDownload.Aliasname))
     reader, err := service.AskStream_MS("group_octopus01", "download_do_download", reqBody)
*/
func AskStreamDown_MS(groupId string, rpcId string, reqParams interface{}) (io.ReadCloser, error) {
	frc, reqUrl, ticket, ticketForLast, err := getNestServiceConnection_MS(groupId, rpcId)
	if err == nil {
		var sc *SimpleRpcClient = nil
		defer func() {
			if sc != nil {
				sc = nil
			}
			ReleaseServiceConnection_MS(frc)
		}()
		sc = frc.CloneRpcClient()
		sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_NAME, ticket)
		sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_LAST_NAME, ticketForLast)
		var reqParamsStream io.Reader
		var ok bool
		if reqParams == nil {
			reqParamsStream = strings.NewReader("")
		} else {
			reqParamsStream, ok = reqParams.(io.Reader)
			if !ok {
				if reflect.TypeOf(reqParams).Kind() == reflect.String {
					reqParamsStream = strings.NewReader(reqParams.(string))
				} else if reflect.TypeOf(reqParams).Kind() == reflect.Struct {
					reqParamsStream = strings.NewReader(webutils.EncodeHttpParams(reflect.ValueOf(reqParams)))
				} else if reflect.TypeOf(reqParams).Kind() == reflect.Ptr {
					if reflect.TypeOf(reqParams).Elem().Kind() == reflect.Struct {
						reqParamsStream = strings.NewReader(webutils.EncodeHttpParams(reflect.ValueOf(reqParams)))
					}
				}
			}
		}
		respReader, err2 := sc.DoStreamDown(reqUrl, reqParamsStream)
		return respReader, err2
	} else {
		log4go.ErrorLog(message.ERR_MSSV_39029, err.Error())
	}
	return nil, err
}

/*
说明：【微服务】，提供上传功能，不只是文件，可以上传流式数据
	@params groupId  根据 groupId 和 rpcId 获取一个rpc连接
	@params rpcId
	@params streamUpInfos 设置stream 流数据（如文件上传），对应类：[]service.StreamUpInfo（fieldName，流对象）
	@params textParamsBean 设置文本数据
	@return (string, error) 返回结果
*/
func AskStreamUp_MS(groupId string, rpcId string, streamUpInfos []StreamUpInfo, textParamsBean interface{}) (string, error) {

	paramKey_datas := make([]string, 0, 1)
	paramValue_datas := make([]string, 0, 1)
	if reflect.TypeOf(textParamsBean).Kind() == reflect.Struct ||
		(reflect.TypeOf(textParamsBean).Kind() == reflect.Ptr &&
			reflect.TypeOf(textParamsBean).Elem().Kind() == reflect.Struct) {
		paramsMap := webutils.EncodeHttpParamsMap(reflect.ValueOf(textParamsBean))
		for k, v := range paramsMap {
			for _, v1 := range v {
				paramKey_datas = append(paramKey_datas, k)
				paramValue_datas = append(paramValue_datas, v1)
			}
		}
	}

	if len(paramKey_datas) == len(paramValue_datas) {
		frc, reqUrl, ticket, ticketForLast, err := getNestServiceConnection_MS(groupId, rpcId)
		if err == nil {
			var sc *SimpleRpcClient = nil
			defer func() {
				if sc != nil {
					sc = nil
				}
				ReleaseServiceConnection_MS(frc)
			}()
			sc = frc.CloneRpcClient()
			uo := sc.PreUpload(reqUrl)
			for _, streamUpInfo := range streamUpInfos {
				uo.AddStream(streamUpInfo.StreamKeyName, streamUpInfo.StreamReader)
			}
			for num, _ := range paramKey_datas {
				uo.AddParam(paramKey_datas[num], paramValue_datas[num])
			}
			sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_NAME, ticket)
			sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_LAST_NAME, ticketForLast)
			body, err2 := uo.DoUpload()
			return body, err2
		} else {
			log4go.ErrorLog(message.ERR_MSSV_39029, err.Error())
		}
		return "", err
	} else {
		err := errors.New(message.ERR_MSSV_39030.String())
		log4go.ErrorLog(err.Error())
		return "", err
	}

	return "", nil
}

/*
用于上传流方法（AskStreamUp_MS）时，记录上传数据流的信息，
*/
type StreamUpInfo struct {
	StreamKeyName string    // 上传的每个数据流别名
	StreamReader  io.Reader // 需要上传的数据流
}

/*
说明：【微服务】，提供上传文件功能
 @params groupId  根据 groupId 和 rpcId 获取一个rpc连接
 @params rpcId
 @param uploadFileUpInfos    需要上传文件的信息 []service.UploadFileUpInfo（fieldName，路径，别名等）
 @param textParamBeans    需要请求的文本参数，必须时结构体形式
 @return (string, error) 返回结果

例如：
 d1, err := service.AskUpload_MS("group_octopus01", "upload_uploadfile",
            &uploadUploadFileBean, &uploadTextParamBean)
*/
func AskUpload_MS(groupId string, rpcId string, uploadFileUpInfos []UploadFileUpInfo, textParamsBean interface{}) (string, error) {

	fileKeyNames := make([]string, 0, 1)
	filePaths := make([]string, 0, 1)
	fileRenames := make([]string, 0, 1)

	for _, ufuiBean := range uploadFileUpInfos {
		if ufuiBean.FileKeyName == "" {
			continue
		}
		if ufuiBean.FilePath == "" {
			continue
		}
		fileKeyNames = append(fileKeyNames, ufuiBean.FileKeyName)
		filePaths = append(filePaths, ufuiBean.FilePath)
		fileRenames = append(fileRenames, ufuiBean.FileRename)
	}

	paramKey_datas := make([]string, 0, 1)
	paramValue_datas := make([]string, 0, 1)
	if reflect.TypeOf(textParamsBean).Kind() == reflect.Struct ||
		(reflect.TypeOf(textParamsBean).Kind() == reflect.Ptr &&
			reflect.TypeOf(textParamsBean).Elem().Kind() == reflect.Struct) {
		paramsMap := webutils.EncodeHttpParamsMap(reflect.ValueOf(textParamsBean))
		for k, v := range paramsMap {
			for _, v1 := range v {
				paramKey_datas = append(paramKey_datas, k)
				paramValue_datas = append(paramValue_datas, v1)
			}
		}
	}
	return askUpload_MS_nest(groupId, rpcId, fileKeyNames, filePaths, fileRenames,
		paramKey_datas, paramValue_datas)
}

/*
用于上传文件方法（AskUpload_MS）时的第三个参数，记录上传文件的信息，
*/
type UploadFileUpInfo struct {
	FileKeyName string // 上传文件时，在form表单中定义的名称：fieldName
	FilePath    string // 需要上传的本地文件路径
	FileRename  string // 上传文件的别名，重命名，如果为 “” 表示不需要更改。
}

/*
说明：【微服务】，提供上传文件功能
 @params groupId  根据 groupId 和 rpcId 获取一个rpc连接
 @params rpcId
 @param ParamKey_files    文件参数key
 @param fullFilePath_files  本地上传文件路径
 @param renameValue_files  重命名文件名称
 @param paramKey_datas  字符串参数key
 @param paramValue_datas  字符串参数值
 @return (string, error) 返回结果
*/
func askUpload_MS_nest(groupId string, rpcId string, ParamKey_files []string, fullFilePath_files []string, renameValue_files []string,
	paramKey_datas []string, paramValue_datas []string) (string, error) {
	if len(ParamKey_files) == len(fullFilePath_files) && len(ParamKey_files) == len(renameValue_files) {
		if len(paramKey_datas) == len(paramValue_datas) {
			frc, reqUrl, ticket, ticketForLast, err := getNestServiceConnection_MS(groupId, rpcId)
			if err == nil {
				var sc *SimpleRpcClient = nil
				defer func() {
					if sc != nil {
						sc = nil
					}
					ReleaseServiceConnection_MS(frc)
				}()
				sc = frc.CloneRpcClient()
				uo := sc.PreUpload(reqUrl)
				for num, _ := range ParamKey_files {
					uo.AddFile(ParamKey_files[num], fullFilePath_files[num], renameValue_files[num])
				}
				for num, _ := range paramKey_datas {
					uo.AddParam(paramKey_datas[num], paramValue_datas[num])
				}
				sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_NAME, ticket)
				sc.AddHeader(SERVICE_MS_REQ_HEADER_TICKET_LAST_NAME, ticketForLast)
				body, err2 := uo.DoUpload()
				return body, err2
			} else {
				log4go.ErrorLog(message.ERR_MSSV_39029, err.Error())
			}
			return "", err
		} else {
			err := errors.New(message.ERR_MSSV_39030.String())
			log4go.ErrorLog(err.Error())
			return "", err
		}

	} else {
		err := errors.New(message.ERR_MSSV_39030.String())
		log4go.ErrorLog(err.Error())
		return "", err
	}

}

/*
【微服务】从微服务连接池中获取一个对象，并且根据rpcId返回请求地址reqUrl
@param groupId  group标示号
@param rpcId   rpc标示号
@return *FlexRpcClient  rpc客户端连接
@return string  实际请求的相对地址 reqUrl
@return error  如果异常，就抛出错误信息
*/
func getNestServiceConnection_MS(groupId string, rpcId string) (*FlexRpcClientWrapper, string, string, string, error) {
	var smg *ServiceManagerGroup = nil
	for _, clients := range clientsGlobalMap {
		if clients.dataMap[groupId] != nil {
			smg = clients.dataMap[groupId]
			break
		}
	}

	var err error
	if smg != nil {
		sc := GetServiceConnection_MS(groupId)
		if sc != nil {
			addr := smg.Addrs[merge_ip_and_port(sc.Ip, sc.Port)]
			if addr == nil {
				ReleaseServiceConnection_MS(sc)
				err = errors.New(fmt.Sprintf(message.ERR_MSSV_39031.String(), sc.Ip, sc.Port))
				log4go.ErrorLog(err.Error())
			} else {
				webContext := addr.WebContext
				ticket := addr.Ticket
				ticketForLast := addr.TicketForLast
				rpc := smg.Rpcs[rpcId]
				if rpc == nil {
					ReleaseServiceConnection_MS(sc)
					err = errors.New(fmt.Sprintf(message.ERR_MSSV_39032.String(), rpcId))
					log4go.ErrorLog(err.Error())
				} else {
					subUrl := rpc.Url
					return sc, fmt.Sprintf("/%s%s", webContext, subUrl), ticket, ticketForLast, nil
				}
			}
		} else {
			err = errors.New(fmt.Sprintf(message.ERR_MSSV_39033.String(), groupId))
		}
	}
	return nil, "", "", "", err
}

func ValidTicketServer_MS(ticket string, ticketForLast string, rootAppName string) bool {
	if ticket == "" && ticketForLast == "" {
		return false
	} else if serverGlobalMap[rootAppName].ticket == ticket || serverGlobalMap[rootAppName].ticket == ticketForLast ||
		serverGlobalMap[rootAppName].ticketForLast == ticket || serverGlobalMap[rootAppName].ticketForLast == ticketForLast {
		return true
	} else {
		return false
	}
}

/*
  是否有微服务的服务端运行
*/
func IsServerRunable_MS(rootAppname string) bool {
	if serverGlobalMap[rootAppname] != nil && serverGlobalMap[rootAppname].manager != nil &&
		(serverGlobalMap[rootAppname].ticket != "" || serverGlobalMap[rootAppname].ticketForLast != "") {
		return true
	} else {
		return false
	}
}

/*
  判断该url是不是rpc服务
*/
func IsRpc_MS(url0 string, rootAppname string) bool {
	if serverGlobalMap[rootAppname] != nil {
		if serverGlobalMap[rootAppname].rpcUrlMap != nil {
			if serverGlobalMap[rootAppname].rpcUrlMap[url0] == 1 {
				return true
			}
		}
	}
	return false
}
