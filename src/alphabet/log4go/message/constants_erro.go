// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package message

/*
[ERR_CORE_39001]: 如果连接池中无法获取可使用的对象，报错。
*/
var ERR_CORE_39001 MessageType = MessageType{Id: "ERR_CORE_39001",
	Cn: "在连接池基础类中，从对象池中获取对象失败，当前对象是 nil，说明对象创建失败，需要检查对象创建方法（ IObjectFactory.Create() ） 。（本次是 第 %d 次尝试）",
	En: "GetObject in pool is error  , current object is nil . (try get is %d times.)"}

/*
  [ERR_CACH_39002]: 初始化redis的cache连接池时，配置文件加载可能有异常
*/
var ERR_CACH_39002 MessageType = MessageType{Id: "ERR_CACH_39002",
	Cn: "初始化redis时，读取toml配置文件（%s）失败，具体报错信息：%s",
	En: "When init redis, configfile of redis is error , configfile = %s , error info: %s "}

/*
   [ERR_CACH_39003]: 初始化redis的cache连接池时，创建cache对象失败，可能redisi无法连接
*/
var ERR_CACH_39003 MessageType = MessageType{Id: "ERR_CACH_39003",
	Cn: "在Cache连接池中，调用CacheObjectFactory.Create() 创建cache对象失败，返回错误信息： %s",
	En: "The creation of objects in the cache pool is failed . err info: %s "}

/*
   [ERR_CACH_39004]: 初始化redis的cache连接池时，创建cache对象失败，可能redisi无法连接
*/
var ERR_CACH_39004 MessageType = MessageType{Id: "ERR_CACH_39004",
	Cn: "在Cache连接池中，Cache对象验证失败，返回错误信息： %s",
	En: "Cache object in cache pool is invalid  . err info: %s "}

/*
   [ERR_SESS_39005]: 在使用redis session 模式时，获取cache对象报错
*/
var ERR_SESS_39005 MessageType = MessageType{Id: "ERR_SESS_39005",
	Cn: "在RedisSession 模式下，创建session时，获取cache对象报错，返回错误信息： %s",
	En: "When creating session in RedisSession, getting cache object is error  . err info: %s "}

/*
   [ERR_SESS_39006]: 在使用redis session 模式时，获取cache对象报错
*/
var ERR_SESS_39006 MessageType = MessageType{Id: "ERR_SESS_39006",
	Cn: "在RedisSession 模式下，调用Session.Get()时，获取cache对象报错，返回错误信息： %s",
	En: "When using Session.Get() in RedisSession, getting cache object is error  . err info: %s "}

/*
   [ERR_SESS_39007]: 在使用redis session 模式时，获取cache对象报错
*/
var ERR_SESS_39007 MessageType = MessageType{Id: "ERR_SESS_39007",
	Cn: "在RedisSession 模式下，调用Session.Set()时，获取cache对象报错，返回错误信息： %s",
	En: "When using Session.Set() in RedisSession, getting cache object is error  . err info: %s "}

/*
   [ERR_SESS_39008]: 在使用redis session 模式时，获取cache对象报错
*/
var ERR_SESS_39008 MessageType = MessageType{Id: "ERR_SESS_39008",
	Cn: "在RedisSession 模式下，调用Session.Maps()时，获取cache对象报错，返回错误信息： %s",
	En: "When using Session.Maps() in RedisSession, getting cache object is error  . err info: %s "}

/*
   [ERR_SESS_39009]: 在使用redis session 模式时，获取cache对象报错
*/
var ERR_SESS_39009 MessageType = MessageType{Id: "ERR_SESS_39009",
	Cn: "在RedisSession 模式下，调用Session.Clear()时，获取cache对象报错，返回错误信息： %s",
	En: "When using Session.Clear() in RedisSession, getting cache object is error  . err info: %s "}

/*
   [ERR_CORE_39010]: 创建etcd连接时报错
*/
var ERR_CORE_39010 MessageType = MessageType{Id: "ERR_CORE_39010",
	Cn: "创建etcd连接时，报错，返回错误信息： %s",
	En: "Creating etcdConnection is error  . err info: %s "}

/*
   [ERR_CORE_39011]: 创建etcd连接时报错
*/
var ERR_CORE_39011 MessageType = MessageType{Id: "ERR_CORE_39011",
	Cn: "关闭etcd连接时（ec.cli.Close()），报错，返回错误信息： %s",
	En: "Closing etcdConnection is error when Calling ec.cli.Close() . err info: %s "}

