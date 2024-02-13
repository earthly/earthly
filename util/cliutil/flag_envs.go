package cliutil

import (
	"reflect"

	"github.com/urfave/cli/v2"
)

func GetValidEnvNames(app *cli.App) map[string]struct{} {
	envs := map[string]struct{}{}
	for _, envName := range getValidEnvNamesFromCommands(app.Commands) {
		envs[envName] = struct{}{}
	}

	// and root level flags
	for _, flg := range app.Flags {
		for _, envName := range getEnvs(flg) {
			envs[envName] = struct{}{}
		}
	}
	return envs
}

func getValidEnvNamesFromCommands(cmds []*cli.Command) []string {
	envs := []string{}
	for _, cmd := range cmds {
		for _, flg := range cmd.Flags {
			envs = append(envs, getEnvs(flg)...)
		}
	}
	return envs
}

// it's not obvious if urfave supports converting a "flag interface" to a "genericFlag"
func getEnvs(fl cli.Flag) []string {
	fv := reflect.ValueOf(fl)
	for fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return nil
		}
		fv = fv.Elem()
	}
	field := fv.FieldByName("EnvVars")
	if !field.IsValid() {
		return nil
	}

	return field.Interface().([]string)
}
