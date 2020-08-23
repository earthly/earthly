package earthfile2llb

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

const debuggerPath = "/usr/bin/earth_debugger"

func splitWildcards(name string) (string, string) {
	i := 0
	for ; i < len(name); i++ {
		ch := name[i]
		if ch == '\\' {
			i++
		} else if ch == '*' || ch == '?' || ch == '[' {
			break
		}
	}
	if i == len(name) {
		return name, ""
	}

	base := path.Base(name[:i])
	if name[:i] == "" || strings.HasSuffix(name[:i], string(filepath.Separator)) {
		base = ""
	}
	return path.Dir(name[:i]), base + name[i:]
}

func withShell(args []string, withShell bool) []string {
	if withShell {
		return []string{"/bin/sh", "-c", strings.Join(args, " ")}
	}
	return args
}

func strWithEnvVars(args []string, envVars []string, withShell bool, withDebugger bool) string {
	var cmdParts []string
	cmdParts = append(cmdParts, strings.Join(envVars, " "))
	if withDebugger {
		cmdParts = append(cmdParts, debuggerPath)
	}
	if withShell {
		var escapedArgs []string
		for _, arg := range args {
			escapedArgs = append(escapedArgs, escapeShellSingleQuotes(arg))
		}
		cmdParts = append(cmdParts, "/bin/sh", "-c")
		cmdParts = append(cmdParts, fmt.Sprintf("'%s'", strings.Join(escapedArgs, " ")))
	} else {
		cmdParts = append(cmdParts, args...)
	}
	return strings.Join(cmdParts, " ")
}

type shellWrapFun func(args []string, envVars []string, withShell bool, withDebugger bool) []string

func withShellAndEnvVars(args []string, envVars []string, withShell bool, withDebugger bool) []string {
	return []string{
		"/bin/sh", "-c",
		strWithEnvVars(args, envVars, withShell, withDebugger),
	}
}

func withDockerdWrapOld(args []string, envVars []string, withShell bool, withDebugger bool) []string {
	return []string{
		"/bin/sh", "-c",
		"/bin/sh <<EOF" +
			"#!/bin/sh\n" +
			// Start dockerd.
			// TODO: vfs is extremely inefficient due to lack of CoW capabilities.
			//       Unfortunately, it's the only thing that works for now. Should explore
			//       some more combinations in the future, once buildkitd supports other
			//       storage drivers other than overlayfs.
			"dockerd-entrypoint.sh dockerd -s vfs &>/var/log/docker.log &\n" +
			"dockerd_pid=\"\\$!\"\n" +
			// Wait for dockerd to start up.
			"let i=1\n" +
			"while ! docker ps &>/dev/null ; do\n" +
			"sleep 1\n" +
			"if [ \"\\$i\" -gt \"30\" ] ; then\n" +
			"exit 1\n" +
			"fi\n" +
			"let i+=1\n" +
			"done\n" +
			// Run provided args.
			strWithEnvVars(args, envVars, withShell, withDebugger) + "\n" +
			"exit_code=\"\\$?\"\n" +
			// Shut down dockerd.
			"kill \"\\$dockerd_pid\" &>/dev/null\n" +
			"let i=1\n" +
			"while kill -0 \"\\$dockerd_pid\" &>/dev/null ; do\n" +
			"sleep 1\n" +
			"let i+=1\n" +
			"if [ \"\\$i\" -gt \"10\" ]; then\n" +
			"kill -9 \"\\$dockerd_pid\" &>/dev/null\n" +
			"fi\n" +
			"done\n" +
			// Exit with right code.
			"exit \"\\$exit_code\"\n" +
			"EOF",
	}
}

func escapeShellSingleQuotes(arg string) string {
	return strings.Replace(arg, "'", "'\"'\"'", -1)
}