/*
   [ERR_CORE_39012]: etcd连接可以已经关闭
*/
var ERR_CORE_39012 MessageType = MessageType{Id: "ERR_CORE_39012",
	Cn: "etcd连接已经关闭，请检查连接。",
	En: "Etcd Connection maybe closed . It's not working! "}

/*
   [ERR_CORE_39013]: 当使用etcd的批量更新时，部分key被重复更新了，报错了。
*/
var ERR_CORE_39013 MessageType = MessageType{Id: "ERR_CORE_39013",
	Cn: "当使用etcd的批量更新时，部分key被重复更新了。具体报错：%s",
	En: "when batch update in etcd , This key is duplicate. error info: %s "}

/*
  [ERR_MSSV_39014]: 微服务的服务端初始化失败，创建对象server.manager是nil
*/
var ERR_MSSV_39014 MessageType = MessageType{Id: "ERR_MSSV_39014",
	Cn: "微服务的服务端初始化失败，创建对象server.manager 是 nil。",
	En: "When init, The server of MicroServive is failed . The object[server.manager] is nil. "}

/*
   [ERR_MSSV_39015]: 微服务的服务端初始化失败，配置文件不存在 或者文件读取异常。
*/
var ERR_MSSV_39015 MessageType = MessageType{Id: "ERR_MSSV_39015",
	Cn: "微服务的服务端初始化失败，配置文件不存在 或者文件读取异常。 %s",
	En: "When init, The server of MicroServive is failed . this configfile is not exist or this format of configfile is error. %s "}

/*
   [ERR_MSSV_39016]: 微服务的客户端初始化失败，配置文件不存在 或者文件读取异常。
*/
var ERR_MSSV_39016 MessageType = MessageType{Id: "ERR_MSSV_39016",
	Cn: "微服务的客户端初始化失败，配置文件不存在 或者文件读取异常。 %s",
	En: "When init, The client of MicroServive is failed . this configfile is not exist or this format of configfile is error. %s "}

/*
   [ERR_CORE_39017]: etcd put方式执行异常。
*/
var ERR_CORE_39017 MessageType = MessageType{Id: "ERR_CORE_39017",
	Cn: "etcd put方式执行异常。报错信息：%s",
	En: "etcd put() is error. error info : %s "}

/*
   [ERR_CORE_39018]: etcd delete方式执行异常。
*/
var ERR_CORE_39018 MessageType = MessageType{Id: "ERR_CORE_39018",
	Cn: "etcd delete方式执行异常。报错信息：%s",
	En: "etcd delete() is error. error info : %s "}

/*
   [ERR_CORE_39019]: etcd Grant方式执行异常。
*/
var ERR_CORE_39019 MessageType = MessageType{Id: "ERR_CORE_39019",
	Cn: "etcd Grant方式执行异常。报错信息：%s",
	En: "etcd Grant() is error. error info : %s "}

/*
   [ERR_CORE_39020]: etcd Watch时，抛出异常，将退出监听。
*/
var ERR_CORE_39020 MessageType = MessageType{Id: "ERR_CORE_39020",
	Cn: "etcd watch时，抛出异常，将退出监听。报错信息：%s",
	En: "etcd Watch() is error.It will be exit. error info : %s "}

/*
   [ERR_MSSV_39021]: 微服务注册group和rpcs时，删除历史group失败
*/
var ERR_MSSV_39021 MessageType = MessageType{Id: "ERR_MSSV_39021",
	Cn: "微服务在注册Group和Rpcs时，删除历史Group时失败，错误信息：%s",
	En: "MicroService register Group and Rpcs . It is failed when delete the old group. error info : %s "}

/*
   [ERR_MSSV_39022]: 微服务在注册Group和Rpcs时，保存Group和Rpcs时失败
*/
var ERR_MSSV_39022 MessageType = MessageType{Id: "ERR_MSSV_39022",
	Cn: "微服务在注册Group和Rpcs时，保存Group和Rpcs时失败，错误信息：%s",
	En: "MicroService register Group and Rpcs . It is failed when save group . error info : %s "}

