// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package mock

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func Test_MockWebAction(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		t.Log("forms= ", r.Form)
		t.Log("p2=", r.FormValue("p2"))
		cookie, cookieErr := r.Cookie("alphabet09-session-id")
		if cookieErr == nil {
			t.Logf("cookie: %s = %s", cookie.Name, cookie.Value)
		} else {
			t.Log("cookieErr=", cookieErr)
		}

		io.WriteString(w, "Hello world!")
	}

	header1 := NewMockWebHeader()
	AddCookies(header1, map[string]string{"alphabet09-session-id": "364d2a29787b1eed5b386e6bf51638ad41e42d9a"})

	code, respBody := MockWebAction("post", "hello/doit?s1=gvv1&s2=gvv2", "p1=value1&p2=value2", nil, header1, handler)
	t.Log(code, "   ", respBody)

	u, _ := url.Parse("http://localhost:888/hello/doit?s1=gvv1&s2=gvv2")
	t.Log("u--------->>>>   ", "http://localhost:888/hello/doit?s1=gvv1&s2=gvv2")
	t.Log("u.RequestURI()--------->>>>   ", u.RequestURI())
	t.Log("u.Path--------->>>>   ", u.Path)
	t.Log("u.Host--------->>>>   ", u.Host)

	u2, _ := url.Parse("/gift/doit?s1=gvv1&s2=gvv2")
	t.Log("u2--------->>>>   ", "gift/doit?s1=gvv1&s2=gvv2")
	t.Log("u2.RequestURI()--------->>>>   ", u2.RequestURI())
	t.Log("u2.Path--------->>>>   ", u2.Path)
	t.Log("u2.Host--------->>>>   ", u2.Host)

	path := u2.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasPrefix(path, "/web2") {
		path = "/web2" + path
	}
	t.Log("u2.Path--------->>>>   ", path)
}
