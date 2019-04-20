// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*
判断文件或者文件夹是否存在

@param fullPathName  文件名

@return bool  返回true表示存在。
*/
func ExistFile(fullPathName string) bool {
	_, err := os.Stat(fullPathName)
	return err == nil || os.IsExist(err)
}

/*
获取文件路径中的具体文件名信息

@param fullPathName  文件路径名 ，例如：  xxx/yyy/abc.zip

@return string  返回文件名信息。  例如：abc.zip

*/
func BaseFileName(fullPathName string) string {
	return filepath.Base(fullPathName)
}

/*
判断是否是文件夹

@param fullPathName  文件夹名

@return bool  返回true表示存在。
*/
func IsFolder(fullPathName string) bool {
	folder, err := os.Stat(fullPathName)
	if err != nil {
		return false
	} else if folder.IsDir() {
		return true
	} else {
		return false
	}
}

// 拷贝文件，从src文件拷贝到dst文件
//
// @param  src  源文件
//
// @param  dst  目标文件
//
func CopyFile(src, dst string) (w int64, err error) {

	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}

/*
字符串截取

@param s  字符串

@param pos  截取的起始位置

@param length  截取的长度

@return string  返回截取的字符串结果。

*/
func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

/*
根据文件或者文件夹字符串信息，获取她的上一级目录。文件夹不能以`/`结尾。默认的文件分割符是`/`。

@param dirctory  文件或文件夹目录字符串

@return string  返回上一级目录的路径。

*/
func GetParentDir(dirctory string) string {
	return Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}