/*
   [ERR_MSSV_39023]: 微服务注册中心，根据GroupId获取Group信息失败
*/
var ERR_MSSV_39023 MessageType = MessageType{Id: "ERR_MSSV_39023",
	Cn: "微服务注册中心，根据GroupId获取Group信息失败，错误信息：%s",
	En: "In register center of MicroService , It is failed when query group by groupId . error info : %s "}

/*
   [ERR_MSSV_39024]: 更新微服务客户端失败，clients.managerMap[clients.groupidToFilenameMap[groupId]] 是 nil
*/
var ERR_MSSV_39024 MessageType = MessageType{Id: "ERR_MSSV_39024",
	Cn: "更新微服务客户端失败，clients.managerMap[clients.groupidToFilenameMap[groupId]] 是 nil ，groupId=%s, clients.groupidToFilenameMap[groupId]=%s",
	En: "It is failed when update MicroService client , because  clients.managerMap[clients.groupidToFilenameMap[groupId]] is nil . groupId=%s, clients.groupidToFilenameMap[groupId]=%s "}

/*
   [ERR_MSSV_39025]: 微服务服务端发送票据失败
*/
var ERR_MSSV_39025 MessageType = MessageType{Id: "ERR_MSSV_39025",
	Cn: "微服务服务端发送票据失败，错误信息：%s",
	En: "It is failed when the server of MicroService send ticket . error info : %s "}

/*
   [ERR_MSSV_39026]: 微服务的rpc连接池异常，查不到对应的rpc对象
*/
var ERR_MSSV_39026 MessageType = MessageType{Id: "ERR_MSSV_39026",
	Cn: "微服务的rpc连接池异常，根据 addr=%s 查不到可用rpc对象。",
	En: "It is not found rpc object when query addr=%s ."}

/*
   [ERR_MSSV_39027]: rpc连接池异常，计数有问题
*/
var ERR_MSSV_39027 MessageType = MessageType{Id: "ERR_MSSV_39027",
	Cn: "rpc连接池异常，计数有问题，len(fspg.poolMap[groupId])=%d ，fspg.pointerMap[groupId])=%d 可能不一致，或者不大于 currNum=%d。",
	En: "The count of RPC connection pool  has a problem . len(fspg.poolMap[groupId])=%d ，fspg.pointerMap[groupId])=%d , currNum=%d ."}

/*
   [ERR_MSSV_39028]: 微服务连接池中，rpc对象创建失败
*/
var ERR_MSSV_39028 MessageType = MessageType{Id: "ERR_MSSV_39028",
	Cn: "微服务连接池中，rpc对象创建失败。 错误信息是：%s ",
	En: "In a micro service connection pool, the RPC object is created to fail. error info : %s "}

/*
   [ERR_MSSV_39029]: 微服务客户端获取rpc连接失败，无法请求服务端。
*/
var ERR_MSSV_39029 MessageType = MessageType{Id: "ERR_MSSV_39029",
	Cn: "微服务客户端获取rpc连接失败，无法请求服务端。 错误信息是：%s ",
	En: "The MicroService client failed to get the RPC connection because it failed to request the server . error info : %s "}

/*
   [ERR_MSSV_39030]: 当微服务客户端调用upload方法，上传数据时，3个参数信息的长度不一致。 len(ParamKey_files) != len(fullFilePath_files) != renameValue_files
*/
var ERR_MSSV_39030 MessageType = MessageType{Id: "ERR_MSSV_39030",
	Cn: "当微服务客户端调用upload方法，上传数据时，3个参数信息的长度不一致。 len(ParamKey_files) != len(fullFilePath_files) != renameValue_files ",
	En: "When the microservice client calls the upload method and uploads data, the length of the 3 parameter is not same  . len(ParamKey_files) != len(fullFilePath_files) != renameValue_files "}

/*
   [ERR_MSSV_39031]: 当微服务客户端获取rpc连接时，根据地址（ip_port）未找到对应的地址 .
*/
var ERR_MSSV_39031 MessageType = MessageType{Id: "ERR_MSSV_39031",
	Cn: "当微服务客户端获取rpc连接时，根据地址（ip_port）未找到对应的地址 （smg.Addrs[merge_ip_and_port(sc.Ip, sc.Port)] 是 nil ），Ip=%s , port=%d",
	En: "When the microservice client acquires the RPC connection, the RPC address is not found according to the address (ip_port).smg.Addrs[merge_ip_and_port(sc.Ip, sc.Port)] is nil . Ip=%s , port=%d "}

