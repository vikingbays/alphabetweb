// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package binary2

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// 定义数据是binary.LittleEndian还是binary.BigEndian
var Binary_Order binary.ByteOrder = binary.LittleEndian

//
// 把数据转换成二进制存储
//

/*
  把int8类型转换成二进制，存储位数：1位。例如： WriteInt8_Binary(7) 输出是：07 （设置 LittleEndian 方式）
*/
func WriteInt8_Binary(int8Data int8) []byte {

	return writeCommonType_Binary(int8Data)

}

/*
  把int16类型转换成二进制，存储位数：2位。例如： WriteInt16_Binary(10) 输出是：0A 00 （设置 LittleEndian 方式）
*/
func WriteInt16_Binary(int16Data int16) []byte {

	return writeCommonType_Binary(int16Data)

}

/*
  把int类型转换成二进制,int类型默认作为int32处理
*/
func WriteInt_Binary(intData int) []byte {

	return WriteInt32_Binary(int32(intData))

}

/*
  把int32类型转换成二进制，存储位数：4位。例如： WriteInt32_Binary(100) 输出是：64 00 00 00 （设置 LittleEndian 方式）
*/
func WriteInt32_Binary(int32Data int32) []byte {

	return writeCommonType_Binary(int32Data)

}

/*
  把int64类型转换成二进制，存储位数：8位。例如： WriteInt64_Binary(1000) 输出是：E8 03 00 00 00 00 00 00 （设置 LittleEndian 方式）
*/
func WriteInt64_Binary(int64Data int64) []byte {

	return writeCommonType_Binary(int64Data)

}

/*
  把float32类型转换成二进制，存储位数：4位。例如：WriteFloat32_Binary(32.889) 输出是：56 8E 03 42 （设置 LittleEndian 方式）
*/
func WriteFloat32_Binary(float32Data float32) []byte {

	return writeCommonType_Binary(float32Data)

}

/*
  把float64类型转换成二进制，存储位数：8位。例如：WriteFloat64_Binary(4454.889) 输出是：25 06 81 95 E3 66 B1 40 （设置 LittleEndian 方式）
*/
func WriteFloat64_Binary(float64Data float64) []byte {

	return writeCommonType_Binary(float64Data)

}

/*
  把string类型转换成二进制，存储位数：实际就是字符串转[]byte的长度。
*/
func WriteString_Binary(str string) []byte {
	return []byte(str)

}

/*
   通用处理方法，用于把int8,int16,int32,int64,float32,float64转换成bytes
*/
func writeCommonType_Binary(commnonData interface{}) []byte {
	w := new(bytes.Buffer)
	err := binary.Write(w, Binary_Order, commnonData)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return nil
	} else {
		return w.Bytes()
	}
}

/*
   打印数据信息
*/
func PrintBinaryDatas(byteDatas []byte) string {
	str := ""
	if byteDatas == nil {
		return ""
	}
	for i, bd := range byteDatas {
		if i == 0 {
			str = fmt.Sprintf("%02X", bd)
		} else {
			str = fmt.Sprintf("%s %02X", str, bd)
		}
	}
	return str
}

/*
   把二进制转换成数据，例如：ReadBinary_Int8([]byte{0x07})，输出：7
*/
func ReadBinary_Int8(commnonByte []byte) int8 {
	var data int8
	ReadBinary_CommonType(commnonByte, &data)
	return data
}

/*
   把二进制转换成数据，例如：ReadBinary_Int16([]byte{0x0A,0x00})，输出：10
*/
func ReadBinary_Int16(commnonByte []byte) int16 {
	var data int16
	ReadBinary_CommonType(commnonByte, &data)
	return data
}

func ReadBinary_Int(commnonByte []byte) int {
	var data int32
	ReadBinary_CommonType(commnonByte, &data)
	return int(data)
}

/*
   把二进制转换成数据，例如：ReadBinary_Int32([]byte{0x64,0x00,0x00,0x00})，输出：100
*/
func ReadBinary_Int32(commnonByte []byte) int32 {
	var data int32
	ReadBinary_CommonType(commnonByte, &data)
	return data
}

/*
   把二进制转换成数据，例如：ReadBinary_Int64([]byte{0xE8 ,0x03 ,0x00 ,0x00 ,0x00 ,0x00 ,0x00 ,0x00})，输出：1000
*/
func ReadBinary_Int64(commnonByte []byte) int64 {
	var data int64
	ReadBinary_CommonType(commnonByte, &data)
	return data
}

/*
   把二进制转换成数据，例如：ReadBinary_Float32([]byte{0x56 ,0x8E ,0x03 ,0x42})，输出：32.889
*/
func ReadBinary_Float32(commnonByte []byte) float32 {
	var data float32
	ReadBinary_CommonType(commnonByte, &data)
	return data
}

/*
   把二进制转换成数据，例如：ReadBinary_Float64([]byte{0x25 ,0x06 ,0x81 ,0x95 ,0xE3 ,0x66 ,0xB1 ,0x40})，输出：4454.889
*/
func ReadBinary_Float64(commnonByte []byte) float64 {
	var data float64
	ReadBinary_CommonType(commnonByte, &data)
	return data
}

/*
   把二进制转换成数据，例如：ReadBinary_String([]byte{0xE6 ,0x82 ,0xA8 ,0xE5 ,0xA5 ,0xBD ,0xEF ,0xBC ,0x8C ,0xE5 ,0x8D ,0x97 ,0xE4 ,0xBA ,0xAC ,0xE3 ,0x80 ,0x82})，输出：您好，南京。
*/
func ReadBinary_String(strByte []byte) string {
	return string(strByte)
}

/*
  通用处理方法，用于把bytes转换成int8,int16,int32,int64,float32,float64
*/
func ReadBinary_CommonType(commnonByte []byte, data interface{}) {
	err := binary.Read(bytes.NewReader(commnonByte), Binary_Order, data)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
}
