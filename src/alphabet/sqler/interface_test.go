// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sqler

import (
	"fmt"
	"reflect"
	"testing"
)

//go test -v -test.run Test_sampe002
func Test_sampe002(t *testing.T) {

}

//go test -v -test.run Test_interface001
func Test_interface001(t *testing.T) {
	s1 := make([]int, 10, 20)
	test1(&s1)
	fmt.Println(s1)

	fmt.Println(len(s1))

	s2 := make(map[string]int)
	test2(s2)
	fmt.Println(s2)
	fmt.Println(s2["aaaaa"])

	s3 := make(map[string]sample1)
	test2(s3)
	fmt.Println(s3)

	type1 := reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(1))
	fmt.Println(type1)
	fmt.Println(type1.Elem().Kind())
	s4 := reflect.MakeMap(type1)
	fmt.Println(s4.Kind())
	fmt.Println(s4.Type().Key().Kind())
	fmt.Println(s4.Type().Elem().Kind())

	s4.SetMapIndex(reflect.ValueOf("ccccc"), reflect.ValueOf(3333))

	fmt.Println(s4)

}

func test2(resultMap interface{}) {

	resultValue := reflect.ValueOf(resultMap)
	resultType := reflect.TypeOf(resultMap)
	fmt.Println(resultType.Kind())
	fmt.Println(resultType.Elem())
	//resultMap["bbbb"] = 222

	if resultType.Kind() == reflect.Map && resultType.Elem().Kind() == reflect.Int {

		fmt.Println(">>>>>>", resultType.Key())
		fmt.Println(">>>>>>", resultType.Elem())

		//resultMap["aaaa"] = 111
		//resultMap["bbbb"] = 222
		resultValue.SetMapIndex(reflect.ValueOf("aaaaa"), reflect.ValueOf(1111))
		resultValue.SetMapIndex(reflect.ValueOf("bbbbb"), reflect.ValueOf(2222))
	} else {
		fmt.Println("resultType.Elem().NumField() = ", resultType.Elem().NumField())
		fmt.Println(resultType.Elem().FieldByName("uname"))
		//resultType.Elem().
	}

}

func test1(result interface{}) {
	resultValue := reflect.ValueOf(result)
	resultType := reflect.TypeOf(result)
	fmt.Println(resultType.Kind())
	fmt.Println(resultType.Elem())
	if resultType.Kind() == reflect.Ptr && resultType.Elem().Kind() == reflect.Slice && resultType.Elem().Elem().Kind() == reflect.Int {

		if resultValue.Elem().CanSet() {
			resultValue.Elem().SetLen(0)
			resultValue.Elem().SetCap(2)
		}

		fmt.Println(resultValue.Elem())
		resultValue.Elem().Set(reflect.Append(resultValue.Elem(), reflect.ValueOf(111111)))

		resultValue.Elem().Set(reflect.Append(resultValue.Elem(), reflect.ValueOf(222222)))

		resultValue.Elem().Set(reflect.Append(resultValue.Elem(), reflect.ValueOf(333333)))

		/*
			r1 := result.(*[]int)
			*r1 = append(*r1, 111)
			//r1 = append(r1, 222)
			fmt.Println(r1)
		*/
	}

}

type sample1 struct {
	uid   int
	uname string
}