/*
   [ERR_MSSV_39032]: 当微服务客户端获取rpc连接时，根据地址（ip_port）未找到对应的地址 .
*/
var ERR_MSSV_39032 MessageType = MessageType{Id: "ERR_MSSV_39032",
	Cn: "当微服务客户端获取rpc连接时，根据rpcId 未找到rpc对象 ，rpcId=%s ",
	En: "When the microservice client acquires the RPC connection, no RPC object is found according to rpcId. rpcId=%s "}

/*
   [ERR_MSSV_39033]: 当微服务客户端获取rpc连接时，根据groupId 未找到rpc连接池 .
*/
var ERR_MSSV_39033 MessageType = MessageType{Id: "ERR_MSSV_39033",
	Cn: "当微服务客户端获取rpc连接时，根据groupId 未找到rpc连接池 ，groupId=%s ",
	En: "When the microservice client acquires the RPC connection, the RPC connection pool is not found according to groupId. groupId=%s "}

/*
   [ERR_MSSV_39034]: 当微服务客户端创建网络连接时报错 .
*/
var ERR_MSSV_39034 MessageType = MessageType{Id: "ERR_MSSV_39034",
	Cn: "当微服务客户端创建网络连接时报错 . protocol: %s , addr: %s , error info: %s",
	En: "It is failed when the MicroService client creates a network connection. protocol: %s , addr: %s , error info: %s"}

/*
   [ERR_MSSV_39035]: 当微服务客户端关闭客户端连接时，发现 SimpleRpcClient.conn 是 nil .
*/
var ERR_MSSV_39035 MessageType = MessageType{Id: "ERR_MSSV_39035",
	Cn: "当微服务客户端关闭客户端连接时，发现 SimpleRpcClient.conn 是 nil . ",
	En: "When the microservice client closes the client connection, SimpleRpcClient.conn is nil. "}

/*
   [ERR_MSSV_39036]: 当微服务客户端请求服务端时，返回 401 。说明使用票据进行权限认证时，验证失败。
*/
var ERR_MSSV_39036 MessageType = MessageType{Id: "ERR_MSSV_39036",
	Cn: "当微服务客户端请求服务端时，返回 401 。说明使用票据进行权限认证时，验证失败。 ",
	En: "When the MicroService client requests the server, it returns 401. Indicates that validation is failed when the server authenticate the client request . "}

/*
   [ERR_MSSV_39037]: 当微服务客户端请求服务端时，返回 404。
*/
var ERR_MSSV_39037 MessageType = MessageType{Id: "ERR_MSSV_39037",
	Cn: "当微服务客户端请求服务端时，返回 404 。",
	En: "When the MicroService client requests the server, it returns 404 . "}

/*
 [ERR_MSSV_39038]: 当微服务客户端请求服务端时，返回 404。
*/
var ERR_MSSV_39038 MessageType = MessageType{Id: "ERR_MSSV_39038",
	Cn: "当微服务客户端请求服务端时，返回 500 。",
	En: "When the MicroService client requests the server, it returns 500 . "}

/*
 [ERR_MSSV_39039]: 当微服务客户端请求服务端时，返回的 Content-Type 必须是 'application/octet-stream' , 实际不是。
*/
var ERR_MSSV_39039 MessageType = MessageType{Id: "ERR_MSSV_39039",
	Cn: "当微服务客户端请求服务端时，返回的 Content-Type 必须是 'application/octet-stream' ，但是当前值是： %s 。",
	En: "When the MicroService client requests the server, Content-Type must be 'application/octet-stream' . but Content-Type is %s . "}

/*
 [ERR_MSSV_39040]: 当微服务客户端请求服务端时，返回 404。
*/
var ERR_MSSV_39040 MessageType = MessageType{Id: "ERR_MSSV_39040",
	Cn: "当解析微服务返回结果时,ioutil.ReadAll(resp.Body) 报错，错误信息： %s 。",
	En: "When the analysis of the microservice returns the result, it is error when ioutil.ReadAll(resp.Body) . but sc.client.Head is successful. error info: %s"}

