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
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

const kernBoottimeMIB = "kern.boottime"

func BootTime() (time.Time, error) {
	tv, err := unix.SysctlTimeval(kernBoottimeMIB)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get host uptime: %w", err)
	}

	bootTime := time.Unix(int64(tv.Sec), int64(tv.Usec)*int64(time.Microsecond))
	return bootTime, nil
}
