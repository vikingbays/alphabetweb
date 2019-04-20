// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package webutils

import (
	"reflect"
	"testing"
)

var nameDoc string = "name名称"

type Bean1 struct {
	Id         int    `alias:"id0"    doc:"id编码"`   //ssssss
	Name       string `alias:"name0"   doc:"sdsds"` //sasasasa
	Male       bool
	Id32       int32
	Id64       int64
	Fee32      float32
	Fee64      float64
	IdArr      []int
	Id32Arr    []int32
	Id64Arr    []int64
	Fee32Arr   []float32
	Fee64Arr   []float64
	NameArr    []string
	NameArrArr [][]string
	MaleArr    []bool
	IdMap      map[int]bool
	SbeanMap   map[string]SubBean1
	Sbean02Map map[string]*SubBean1
	NameArrMap map[string][]string
	NameMapMap map[string]map[string]string
	Sbean      SubBean1
	Sbean02    *SubBean1
	SbeanArr   []SubBean1
	SbeanArr02 []*SubBean1

	SabeanMap   map[string]SubBean1
	Sabean02Map map[string]*SubBean1
	Sabean      SubBean1
	Sabean02    *SubBean1
	SabeanArr   []SubBean1
	SabeanArr02 []*SubBean1
}

type SubBean1 struct {
	Id         int    `alias:"id0"    doc:"id编码"`   //ssssss
	Name       string `alias:"name0"   doc:"sdsds"` //sasasasa
	NameArr    []string
	Abean      AppleBean1
	Abean02    *AppleBean1
	AbeanMap   map[string]AppleBean1
	AbeanMap02 map[string]*AppleBean1
}

type AppleBean1 struct {
	Id      int    `alias:"aid"    doc:"id编码"`   //ssssss
	Name    string `alias:"aname"   doc:"sdsds"` //sasasasa
	NameArr []string
}

//go test -v -test.run Test_Bean_Encode_reflect
func Test_Bean_Encode_reflect(t *testing.T) {
	// Name: []string{"aaaa", "bbbb"}
	b := &Bean1{Id: 100, Name: "aaaa",
		Male: true,
		Id32: 32, Id64: 64,
		Fee32: 10.32, Fee64: 10.64,
		IdArr: []int{101, 102}, NameArr: []string{"zzzz0001", "zzzz0002"},
		NameArrArr: [][]string{[]string{"name_X1_001", "name_X1_002"}, []string{"name_Y1_Z91", "name_Y1_Z92", "name_Y1_Z93"}},
		MaleArr:    []bool{true, false, true},
		Id32Arr:    []int32{1032, 1132}, Id64Arr: []int64{1064, 1164, 1264, 1364},
		Fee32Arr: []float32{90.32, 92.32}, Fee64Arr: []float64{90.64, 92.64},
		IdMap: map[int]bool{11: true, 22: false},
		SbeanMap: map[string]SubBean1{"a1": SubBean1{Id: 106, NameArr: []string{"a1_106_01", "a1_106_02"}},
			"b1": SubBean1{Id: 108, NameArr: []string{"b1_108_01", "b1_108_02"}}},
		NameArrMap: map[string][]string{"s1": []string{"s1_zz001", "s1_zz002"},
			"t1": []string{"t1_zz001", "t1_zz002"}},
		NameMapMap: map[string]map[string]string{"p1": map[string]string{"m0101": "p1_m0101_aa001", "m0102": "p1_m0102_aa002"},
			"p2": map[string]string{"m0201": "p2_m0201_aa001", "m0202": "p2_m0202_aa002"}},
		Sbean:   SubBean1{Id: 222, NameArr: []string{"11QQ", "22WW"}},
		Sbean02: &SubBean1{Id: 999, NameArr: []string{"ZZAA", "NNYY"}},
		//Sbean02: nil,
		SbeanArr: []SubBean1{SubBean1{Id: 0001, NameArr: []string{"000111QQ", "000122WW"}},
			SubBean1{Id: 0002, NameArr: []string{"000211QQ", "000222WW"}}},
		SbeanArr02: []*SubBean1{&SubBean1{Id: 9001, NameArr: []string{"900111QQ", "900122WW"}},
			&SubBean1{Id: 9002, NameArr: []string{"900211QQ", "900222WW"}}},
	}
	//data := reflect.ValueOf(b)

	flagOfPrintEncodeHttpParams = true
	paramStrings := EncodeHttpParams(reflect.ValueOf(b))
	paramStrings = paramStrings[1:]
	t.Log(paramStrings)

	paramMap := EncodeHttpParamsMap(reflect.ValueOf(b))
	t.Log(paramMap)
}