/*
 [ERR_MSSV_39041]: 当微服务客户端进行文件上传操作时，上传文件异常，该文件路径不存在 。
*/
var ERR_MSSV_39041 MessageType = MessageType{Id: "ERR_MSSV_39041",
	Cn: "当微服务客户端进行文件上传操作时，上传文件异常，该文件路径不存在 。 AddFile[%d] :  fullFilePath=[%s] ",
	En: "When the MicroService client uploads the file, the path of upload file does not exist. AddFile[%d] :  fullFilePath=[%s] "}

/*
 [ERR_MSSV_39042]: 当微服务客户端进行文件上传操作时，传入参数为空。
*/
var ERR_MSSV_39042 MessageType = MessageType{Id: "ERR_MSSV_39042",
	Cn: "当微服务客户端进行文件上传操作时，传入参数为空。  uo.fileLength=%d , paramKey = %s  ， fullFilePath =%s ",
	En: "When the microservice client uploads the file, the parameter is empty. uo.fileLength=%d , paramKey = %s  ， fullFilePath =%s "}

/*
 [ERR_MSSV_39043]: 当微服务客户端进行文件上传操作时，传入参数为空。
*/
var ERR_MSSV_39043 MessageType = MessageType{Id: "ERR_MSSV_39043",
	Cn: "当微服务客户端进行文件上传操作时，传入参数为空。  uo.fileLength=%d , paramKey = %s  ， paramValue =%s ",
	En: "When the microservice client uploads the file, the parameter is empty. uo.fileLength=%d , paramKey = %s  ， paramValue =%s "}

/*
 [ERR_MSSV_39044]: 当微服务客户端进行文件上传操作时，传入参数为空。
*/
var ERR_MSSV_39044 MessageType = MessageType{Id: "ERR_MSSV_39044",
	Cn: "当微服务客户端进行文件上传操作时，上传已经完成，但是返回报错，错误编码：%d。错误信息：%s ",
	En: "When the MicroService client uploads the file, the upload is completed, but responsedata is error . httpcode=%d , errorInfo = %s "}

/*
 [ERR_WEB0_39045]: 当Web请求处理完成后，使用Forward方式进行页面跳转时，没有找到可以使用的gml模板.
*/
var ERR_WEB0_39045 MessageType = MessageType{Id: "ERR_WEB0_39045",
	Cn: "当Web请求处理完成后，使用Forward方式进行页面跳转时，没有找到可以使用的gml模板，其中 appname=%s , aliasOfPath=%s ",
	En: "When the Web request processing is completed,  gml template can not be found when using forward . appname=%s, aliasOfPath=%s. "}

/*
 [ERR_WEB0_39046]: 当Web请求处理完成后，使用Forward方式进行页面跳转时，报错.
*/
var ERR_WEB0_39046 MessageType = MessageType{Id: "ERR_WEB0_39046",
	Cn: "当Web请求处理完成后，使用Forward方式进行页面跳转时，报错，其中 appname=%s , aliasOfPath=%s , errorinfo=%s ",
	En: "When the Web request processing is completed,  it is error when using forward . appname=%s, aliasOfPath=%s , errorinfo=%s . "}

/*
 [ERR_WEB0_39047]: 当Web请求处理完成后，使用Json方式进行数据返回时，解析Json格式报错 。
*/
var ERR_WEB0_39047 MessageType = MessageType{Id: "ERR_WEB0_39047",
	Cn: "当Web请求处理完成后，使用Json方式进行数据返回时，解析Json格式报错 。 errorinfo=%s ",
	En: "When the Web request processing is completed,  the format of json data is error when using json . errorinfo=%s  "}

/*
 [ERR_WEB0_39048]: 当获取表单数据时，解析报错。
*/
var ERR_WEB0_39048 MessageType = MessageType{Id: "ERR_WEB0_39048",
	Cn: "当获取表单数据时，解析报错。 错误信息：%s ",
	En: "When the form data is obtained, the formdata is error . errorinfo=%s  "}

/*
 [ERR_WEB0_39049]: 当获取Multipart表单数据时，根据文件名称获取上传文件数据报错。
*/
var ERR_WEB0_39049 MessageType = MessageType{Id: "ERR_WEB0_39049",
	Cn: "当获取Multipart表单数据时，根据文件名称获取上传文件数据报错。 错误信息：%s ",
	En: "When the Multipart form data is obtained, the upload file is error according to the file name . errorinfo=%s  "}

