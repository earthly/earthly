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
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func MachineID() (string, error) {
	return getMachineGUID()
}

func getMachineGUID() (string, error) {
	const key = registry.LOCAL_MACHINE
	const path = `SOFTWARE\Microsoft\Cryptography`
	const name = "MachineGuid"

	k, err := registry.OpenKey(key, path, registry.READ|registry.WOW64_64KEY)
	if err != nil {
		return "", fmt.Errorf(`failed to open HKLM\%v: %w`, path, err)
	}
	defer k.Close()

	guid, _, err := k.GetStringValue(name)
	if err != nil {
		return "", fmt.Errorf(`failed to get value of HKLM\%v\%v: %w`, path, name, err)
	}

	return guid, nil
}
