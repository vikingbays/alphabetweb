// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/service"
)

/*
  在微服务场景下，验证票据合法性
*/
func authTicketServer_MS(context *Context, url string) bool {
	var flag bool = true
	if service.IsServerRunable_MS(context.RootAppname) {
		if service.IsRpc_MS(url, context.RootAppname) { // 判断是不是rpc服务
			header := context.Request.Header
			t1 := header.Get(service.SERVICE_MS_REQ_HEADER_TICKET_NAME)
			t2 := header.Get(service.SERVICE_MS_REQ_HEADER_TICKET_LAST_NAME)
			flag = service.ValidTicketServer_MS(t1, t2, context.RootAppname)
		}
	}

	if !flag {
		context.Return.Forward401()
	}
	return flag
}