/*
 [ERR_WEB0_39050]: 当获取Multipart表单数据时，创建本地文件时出错。
*/
var ERR_WEB0_39050 MessageType = MessageType{Id: "ERR_WEB0_39050",
	Cn: "当获取Multipart表单数据时，创建本地文件时出错。storeFolder=%s , filename=%s , 错误信息：%s ",
	En: "When the Multipart form data is obtained, creating the local file is error . storeFolder=%s , filename=%s , errorinfo=%s  "}

/*
 [ERR_WEB0_39051]: 当获取Multipart表单数据时，数据拷贝到本地文件时出错。
*/
var ERR_WEB0_39051 MessageType = MessageType{Id: "ERR_WEB0_39051",
	Cn: "当获取Multipart表单数据时，数据拷贝到本地文件时出错。storeFolder=%s , filename=%s , 错误信息：%s ",
	En: "When the Multipart form data is obtained, it is error when data is copied to local file . storeFolder=%s , filename=%s , errorinfo=%s  "}

/*
 [ERR_WEB0_39052]: 在获取上传文件路径时，路径地址报错。
*/
var ERR_WEB0_39052 MessageType = MessageType{Id: "ERR_WEB0_39052",
	Cn: "在获取上传文件路径时，路径地址报错。storeFolder=%s , 错误信息：%s ",
	En: "The path address is wrong when getting the upload file path . storeFolder=%s , errorinfo=%s  "}

/*
 [ERR_WEB0_39053]: 在获取上传文件路径时，该路径不是文件夹。
*/
var ERR_WEB0_39053 MessageType = MessageType{Id: "ERR_WEB0_39053",
	Cn: "在获取上传文件路径时，该路径不是文件夹。storeFolder=%s ",
	En: "The path is not folder when getting the upload file path . storeFolder=%s  "}

/*
 [ERR_WEB0_39054]: 在gml模板中，执行 IncludeHTML/IncludeText  获取页面报错（404/500）。
*/
var ERR_WEB0_39054 MessageType = MessageType{Id: "ERR_WEB0_39054",
	Cn: "在gml模板中，执行 IncludeHTML/IncludeText 获取页面报错（404/500）。urlStr=%s , httpcode=%d , path=%s ",
	En: "When executing IncludeHTML/IncludeText in gml template , it is error . urlStr=%s , httpcode=%d , path=%s "}

/*
 [ERR_WEB0_39055]: 在gml模板中，执行 IncludeHTML/IncludeText 时 没有找到 servHTTP 方法。
*/
var ERR_WEB0_39055 MessageType = MessageType{Id: "ERR_WEB0_39055",
	Cn: "在gml模板中，执行 IncludeHTML/IncludeText 时 没有找到 servHTTP 方法。urlStr=%s , path=%s ",
	En: "When executing IncludeHTML/IncludeText in gml template , the method of servHTTP is not found. urlStr=%s , path=%s "}

/*
 [ERR_WEB0_39056]: 在gml模板中，执行 IncludeHTML/IncludeText 时 报错。
*/
var ERR_WEB0_39056 MessageType = MessageType{Id: "ERR_WEB0_39056",
	Cn: "在gml模板中，执行 IncludeHTML/IncludeText 时 报错。urlStr=%s ",
	En: "When executing IncludeHTML/IncludeText in gml template , it is error .  urlStr=%s "}

/*
 [ERR_WEB0_39057]: 在获取404/500/index页面时报错，文件不存在。
*/
var ERR_WEB0_39057 MessageType = MessageType{Id: "ERR_WEB0_39057",
	Cn: "在获取%s页面时报错，文件不存在。path=%s , errorinfo=%s",
	En: "In obtaining the %s page times error, the file does not exist.  path=%s , errorinfo=%s"}

/*
 [ERR_WEB0_39058]: 装载模板时，文件不存在。
*/
var ERR_WEB0_39058 MessageType = MessageType{Id: "ERR_WEB0_39058",
	Cn: "装载模板时，文件不存在。path=%s , errorinfo=%s",
	En: "When loading template, the file does not exist.  path=%s , errorinfo=%s"}

/*
 [ERR_WEB0_39059]: 装载模板时，初始化出错。
*/
var ERR_WEB0_39059 MessageType = MessageType{Id: "ERR_WEB0_39059",
	Cn: "装载模板时，初始化出错。path=%s , errorinfo=%s",
	En: "When loading template, it is error of initialization.  path=%s , errorinfo=%s"}

