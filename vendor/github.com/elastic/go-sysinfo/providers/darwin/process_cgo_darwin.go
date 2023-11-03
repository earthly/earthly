// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//go:build (amd64 && cgo) || (arm64 && cgo)
// +build amd64,cgo arm64,cgo

package darwin

// #cgo LDFLAGS:-lproc
// #include <sys/sysctl.h>
// #include <libproc.h>
import "C"

import (
	"errors"
	"unsafe"
)

//go:generate sh -c "go tool cgo -godefs defs_darwin.go > ztypes_darwin.go"

func getProcTaskAllInfo(pid int, info *procTaskAllInfo) error {
	size := C.int(unsafe.Sizeof(*info))
	ptr := unsafe.Pointer(info)

	n, err := C.proc_pidinfo(C.int(pid), C.PROC_PIDTASKALLINFO, 0, ptr, size)
	if err != nil {
		return err
	} else if n != size {
		return errors.New("failed to read process info with proc_pidinfo")
	}

	return nil
}

func getProcVnodePathInfo(pid int, info *procVnodePathInfo) error {
	size := C.int(unsafe.Sizeof(*info))
	ptr := unsafe.Pointer(info)

	n := C.proc_pidinfo(C.int(pid), C.PROC_PIDVNODEPATHINFO, 0, ptr, size)
	if n != size {
		return errors.New("failed to read vnode info with proc_pidinfo")
	}

	return nil
}
