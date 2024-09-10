// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import (
	"strconv"
	"strings"

	"github.com/byted-apaas/server-common-go/constants"
)

// PathReplace 路径替换工具
type PathReplace struct {
	path string
}

func NewPathReplace(path string) *PathReplace {
	return &PathReplace{path: path}
}

func (p *PathReplace) Path() string {
	return p.path
}

func (p *PathReplace) Namespace(namespace string) *PathReplace {
	p.path = strings.Replace(p.path, constants.ReplaceNamespace, namespace, 1)
	return p
}

func (p *PathReplace) ObjectAPIName(objectAPIName string) *PathReplace {
	p.path = strings.Replace(p.path, constants.ReplaceObjectAPIName, objectAPIName, 1)
	return p
}

func (p *PathReplace) ObjectAPINameV3(objectAPIName string) *PathReplace {
	p.path = strings.Replace(p.path, constants.ReplaceObjectAPINameV3, objectAPIName, 1)
	return p
}

func (p *PathReplace) FieldAPIName(fieldAPIName string) *PathReplace {
	p.path = strings.Replace(p.path, constants.ReplaceFieldAPIName, fieldAPIName, 1)
	return p
}

func (p *PathReplace) RecordID(recordID int64) *PathReplace {
	p.path = strings.Replace(p.path, constants.ReplaceRecordID, strconv.FormatInt(recordID, 10), 1)
	return p
}

func (p *PathReplace) FileID(fileID string) *PathReplace {
	p.path = strings.Replace(p.path, constants.ReplaceFileID, fileID, 1)
	return p
}

func (p *PathReplace) FunctionAPIName(functionAPIName string) *PathReplace {
	p.path = strings.Replace(p.path, constants.ReplaceFunctionAPIName, functionAPIName, 1)
	return p
}

func (p *PathReplace) ExecutionID(instanceID int64) *PathReplace {
	p.path = strings.Replace(p.path, constants.ReplaceExecutionID, strconv.FormatInt(instanceID, 10), 1)
	return p
}

func (p *PathReplace) APIName(APIName string) *PathReplace {
	p.path = strings.Replace(p.path, constants.ReplaceAPIName, APIName, 1)
	return p
}
