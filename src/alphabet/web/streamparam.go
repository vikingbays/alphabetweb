// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/log4go"
	"io"
	"mime/multipart"
)

type StreamParameterHandler struct {
	mr   *multipart.Reader
	init bool
	part *multipart.Part
}

func (sph *StreamParameterHandler) Next() bool {
	if sph.init {
		sph.init = false
		if sph.part == nil {
			return false
			/*
				mp, err := sph.mr.NextPart()
				if err == io.EOF {
					return false
				} else if err != nil {
					log4go.ErrorLog(err)
					return false
				}
				sph.part = mp
			*/
		}
	} else {
		mp, err := sph.mr.NextPart()
		if err == io.EOF {
			return false
		} else if err != nil {
			log4go.ErrorLog(err)
			return false
		}
		sph.part = mp
	}
	if sph.part == nil {
		return false
	} else {
		return true
	}
}

func (sph *StreamParameterHandler) GetKeyName() string {
	if sph.part != nil {
		return sph.part.FormName()
	}
	return ""
}

func (sph *StreamParameterHandler) GetAliasName() string {
	if sph.part != nil {
		return sph.part.FileName()
	}
	return ""
}

type StreamParameter struct {
	Key string
}

func (sph *StreamParameterHandler) Read(bytes []byte) (n int, err error) {
	n, err = sph.part.Read(bytes)
	return
}
