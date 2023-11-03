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

//go:build (amd64 && !cgo) || (arm64 && !cgo)

package darwin

import (
	"fmt"

	"github.com/elastic/go-sysinfo/types"
)

func getHostCPULoadInfo() (*cpuUsage, error) {
	return nil, fmt.Errorf("host cpu load requires cgo: %w", types.ErrNotImplemented)
}

// getClockTicks returns the number of click ticks in one jiffie.
func getClockTicks() int {
	return 0
}

func getHostVMInfo64() (*vmStatistics64Data, error) {
	return nil, fmt.Errorf("host vm info requires cgo: %w", types.ErrNotImplemented)
}
