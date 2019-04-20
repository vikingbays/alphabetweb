// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*

service 提供 微服务访问模式。主要包括三个部分：1、微服务的管理，2、rpc调用访问.


---------------------------------------------------------------------------------------------------------------------------------------------

0、service微服务依赖与 etcd 3 ，需要初始化etcd3服务。

假设： micro-service配置
   type="etcd"
   endpoints = ["127.0.0.1:2379"]
   username  = "serv001"
   password  = "123456"
   timeout = 2     # 2秒
   root = "awroot_serv001"

启动etcd3服务：
   ETCDCTL_API=3  /Users/vikingbays/golang/AlphabetwebProject/etcd3/etcd  --name etcd01   --listen-client-urls http://0.0.0.0:2379   --advertise-client-urls http://0.0.0.0:2379  --data-dir /Users/vikingbays/golang/AlphabetwebProject/etcd3/data

初始化
   export  ETCDCTL_API=3
   ./etcd3/etcdctl --endpoints="localhost:2379"  user add  root
   #设置密码【rootpw】

   ./etcd3/etcdctl --endpoints="localhost:2379" auth enable
   ./etcd3/etcdctl --endpoints="localhost:2379" --user root:rootpw role add role1
   ./etcd3/etcdctl --endpoints="localhost:2379" --user root:rootpw  role grant-permission role1 --prefix=true  readwrite  awroot_serv001

   ./etcd3/etcdctl --endpoints="localhost:2379" --user root:rootpw  user add serv001
   #设置密码【123456】

   ./etcd3/etcdctl --endpoints="localhost:2379" --user root:rootpw  user grant-role  serv001  role1

验证：
   ./etcd3/etcdctl --endpoints="localhost:2379" --user serv001:123456   get  --prefix   awroot_serv001

---------------------------------------------------------------------------------------------------------------------------------------------


1、微服务的管理，提供针对RPC服务的服务注册、服务发现，服务监控等功能。整个架构设计中，服务端注册只能指向一个注册中心。客户端注册可以指向多个注册中心，也就是服务提供来源与不同域。

1.1）、配置文件说明:

服务端配置和客户端配置都是采用toml文件注册，

服务端配置信息如下： 命名是：ms_server_config.toml

  #register  一个server只能注册到一个服务中心，用于服务端注册到服务中心
  [[server]]                    # 只能配置一个
  groupId="group_octopus01"
  groupName="group_octopus01"
  protocol="rpc_unix"            #rpc_tcp,rpc_unix,rpc_tcp_ssl
  ipType="ipv4"                  # 使用ipv4 还是 ipv6
  ip="${project}/sos_rpc.sock"   # 如果配置 * 表示任意指定有效网卡，
                                 # 如果协议是rpc_unix，写成： ip="${project}/sos_rpc.sock"
                                 # 如果协议是rpc_tcp，写成：127.0.01
  port=0

  #[[server]]                    # 只举例
  #groupId="group_octopus01"
  #groupName="group_octopus01"
  #protocol="rpc_tcp"
  #ipType="ipv4"
  #ip="127.0.0.1"
  #port=10777

  #[[server]]                    # 只举例
  #groupId="group_octopus01"
  #groupName="group_octopus01"
  #protocol="rpc_tcp_ssl"
  #ipType="ipv4"
  #ip="127.0.0.1"
  #port=10766

  #[[server]]                    # 只举例
  #groupId="group_octopus01"
  #groupName="group_octopus01"
  #protocol="rpc_unix"
  #ipType="ipv4"
  #ip="${project}/sos_rpc.sock"
  #port=0

  [[server.register]]
  type="etcd"                      # 当前只支持etcd
  endpoints = ["127.0.0.1:2379"]   # etcd服务地址，可以多个
  username  = "serv001"            # etcd的用户名
  password  = "123456"             # etcd的密码
  timeout = 2                      # 2秒
  root = "awroot_serv001"          # 在etcd中配置的根路径

  [[server.rpcs]]                  # 注册具体的rpc服务，可以多个
  rpcId="db_index"                 # 唯一rpc标示号，结合groupId配套使用
  url="/db/index"                  # 实际地址
  desc="db_index...desc...."       # 文本表述
  available=true

  [[server.rpcs]]
  rpcId="db_insert"
  url="/db/insert"
  desc="db_insert...desc...."
  available=true

  [[server.rpcs]]
  rpcId="db_delete"
  url="/db/delete"
  desc="db_delete...desc...."
  available=true

客户端配置信息如下： 命名是：ms_client_config_01.toml

如果同时连接多个服务中心，那么就需要ms_client_config_01.toml，ms_client_config_02.toml，ms_client_config_03.toml
  [[client]]
  groupIds=["group_octopus01"]    # 可以使用的groupid信息，可以多个
  maxPoolSize=20
  reqPerConn=20                   # 如果是http2才有用，支持一个连接多路由复用


  [[client.register]]
  type = "etcd"                   # 当前只支持etcd
  endpoints = ["127.0.0.1:2379"]  # etcd服务地址，可以多个
  username  = "serv001"           # etcd的用户名
  password  = "123456"            # etcd的密码
  timeout = 2                     # 2秒
  root = "awroot_serv001"         # 在etcd中配置的根路径

1.2）、微服务调用说明:

client调用server端，是基于rpc调用访问封装的。

Json请求：
  func AskJson_MS(groupId string, rpcId string, reqParams interface{}, returnRespJson interface{}) error

 【微服务】，采用post方式，请求结果返回为json数据，并序列化到 respJson 数据结构上。

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

Stream请求（支持响应端流下载）：
  func AskStreamDown_MS(groupId string, rpcId string, reqParams interface{}) (io.ReadCloser, error)

 【微服务】，下载数据流，在某些场景下AskDownload_MS是他的简化版本（直接下载到本地文件），而 AskStreamDown_MS 返回的是是一个数据流，适合大数据传输。


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