/*
 [ERR_WEB0_39060]: 文件变化的监听服务启动异常。
*/
var ERR_WEB0_39060 MessageType = MessageType{Id: "ERR_WEB0_39060",
	Cn: "文件变化的监听服务启动异常。 errorinfo=%s",
	En: "A file change monitor service starts an exception. errorinfo=%s"}

/*
 [ERR_WEB0_39061]: 文件变化的监听服务运行中监控到报错。
*/
var ERR_WEB0_39061 MessageType = MessageType{Id: "ERR_WEB0_39061",
	Cn: "文件变化的监听服务运行中监控到报错。 errorinfo=%s",
	En: "A file change monitor service found an error. errorinfo=%s"}

/*
 [ERR_WEB0_39062]: 文件变化的监听服务运行中监控到报错。
*/
var ERR_WEB0_39062 MessageType = MessageType{Id: "ERR_WEB0_39062",
	Cn: "文件变化的监听服务在添加一个文件夹时报错。 folder=%s , errorinfo=%s",
	En: "A file change monitor service adds a folder to times wrong. folder=%s , errorinfo=%s"}

/*
 [ERR_CMD0_39063]: 在准备生成可执行代码时报错，可能文件不存在。
*/
var ERR_CMD0_39063 MessageType = MessageType{Id: "ERR_CMD0_39063",
	Cn: "在准备生成可执行代码时报错，可能文件不存在。path=%s , errorinfo=%s",
	En: "When preparing to generate executable code times error, maybe the file does not exist. path=%s , errorinfo=%s"}

/*
 [ERR_CMD0_39064]: 进程号文件内容非数字，请查看路径信息。
*/
var ERR_CMD0_39064 MessageType = MessageType{Id: "ERR_CMD0_39064",
	Cn: "进程号文件内容非数字，请查看路径信息。path=%s , errorinfo=%s",
	En: "The content of the processid file is not numeric. Please view the path information. path=%s ,  errorinfo=%s"}

/*
 [ERR_CMD0_39065]: 进程号文件内容非数字，请查看路径信息。
*/
var ERR_CMD0_39065 MessageType = MessageType{Id: "ERR_CMD0_39065",
	Cn: "没有找到进程号文件，路径不存在。path=%s , errorinfo=%s",
	En: "processid file is not found, path does not exist. path=%s ,  errorinfo=%s"}

/*
 [ERR_SQLR_39066]: 创建数据库连接失败。
*/
var ERR_SQLR_39066 MessageType = MessageType{Id: "ERR_SQLR_39066",
	Cn: "创建数据库连接失败。 c.DriverName=%s, c.DataSourceName=%s , error info = %s ",
	En: "It is failed to create a connection of databse. c.DriverName=%s, c.DataSourceName=%s , error info = %s "}

/*
 [ERR_SQLR_39067]: 获取连接池失败，请查看配置是否存在。
*/
var ERR_SQLR_39067 MessageType = MessageType{Id: "ERR_SQLR_39067",
	Cn: "获取连接池[name=%s]失败，请查看配置是否存在。",
	En: "Failed to get connection pool [name=%s]. Please check whether configuration exists."}

/*
 [ERR_SQLR_39068]: 获取连接池失败，请查看配置数据库配置是否正确。
*/
var ERR_SQLR_39068 MessageType = MessageType{Id: "ERR_SQLR_39068",
	Cn: "获取连接池[name=%s]失败，请查看配置数据库配置是否正确。",
	En: "Failed to get connection pool [name=%s]. Please see if the configuration database is configured correctly. "}

/*
 [ERR_SQLR_39069]: 运行 datafinder.StartTrans() 时，报错。
*/
var ERR_SQLR_39069 MessageType = MessageType{Id: "ERR_SQLR_39069",
	Cn: "运行 datafinder.StartTrans() 时，报错。errinfo=%s",
	En: "It is error when call datafinder.StartTrans(). errinfo=%s"}

/*
 [ERR_SQLR_39070]: 运行 datafinder.CommitTrans() 时，报错。
*/
var ERR_SQLR_39070 MessageType = MessageType{Id: "ERR_SQLR_39070",
	Cn: "运行 datafinder.CommitTrans() 时，报错。errinfo=%s",
	En: "It is error when call datafinder.CommitTrans(). errinfo=%s"}

