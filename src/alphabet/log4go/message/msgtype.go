// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package message

import (
	"alphabet/env"
	"fmt"
)

type MessageType struct {
	Id string // 唯一标示号，消息编码
	Cn string // 中文信息，可以适配fmt.Sprintf 变量
	En string // 英文信息，可以适配fmt.Sprintf 变量
}

func (mt MessageType) String() string {
	if env.Env_Web_I18n == "zh" {
		return fmt.Sprintf("Code(%s),%s", mt.Id, mt.Cn)
	} else {
		return fmt.Sprintf("Code(%s),%s", mt.Id, mt.En)
	}
}
