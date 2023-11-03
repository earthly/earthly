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

/*
#cgo LDFLAGS:-lproc
#include <sys/sysctl.h>
#include <mach/mach_time.h>
#include <mach/mach_host.h>
#include <unistd.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func getHostCPULoadInfo() (*cpuUsage, error) {
	var count C.mach_msg_type_number_t = C.HOST_CPU_LOAD_INFO_COUNT
	var cpu cpuUsage
	status := C.host_statistics(C.host_t(C.mach_host_self()),
		C.HOST_CPU_LOAD_INFO,
		C.host_info_t(unsafe.Pointer(&cpu)),
		&count)

	if status != C.KERN_SUCCESS {
		return nil, fmt.Errorf("host_statistics returned status %d", status)
	}

	return &cpu, nil
}

// getClockTicks returns the number of click ticks in one jiffie.
func getClockTicks() int {
	return int(C.sysconf(C._SC_CLK_TCK))
}

func getHostVMInfo64() (*vmStatistics64Data, error) {
	var count C.mach_msg_type_number_t = C.HOST_VM_INFO64_COUNT

	var vmStat vmStatistics64Data
	status := C.host_statistics64(
		C.host_t(C.mach_host_self()),
		C.HOST_VM_INFO64,
		C.host_info_t(unsafe.Pointer(&vmStat)),
		&count)

	if status != C.KERN_SUCCESS {
		return nil, fmt.Errorf("host_statistics64 returned status %d", status)
	}

	return &vmStat, nil
}
