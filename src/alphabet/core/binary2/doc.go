// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*

实现数据和二进制互转。

支持：string , int , int8 , int16 , int32 , int64 , float32 , float64 转换成 []byte 。其中：int 与int32等同。

数据长度：
    string          // (根据数据大小) 采用UTF-8 数据格式存储
    int             // 4 位
    int8            // 1 位
    int16           // 2 位
    int32           // 4 位
    int64           // 8 位
    float32         // 4 位
    float64         // 8 位


提供三种功能：
	WriteXxx_Binary    //把数据转成二机制。     例如：WriteInt16_Binary(10))                  －>输出： 0x0A, 0x00
	ReadBinary_Xxx     //从二进制转成具体数据。  例如：ReadBinary_Int16([]byte{0x0A, 0x00})    －>输出： 10
	PrintBinaryDatas   //打印二进制数据。       例如：PrintBinaryDatas(WriteInt16_Binary(10))


*/
package binary2