/*
 [ERR_SQLR_39071]: 运行 datafinder.RollbackTrans() 时，报错。
*/
var ERR_SQLR_39071 MessageType = MessageType{Id: "ERR_SQLR_39071",
	Cn: "运行 datafinder.RollbackTrans() 时，报错。errinfo=%s",
	En: "It is error when call datafinder.RollbackTrans(). errinfo=%s"}

/*
 [ERR_SQLR_39072]: ResultList 必须是 一个指针，类型是 slice。
*/
var ERR_SQLR_39072 MessageType = MessageType{Id: "ERR_SQLR_39072",
	Cn: "ResultList 必须是 一个指针，类型是 slice。",
	En: "ResultList must be a pointer and slice . ex: &userList "}

/*
 [ERR_SQLR_39073]: resultList的结果集必须是指针，当前不是，请修正 .
*/
var ERR_SQLR_39073 MessageType = MessageType{Id: "ERR_SQLR_39073",
	Cn: "ResultList 必须是指针，当前不是，请修正 .  >>> appName: %s \n >>> sqlname: %s ",
	En: "ResultList must be a pointer and slice . ex: &userList  >>> appName: %s \n >>> sqlname: %s "}

/*
 [ERR_SQLR_39074]: resultList的结果集必须是指针，当前不是，请修正 .
*/
var ERR_SQLR_39074 MessageType = MessageType{Id: "ERR_SQLR_39074",
	Cn: "数据库执行异常，错误信息：%s \n >>> appName: %s \n >>> sqlname: %s \n >>> sql:    %s \n >>> params:  %v ",
	En: "exec sql is error, errorinfo: %s \n >>> appName: %s \n >>> sqlname: %s \n >>> sql:    %s \n >>> params:  %v "}

/*
 [ERR_SQLR_39075]: sql查询结果后，字段转换异常 .
*/
var ERR_SQLR_39075 MessageType = MessageType{Id: "ERR_SQLR_39075",
	Cn: "sql查询结果后，字段转换异常，从string(%s) 转换到 %s 报错，错误信息：%s \n >>> appName: %s \n >>> sqlname: %s \n >>> sql:    %s \n >>> params:  %v ",
	En: "after query ,converting string(%s) to %s is error, errorinfo: %s \n >>> appName: %s \n >>> sqlname: %s \n >>> sql:    %s \n >>> params:  %v "}

/*
 [ERR_SQLR_39076]: resultMap的结果集必须是Map，当前不是，需要修正传参 .
*/
var ERR_SQLR_39076 MessageType = MessageType{Id: "ERR_SQLR_39076",
	Cn: "resultMap的结果集必须是Map，当前不是，需要修正传参 . \n >>> appName: %s \n >>> sqlname: %s  ",
	En: "resultMap must be a map . \n >>> appName: %s \n >>> sqlname: %s "}

/*
 [ERR_SQLR_39077]: resultMap的value数据不是struct，但是定义的 valuecolumn参数是空值，当前需要修正，valuecolumn不能为空 .
*/
var ERR_SQLR_39077 MessageType = MessageType{Id: "ERR_SQLR_39077",
	Cn: "resultMap的value数据不是struct，但是定义的 valuecolumn参数是空值，当前需要修正，valuecolumn不能为空 . \n >>> appName: %s \n >>> sqlname: %s ",
	En: "if The value of resultMap is not a struct , valuecolumn must be a value , not \"\". \n >>> appName: %s \n >>> sqlname: %s "}

/*
 [ERR_SQLR_39078]: resultMap的value数据不是struct，但是定义的 valuecolumn参数是空值，当前需要修正，valuecolumn不能为空 .
*/
var ERR_SQLR_39078 MessageType = MessageType{Id: "ERR_SQLR_39078",
	Cn: "获取连接池对象时，没有对象。非常严重，未知的bug。 ",
	En: "When you get a connection pool object, there is no object.Very serious, unknown bug. "}

var ERR_WEB0_39079 MessageType = MessageType{Id: "ERR_WEB0_39079",
	Cn: "获取Api类型时出错，只能接收struct类型的数据。当前数据类型是: %v",
	En: "This is not the type of Struct When parsing apibean . The current type is %v  "}
