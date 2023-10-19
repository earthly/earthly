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

//go:build amd64 || arm64
// +build amd64 arm64

package darwin

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"golang.org/x/sys/unix"
)

type cpuUsage struct {
	User   uint32
	System uint32
	Idle   uint32
	Nice   uint32
}

func getPageSize() (uint64, error) {
	i, err := unix.SysctlUint32("vm.pagesize")
	if err != nil {
		return 0, fmt.Errorf("vm.pagesize returned %w", err)
	}

	return uint64(i), nil
}

// From sysctl.h - xsw_usage.
type swapUsage struct {
	Total     uint64
	Available uint64
	Used      uint64
	PageSize  uint64
}

const vmSwapUsageMIB = "vm.swapusage"

func getSwapUsage() (*swapUsage, error) {
	var swap swapUsage
	data, err := unix.SysctlRaw(vmSwapUsageMIB)
	if err != nil {
		return nil, err
	}

	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &swap); err != nil {
		return nil, err
	}

	return &swap, nil
}