//go test -v -test.run Test_Bean_Decode_reflect
func Test_Bean_Decode_reflect(t *testing.T) {

	paramBean := &Bean1{}

	paramStringsMap := map[string][]string{
		"id0":                []string{"100"},
		"name0":              []string{"aaaa"},
		"male":               []string{"true"},
		"id32":               []string{"32"},
		"id64":               []string{"64"},
		"fee32":              []string{"99.32"},
		"fee64":              []string{"99.64"},
		"idarr":              []string{"900", "901"},
		"namearr":            []string{"zz001", "zz002", "zz003"},
		"malearr":            []string{"true", "true", "false"},
		"id32arr":            []string{"2032", "2132", "2232"},
		"id64arr":            []string{"2064", "2164"},
		"fee32arr":           []string{"80.32"},
		"fee64arr":           []string{"80.64", "81.64", "82.64"},
		"namearrarr.0":       []string{"aa01", "aa02"},
		"namearrarr.1":       []string{"bb01", "bb02", "bb03"},
		"sbean.id0":          []string{"11001"},
		"sbean.namearr":      []string{"hata001", "hata002"},
		"sbean02.id0":        []string{"22002"},
		"sbean02.namearr":    []string{"tomy001", "tomy002"},
		"sbeanarr.0.id0":     []string{"81001"},
		"sbeanarr.0.namearr": []string{"fish81001", "fish81002"},
		/*"SbeanArr.1.id0":     []string{"82001"},*/
		"sbeanarr.1.namearr":       []string{"shrimp82001", "shrimp82002", "shrimp82003", "shrimp82004"},
		"sbeanarr02.0.id0":         []string{"91001"},
		"sbeanarr02.0.namearr":     []string{"bird91001", "bird91002"},
		"sbeanarr02.1.id0":         []string{"92001"},
		"sbeanarr02.1.namearr":     []string{"plane92001"},
		"sbeanarr02.2.id0":         []string{"93001"},
		"sbeanarr02.2.namearr":     []string{"dart93001", "dart93002", "dart93003"},
		"idmap.11":                 []string{"true"},
		"idmap.12":                 []string{"false"},
		"namearrmap.aa1":           []string{"aa1_v1", "aa1_v2"},
		"namearrmap.ab2":           []string{"ab2_v1", "ab2_v2"},
		"namearrmap.ac3":           []string{"ac3_v1", "ac3_v2", "ac3_v3"},
		"namemapmap.m1.d001":       []string{"d001_v1"},
		"namemapmap.m1.d002":       []string{"d002_v1"},
		"namemapmap.m2.d101":       []string{"d101_v1"},
		"sbeanmap.bean1.id0":       []string{"9001"},
		"sbeanmap.bean1.namearr":   []string{"name01_v1"},
		"sbeanmap.bean2.id0":       []string{"9002"},
		"sbeanmap.bean2.namearr":   []string{"name02_v2"},
		"sbean02map.bean1.id0":     []string{"9021"},
		"sbean02map.bean1.namearr": []string{"name01_v1"},
		"sbean02map.bean2.id0":     []string{"9022"},
		"sbean02map.bean2.namearr": []string{"name02_v2"},
		//"NameArrArr": [][]string{[]string{"aa01", "aa02"}, []string{"bb01", "bb02", "bb03"}},

		"sabean.id0":                []string{"11001"},
		"sabean.name0":              []string{"banana"},
		"sabean.namearr":            []string{"hata001", "hata002"},
		"sabean.abean.aid":          []string{"81"},
		"sabean.abean.aname":        []string{"banana81"},
		"sabean.abean02.aid":        []string{"8102"},
		"sabean.abean02.aname":      []string{"banana8102"},
		"sabean.abeanmap.ss1.aid":   []string{"61"},
		"sabean.abeanmap.ss1.aname": []string{"banana61"},
		"sabean.abeanmap.ss2.aid":   []string{"62"},
		"sabean.abeanmap.ss2.aname": []string{"banana62"},

		"sabean.abeanmap02.tt1.aid":   []string{"6102"},
		"sabean.abeanmap02.tt1.aname": []string{"banana6102"},
		"sabean.abeanmap02.tt2.aid":   []string{"6202"},
		"sabean.abeanmap02.tt2.aname": []string{"banana6202"},
	}
	DecodeHttpParams(paramStringsMap, reflect.ValueOf(paramBean))
	t.Log(paramBean)

	flagOfPrintEncodeHttpParams = true
	t.Log(EncodeHttpParams(reflect.ValueOf(paramBean)))
	//b, _ := json.Marshal(paramBean)
	//t.Log(string(b))

}
