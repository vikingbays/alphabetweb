// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

/*
 web访问时抛出异常
*/
type WebappError struct {
	Err error
}

func (e *WebappError) Error() string {
	return e.Err.Error()
}
