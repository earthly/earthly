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

package windows

import (
	"time"

	"golang.org/x/sys/windows"
)

func BootTime() (time.Time, error) {
	bootTime := time.Now().Add(-1 * windows.DurationSinceBoot())

	// According to GetTickCount64, the resolution of the value is limited to
	// the resolution of the system timer, which is typically in the range of
	// 10 milliseconds to 16 milliseconds. So this will round the value to the
	// nearest second to not mislead anyone about the precision of the value
	// and to provide a stable value.
	bootTime = bootTime.Round(time.Second)
	return bootTime, nil
}
