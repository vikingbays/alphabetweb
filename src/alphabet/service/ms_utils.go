// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*
用于定义服务端管理模块的具体实现，以etcd方式实现。
*/

package service

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

/**
 * 获取当前主机的IP地址，
 * @Param ipType ip类型，传递的值是： ipv4 或  ipv6
 * @Return 返回ip地址
 */
func getCurrentIp(ipType string) string {
	if ipType == "" {
		ipType = "ipv4"
	}

	ipString := "" // 获取本机的IP地址
	host, _ := os.Hostname()
	addrsFromHost, _ := net.LookupHost(host)
	for _, a := range addrsFromHost {
		ip0 := net.ParseIP(a)
		if ip0.To4() != nil {
			if ipType == "ipv4" {
				ipString = a
			}
		} else {
			if ipType == "ipv6" {
				ipString = a
			}
		}
	}
	if ipString == "" {
		addrs, _ := net.InterfaceAddrs()
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					if ipType == "ipv4" {
						ipString = ipnet.IP.String()
						break
					}
				} else {
					if ipType == "ipv6" {
						ipString = ipnet.IP.String()
						break
					}
				}
			}
		}
	}
	return ipString

}

func generateTicket() string {
	ticketFromRand := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20; i++ {
		numFromChar := 126 - rand.Intn(93)
		ticketFromRand = fmt.Sprintf("%s%s", ticketFromRand, string(rune(numFromChar)))
	}
	return ticketFromRand
}

func merge_ip_and_port(ip string, port int) string {
	return fmt.Sprintf("%s_%d", ip, port)
}
