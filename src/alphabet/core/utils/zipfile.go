// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package utils

import (
	"alphabet/log4go"
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*
* 解压zip文件
* @param srcFile 源文件，zip文件
* @param dest 解压的目标路径
*
* @return 解压后的根目录
 */
func DeCompressZip(srcFile string, dest string) string {
	deRootFolder := ""
	zipFile, err0 := zip.OpenReader(srcFile)
	if err0 != nil {
		log4go.ErrorLog("ZipFile is not exist . err:  ", err0.Error())
	}
	defer zipFile.Close()
	if !ExistFile(dest) {
		err := os.MkdirAll(dest, os.ModePerm)
		if err != nil {
			log4go.ErrorLog("Unzip File Error : " + err.Error())
		}
	}
	for _, innerFile := range zipFile.File {
		info := innerFile.FileInfo()
		if info.IsDir() {
			err := os.MkdirAll(dest+"/"+innerFile.Name, os.ModePerm)
			if err != nil {
				log4go.ErrorLog("Unzip File Error : " + err.Error())
			}
		} else {
			srcFile, err := innerFile.Open()
			if err != nil {
				log4go.ErrorLog("Unzip File Error : " + err.Error())
			}
			defer srcFile.Close()

			relpath := innerFile.FileHeader.Name
			pos0 := strings.Index(relpath, "/")
			if pos0 != -1 {
				if pos0 == 0 {
					pos0 := strings.Index(relpath[1:], "/")
					if pos0 != -1 {
						pos0 = pos0 + 1
						deRootFolder = relpath[1:pos0]
					}
				} else {
					deRootFolder = relpath[0:pos0]
				}
			}
			pos := strings.LastIndex(relpath, "/")
			if pos != -1 {
				if !ExistFile(dest + "/" + relpath[0:pos]) {
					err := os.MkdirAll(dest+"/"+relpath[0:pos], os.ModePerm)
					if err != nil {
						log4go.ErrorLog("Unzip File Error : " + err.Error())
					}
				}
			}

			newFile, err := os.Create(dest + "/" + innerFile.FileHeader.Name)
			if err != nil {
				log4go.ErrorLog("Unzip File Error : " + err.Error())
			}
			defer newFile.Close()
			io.Copy(newFile, srcFile)

		}
	}
	return deRootFolder
}

/*
* 压缩目录成为Zip文件
* @param  srcRoot 源目录
* @param  dest    目标路径，存储的zip文件名
* @param  zipRootFolder  定义zip文件里的根路径
* @param  filterFileNames  被过滤的文件，不会加入到zip文件中
 */
func CompressZip(srcRoot string, dest string, zipRootFolder string, filterFileNames []string) {
	d, _ := os.Create(dest)
	defer d.Close()

	w := zip.NewWriter(d)
	defer w.Close()

	err := filepath.Walk(srcRoot, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			for _, ffname := range filterFileNames {
				if strings.HasPrefix(f.Name(), ffname) {
					return nil
				}
			}
			relPath := relativepath(srcRoot, path, f.Name())
			header := new(zip.FileHeader)
			if relPath == "" {
				header.Name = zipRootFolder + "/" + f.Name()
			} else {
				header.Name = zipRootFolder + "/" + relPath + "/" + f.Name()
			}
			writer, err := w.CreateHeader(header)
			if err != nil {
				log4go.ErrorLog(err)
			}
			file, _ := os.Open(path)
			_, err = io.Copy(writer, file)
			defer file.Close()
			if err != nil {
				log4go.ErrorLog(err)
			}
		}

		return nil
	})
	if err != nil {
		log4go.ErrorLog("filepath.Walk() returned %v\n", err)
	}
}

func relativepath(root string, path string, filename string) string {
	pos0 := len(root) + 1
	pos1 := len(path)
	pos2 := len(filename)
	if pos0 >= pos1 || pos0 >= (pos1-pos2-1) {
		return ""
	} else {
		return path[pos0:(pos1 - pos2 - 1)]
	}
}
