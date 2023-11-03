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
	"os"
	"strconv"
	"time"

	"golang.org/x/sys/unix"

	"github.com/elastic/go-sysinfo/types"
)

func (s darwinSystem) Processes() ([]types.Process, error) {
	ps, err := unix.SysctlKinfoProcSlice("kern.proc.all")
	if err != nil {
		return nil, fmt.Errorf("failed to read process table: %w", err)
	}

	processes := make([]types.Process, 0, len(ps))
	for _, kp := range ps {
		pid := kp.Proc.P_pid
		if pid == 0 {
			continue
		}

		processes = append(processes, &process{
			pid: int(pid),
		})
	}

	return processes, nil
}

func (s darwinSystem) Process(pid int) (types.Process, error) {
	p := process{pid: pid}

	return &p, nil
}

func (s darwinSystem) Self() (types.Process, error) {
	return s.Process(os.Getpid())
}

type process struct {
	info *types.ProcessInfo
	pid  int
	cwd  string
	exe  string
	args []string
	env  map[string]string
}

func (p *process) PID() int {
	return p.pid
}

func (p *process) Parent() (types.Process, error) {
	info, err := p.Info()
	if err != nil {
		return nil, err
	}

	return &process{pid: info.PPID}, nil
}

func (p *process) Info() (types.ProcessInfo, error) {
	if p.info != nil {
		return *p.info, nil
	}

	var task procTaskAllInfo
	if err := getProcTaskAllInfo(p.pid, &task); err != nil && err != types.ErrNotImplemented {
		return types.ProcessInfo{}, err
	}

	var vnode procVnodePathInfo
	if err := getProcVnodePathInfo(p.pid, &vnode); err != nil && err != types.ErrNotImplemented {
		return types.ProcessInfo{}, err
	}

	if err := kern_procargs(p.pid, p); err != nil {
		return types.ProcessInfo{}, err
	}

	p.info = &types.ProcessInfo{
		Name: int8SliceToString(task.Pbsd.Pbi_name[:]),
		PID:  p.pid,
		PPID: int(task.Pbsd.Pbi_ppid),
		CWD:  int8SliceToString(vnode.Cdir.Path[:]),
		Exe:  p.exe,
		Args: p.args,
		StartTime: time.Unix(int64(task.Pbsd.Pbi_start_tvsec),
			int64(task.Pbsd.Pbi_start_tvusec)*int64(time.Microsecond)),
	}

	return *p.info, nil
}

func (p *process) User() (types.UserInfo, error) {
	kproc, err := unix.SysctlKinfoProc("kern.proc.pid", p.pid)
	if err != nil {
		return types.UserInfo{}, err
	}

	egid := ""
	if len(kproc.Eproc.Ucred.Groups) > 0 {
		egid = strconv.Itoa(int(kproc.Eproc.Ucred.Groups[0]))
	}

	return types.UserInfo{
		UID:  strconv.Itoa(int(kproc.Eproc.Pcred.P_ruid)),
		EUID: strconv.Itoa(int(kproc.Eproc.Ucred.Uid)),
		SUID: strconv.Itoa(int(kproc.Eproc.Pcred.P_svuid)),
		GID:  strconv.Itoa(int(kproc.Eproc.Pcred.P_rgid)),
		SGID: strconv.Itoa(int(kproc.Eproc.Pcred.P_svgid)),
		EGID: egid,
	}, nil
}

func (p *process) Environment() (map[string]string, error) {
	return p.env, nil
}

func (p *process) CPUTime() (types.CPUTimes, error) {
	var task procTaskAllInfo
	if err := getProcTaskAllInfo(p.pid, &task); err != nil {
		return types.CPUTimes{}, err
	}
	return types.CPUTimes{
		User:   time.Duration(task.Ptinfo.Total_user),
		System: time.Duration(task.Ptinfo.Total_system),
	}, nil
}

func (p *process) Memory() (types.MemoryInfo, error) {
	var task procTaskAllInfo
	if err := getProcTaskAllInfo(p.pid, &task); err != nil {
		return types.MemoryInfo{}, err
	}
	return types.MemoryInfo{
		Virtual:  task.Ptinfo.Virtual_size,
		Resident: task.Ptinfo.Resident_size,
		Metrics: map[string]uint64{
			"page_ins":    uint64(task.Ptinfo.Pageins),
			"page_faults": uint64(task.Ptinfo.Faults),
		},
	}, nil
}

var nullTerminator = []byte{0}

// wrapper around sysctl KERN_PROCARGS2
// callbacks params are optional,
// up to the caller as to which pieces of data they want
func kern_procargs(pid int, p *process) error {
	data, err := unix.SysctlRaw("kern.procargs2", pid)
	if err != nil {
		return nil
	}
	buf := bytes.NewBuffer(data)

	// argc
	var argc int32
	if err := binary.Read(buf, binary.LittleEndian, &argc); err != nil {
		return err
	}

	// exe
	lines := bytes.Split(buf.Bytes(), nullTerminator)
	p.exe = string(lines[0])
	lines = lines[1:]

	// skip nulls
	for len(lines) > 0 {
		if len(lines[0]) == 0 {
			lines = lines[1:]
			continue
		}
		break
	}

	// args
	for i := 0; i < int(argc); i++ {
		p.args = append(p.args, string(lines[0]))
		lines = lines[1:]
	}

	// env vars
	env := make(map[string]string, len(lines))
	for _, l := range lines {
		if len(l) == 0 {
			break
		}

		parts := bytes.SplitN(l, []byte{'='}, 2)
		key := string(parts[0])
		var value string
		if len(parts) == 2 {
			value = string(parts[1])
		}
		env[key] = value
	}
	p.env = env

	return nil
}

func int8SliceToString(s []int8) string {
	buf := bytes.NewBuffer(make([]byte, len(s)))
	buf.Reset()

	for _, b := range s {
		if b == 0 {
			break
		}
		buf.WriteByte(byte(b))
	}
	return buf.String()
}
