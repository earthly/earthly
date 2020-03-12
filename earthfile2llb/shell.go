package earthfile2llb

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

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

func withShell(args []string) []string {
	return []string{"/bin/sh", "-c", strings.Join(args, " ")}
}

func withShellAndEnvVars(args []string, envVars []string) []string {
	var escapedArgs []string
	for _, arg := range args {
		escapedArgs = append(escapedArgs, escapeShellSingleQuotes(arg))
	}
	return []string{
		"/bin/sh", "-c",
		strings.Join([]string{
			strings.Join(envVars, " "),
			"/bin/sh",
			"-c",
			fmt.Sprintf("'%s'", strings.Join(escapedArgs, " ")),
		}, " "),
	}
}

func withDockerdWrap(args []string, envVars []string) []string {
	var escapedArgs []string
	for _, arg := range args {
		escapedArgs = append(escapedArgs, escapeShellSingleQuotes(arg))
	}
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
			strings.Join([]string{
				strings.Join(envVars, " "),
				"/bin/sh",
				"-c",
				fmt.Sprintf("'%s'", strings.Join(escapedArgs, " ")),
			}, " ") + "\n" +
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
			"rm -f /var/run/docker.sock\n" +
			// TODO: This should not be necessary.
			"rm -rf /var/lib/docker/tmp\n" +
			"rm -rf /var/lib/docker/runtimes\n" +
			"find /tmp/earthly -type f -name '*.sock' -rm\n" +
			// Exit with right code.
			"exit \"\\$exit_code\"\n" +
			"EOF",
	}
}

func escapeShellSingleQuotes(arg string) string {
	return strings.Replace(arg, "'", "'\"'\"'", -1)
}

func parseKeyValue(env string) (string, string) {
	parts := strings.SplitN(env, "=", 2)
	v := ""
	if len(parts) > 1 {
		v = parts[1]
	}

	return parts[0], v
}

func addEnv(envVars []string, key, value string) []string {
	// Note that this mutates the original slice.
	found := false
	for i, envVar := range envVars {
		k, _ := parseKeyValue(envVar)
		if k == key {
			envVars[i] = fmt.Sprintf("%s=%s", key, value)
			found = true
			break
		}
	}
	if !found {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}
	return envVars
}
