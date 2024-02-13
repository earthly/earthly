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

func strWithEnvVarsAndDocker(args []string, envVars []string, withShell, withDebugger, forceDebugger, withDocker, isExpression bool, exitCodeFile, outputFile string) string {
	var cmdParts []string
	cmdParts = append(cmdParts, strings.Join(envVars, " "))
	if withDocker {
		cmdParts = append(cmdParts, dockerdWrapperPath, "execute")
	}
	if withDebugger {
		cmdParts = append(cmdParts, debuggerPath)

		if forceDebugger {
			cmdParts = append(cmdParts, "--force")
		}
	}
	if withShell {
		var escapedArgs []string
		if outputFile != "" {
			escapedArgs = append(escapedArgs,
				fmt.Sprintf("exec 1<>'\"'\"%s\"'\"' && ", escapeShellSingleQuotes(outputFile)))
		}
		if isExpression {
			escapedArgs = append(escapedArgs, "echo")
		}
		for _, arg := range args {
			escapedArgs = append(escapedArgs, escapeShellSingleQuotes(arg))
		}
		if exitCodeFile != "" {
			escapedArgs = append(escapedArgs,
				fmt.Sprintf("; echo $? >'\"'\"%s\"'\"'", escapeShellSingleQuotes(exitCodeFile)))
		}
		cmdParts = append(cmdParts, "/bin/sh", "-c")
		cmdParts = append(cmdParts, fmt.Sprintf("'%s'", strings.Join(escapedArgs, " ")))
	} else {
		cmdParts = append(cmdParts, args...)
	}
	return strings.Join(cmdParts, " ")
}

type shellWrapFun func(args []string, envVars []string, withShell, withDebugger, forceDebugger bool) []string

func withShellAndEnvVars(args []string, envVars []string, withShell, withDebugger, forceDebugger bool) []string {
	return []string{
		"/bin/sh", "-c",
		strWithEnvVarsAndDocker(args, envVars, withShell, withDebugger, forceDebugger, false, false, "", ""),
	}
}

func withShellAndEnvVarsExitCode(exitCodeFile string) shellWrapFun {
	return func(args []string, envVars []string, withShell, withDebugger, forceDebugger bool) []string {
		if !withShell {
			panic("unexpected exec mode")
		}
		return []string{
			"/bin/sh", "-c",
			strWithEnvVarsAndDocker(args, envVars, true, withDebugger, false, false, false, exitCodeFile, ""),
		}
	}
}

func withShellAndEnvVarsOutput(outputFile string) shellWrapFun {
	return func(args []string, envVars []string, withShell, withDebugger, forceDebugger bool) []string {
		if !withShell {
			panic("unexpected exec mode")
		}
		return []string{
			"/bin/sh", "-c",
			strWithEnvVarsAndDocker(args, envVars, true, withDebugger, false, false, false, "", outputFile),
		}
	}
}

func expressionWithShellAndEnvVarsOutput(outputFile string) shellWrapFun {
	return func(args []string, envVars []string, withShell, withDebugger, forceDebugger bool) []string {
		if !withShell {
			panic("unexpected exec mode")
		}
		return []string{
			"/bin/sh", "-c",
			strWithEnvVarsAndDocker(args, envVars, true, withDebugger, false, false, true, "", outputFile),
		}
	}
}

func escapeShellSingleQuotes(arg string) string {
	return strings.ReplaceAll(arg, "'", "'\"'\"'")
}