Stream请求（支持响应端流上传）：
  func AskStreamUp_MS(groupId string, rpcId string, streamUpInfos []StreamUpInfo, textParamsBean interface{}) (string, error)

  微服务】，提供上传功能，不只是文件，可以上传流式数据

  @params groupId  根据 groupId 和 rpcId 获取一个rpc连接
  @params rpcId
  @params streamUpInfos 设置stream 流数据（如文件上传），对应类：[]service.StreamUpInfo（fieldName，流对象）
  @params textParamsBean 设置文本数据
  @return (string, error) 返回结果

  例如：
    streamUpInfos := make([]service.StreamUpInfo, len(p2))
    for num, p2_str := range p2 {
      f1, _ := os.Open(p2_str)
      if len(p3) == 0 {
        streamUpInfos[num] = service.StreamUpInfo{StreamKeyName: "unknown_stream", StreamReader: f1}
      } else {
        streamUpInfos[num] = service.StreamUpInfo{StreamKeyName: p3[num] + "_stream", StreamReader: f1}
      }
      defer f1.Close()
    }
    service.AskStreamUp_MS("group_octopus01", "upload_uploadstream", streamUpInfos, &uploadTextParamBean)

数据文件上传：

  func AskUpload_MS(groupId string, rpcId string, uploadFileUpInfos []UploadFileUpInfo, textParamsBean interface{}) (string, error)

  【微服务】，提供上传文件功能

  @params groupId  根据 groupId 和 rpcId 获取一个rpc连接
  @params rpcId
  @param uploadFileUpInfos    需要上传文件的信息 []service.UploadFileUpInfo（fieldName，路径，别名等）
  @param textParamBeans    需要请求的文本参数，必须时结构体形式
  @return (string, error) 返回结果

  例如：
    d1, err := service.AskUpload_MS("group_octopus01", "upload_uploadfile",&uploadUploadFileBean, &uploadTextParamBean)

数据文件下载（预先存储到本地）：

  func AskDownload_MS(groupId string, rpcId string, reqParams interface{}, localFilePath string) (err error)

  【微服务】，提供数据流下载成文件

  @params groupId  根据 groupId 和 rpcId 获取一个rpc连接
  @params rpcId
  @param reqParams  请求体，支持的格式有两种：
                           1） 字符串方式（string） 例如：请求参数信息， a=1&b=2&c=3
                           2） 结构体方式（struct），可以是指针 例如：请求参数信息， &Bean1
  @param localFilePath  需要保存到的本地文件路径。

GET 通用请求：

 func AskCommonGetDatas_MS(groupId string, rpcId string, reqParams interface{}) (string, error)

 【微服务】，采用get方式，请求数据并返回查询结果（字符串形式）。

 @params groupId  根据 groupId 和 rpcId 获取一个rpc连接
 @params rpcId
 @param reqParams 请求体，支持的格式有两种：
		     1） 字符串方式（string） 例如：请求参数信息， a=1&b=2&c=3
		     2） 结构体方式（struct），可以是指针 例如：请求参数信息， &Bean1
 @return string 返回查询结果
 @return err 如果有错误，就抛出。

POST通用请求：

  func AskCommonPostDatas_MS(groupId string, rpcId string, reqParams interface{}) (string, error)
  说明【微服务】，采用post方式，请求数据并返回查询结果（字符串形式）。

  @param groupId 根据 groupId 和 rpcId 获取一个rpc连接
  @param rpcId
  @param reqParams 请求体，支持的格式有两种：
                    1） 字符串方式（string） 例如：请求参数信息， a=1&b=2&c=3
                    2） 结构体方式（struct），可以是指针 例如：请求参数信息， &Bean1
  @return string 返回查询结果
  @return err 如果有错误，就抛出。


---------------------------------------------------------------------------------------------------------------------------------------------


2、rpc调用访问，用于访问远程服务的方法。(当前提供的服务都是采用post方式)。涉及的代码以：“rpc”开头。

2.1). 请求数据返回json格式：
  sc := &SimpleRpcClient{protocol: "rpc_tcp", addr: "127.0.0.1", port: 7777}    ## 创建SimpleRpcClient对象
  err := sc.Connect()                                                           ## 创建连接
  sc.DoJson("/web2/restful/jsoninfo/0/1", "", &respBody)                       ## 调用Json服务，并写入到respBody结果集中，respBody预先定义好格式。
  sc.Close()                                                                    ## 关闭连接

2.2). 文件上传操作：
  errUpload := sc.PreUpload("/web2/upload/uploadfile").
	          AddFile("path", "/sample1/doc/postgresql-42.1.4.jar", "pgsql20009.jar").      ## 设定需要上传的文件
		  AddParam("alias", "pgsql2.jar").                                              ## 设置参数数据，可以多个
		  AddParam("author", "jack2").
		  AddParam("name", "n1").
		  AddParam("name", "n2").
		  DoUpload()                                                                    ## 进行上传操作，实际发起请求的地方

2.3). 文件下载操作：
  sc.DoDownload("/web2/download/do_download", "filepath=/sample1/postgresql-42.1.4.jar&aliasname=pg.jar", "/tu.jar")

2.4). 接近原生数据访问：（byte数据处理）
  reader, errStream := sc.DoStream("/web2/db/query", strings.NewReader("min=0&max=100"))
  if errStream != nil {
	t.Log(errStream)
  }
  bytes, errBytes := ioutil.ReadAll(reader)
  if errBytes != nil {
	t.Log(errBytes)
  }
  defer reader.Close()
  t.Log(string(bytes))

可以调用 go test -v -bench="simplerpcclient_test.go" 测试。

*/
package service
