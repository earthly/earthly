package command

var cmdMap map[Type]string

func init() {
	cmdMap = map[Type]string{
		AddCmd:            Add,
		ArgCmd:            Arg,
		BuildCmd:          Build,
		CacheCmd:          Cache,
		CmdCmd:            Cmd,
		CommandCmd:        Command,
		CopyCmd:           Copy,
		DoCmd:             Do,
		DockerCmd:         Docker,
		EntrypointCmd:     Entrypoint,
		EnvCmd:            Env,
		ExposeCmd:         Expose,
		FromCmd:           From,
		FromDockerfileCmd: FromDockerfile,
		GitCloneCmd:       GitClone,
		HealthcheckCmd:    HealthCheck,
		HostCmd:           Host,
		ImportCmd:         Import,
		LabelCmd:          Label,
		LetCmd:            Let,
		LoadCmd:           Load,
		LocallyCmd:        Locally,
		OnBuildCmd:        OnBuild,
		ProjectCmd:        Project,
		RunCmd:            Run,
		SaveArtifactCmd:   SaveArtifact,
		SaveImageCmd:      SaveImage,
		SetCmd:            Set,
		ShellCmd:          Shell,
		StopSignalCmd:     StopSignal,
		UserCmd:           User,
		VolumeCmd:         Volume,
		WorkdirCmd:        Workdir,
		FunctionCmd:       Function,
	}
}

func CommandToString(cmd Type) string {
	return cmdMap[cmd]
}
