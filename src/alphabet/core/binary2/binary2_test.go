// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package binary2

import (
	"testing"
)

//go test -v -bench="binary2_test.go"
func Test_WriteXxx_Binary(t *testing.T) {

	t.Logf("input: %d , output: %s", 7, PrintBinaryDatas(WriteInt8_Binary(7)))
	t.Logf("input: %d , output: %s", 10, PrintBinaryDatas(WriteInt16_Binary(10)))
	t.Logf("input: %d , output: %s", 100, PrintBinaryDatas(WriteInt_Binary(100)))
	t.Logf("input: %d , output: %s", 100, PrintBinaryDatas(WriteInt32_Binary(100)))
	t.Logf("input: %d , output: %s", 1000, PrintBinaryDatas(WriteInt64_Binary(1000)))
	t.Logf("input: %f , output: %s", 32.889, PrintBinaryDatas(WriteFloat32_Binary(32.889)))
	t.Logf("input: %f , output: %s", 4454.889, PrintBinaryDatas(WriteFloat64_Binary(4454.889)))
	t.Logf("input: %s , output: %s", "您好，南京。", PrintBinaryDatas(WriteString_Binary("您好，南京。")))

	t.Log(PrintBinaryDatas(WriteString_Binary("DAT")))

}

func Test_ReadBinary_Xxx(t *testing.T) {

	t.Logf("input: %d , output: %d", 7, ReadBinary_Int8([]byte{0x07}))
	t.Logf("input: %d , output: %d", 10, ReadBinary_Int16([]byte{0x0A, 0x00}))
	t.Logf("input: %d , output: %d", 100, ReadBinary_Int([]byte{0x64, 0x00, 0x00, 0x00}))
	t.Logf("input: %d , output: %d", 100, ReadBinary_Int32([]byte{0x64, 0x00, 0x00, 0x00}))
	t.Logf("input: %d , output: %d", 1000, ReadBinary_Int64([]byte{0xE8, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
	t.Logf("input: %f , output: %f", 32.889, ReadBinary_Float32([]byte{0x56, 0x8E, 0x03, 0x42}))
	t.Logf("input: %f , output: %f", 4454.889, ReadBinary_Float64([]byte{0x25, 0x06, 0x81, 0x95, 0xE3, 0x66, 0xB1, 0x40}))
	t.Logf("input: %s , output: %s", "您好，南京。", ReadBinary_String([]byte{0xE6, 0x82, 0xA8, 0xE5, 0xA5, 0xBD, 0xEF, 0xBC, 0x8C, 0xE5, 0x8D, 0x97, 0xE4, 0xBA, 0xAC, 0xE3, 0x80, 0x82}))

}
