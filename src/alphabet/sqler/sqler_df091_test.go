// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sqler

import (
	"alphabet/core/abtest"
	"sync"
	"testing"
	"time"
)

// 压力测试

//go test -v -test.run Test_sqler_df091
func Test_sqler_df091(t *testing.T) {
	test_init_sqler_df000_only_register(t)

	Test_sqler_df001(t)

	flag_create_init = false

	var wg sync.WaitGroup

	num := 10
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			dataFinder, _ := GetConnectionForDataFinder(dbsConfigListToml_sqler_df000.Dbs[0].Name)
			defer dataFinder.ReleaseConnection(0)

			for j := 0; j < 20; j++ {
				test_sqler_df001_nest(dataFinder, t)
				test_sqler_df002_nest(dataFinder, t)
				test_sqler_df003_nest(dataFinder, t)
				test_sqler_df011_nest(dataFinder, t)
				test_sqler_df012_nest(dataFinder, t)
				test_sqler_df013_nest(dataFinder, t)

			}
		}()
	}

	wg.Wait()

	abtest.InfoTestLog(t, "Game Over.")

	time.Sleep(time.Second * 2)

}
