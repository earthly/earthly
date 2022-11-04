package earthfile2llb

import (
	"path"
	"strconv"
	"strings"

	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

func (c *Converter) parseMounts(mounts []string) ([]llb.RunOption, error) {
	var runOpts []llb.RunOption
	for _, mount := range mounts {
		mountRunOpts, err := c.parseMount(mount)
		if err != nil {
			return nil, errors.Wrap(err, "parse mount")
		}
		runOpts = append(runOpts, mountRunOpts...)
	}
	return runOpts, nil
}

func (c *Converter) parseMount(mount string) ([]llb.RunOption, error) {
	var state pllb.State
	var mountSource string
	var mountTarget string
	var mountID string
	var mountType string
	var mountMode int
	var mountOpts []llb.MountOption
	sharingMode := llb.CacheMountLocked
	kvPairs := strings.Split(mount, ",")
	for _, kvPair := range kvPairs {
		kvSplit := strings.SplitN(kvPair, "=", 2)
		if len(kvSplit) == 0 {
			return nil, errors.Errorf("invalid mount arg %s", kvPair)
		}
		switch kvSplit[0] {
		case "id":
			if len(kvSplit) != 2 {
				return nil, errors.Errorf("invalid mount arg %s", kvPair)
			}
			mountID = kvSplit[1]
		case "type":
			if len(kvSplit) != 2 {
				return nil, errors.Errorf("invalid mount arg %s", kvPair)
			}
			mountType = kvSplit[1]
		case "source":
			if len(kvSplit) != 2 {
				return nil, errors.Errorf("invalid mount arg %s", kvPair)
			}
			mountSource = kvSplit[1]
		case "target":
			if len(kvSplit) != 2 {
				return nil, errors.Errorf("invalid mount arg %s", kvPair)
			}
			mountTarget = kvSplit[1]
		case "ro", "readonly":
			if len(kvSplit) != 1 {
				return nil, errors.Errorf("invalid mount arg %s", kvPair)
			}
			mountOpts = append(mountOpts, llb.Readonly)
		case "uid":
			return nil, errors.Errorf("not yet supported %s", kvPair)
			// if len(kvSplit) != 2 {
			// 	return nil, errors.Errorf("invalid mount arg %s", kvPair)
			// }
			// var err error
			// uid, err = strconv.ParseInt(kvSplit[1], 10, 64)
			// if err != nil {
			// 	return nil, errors.Errorf("invalid mount arg %s", kvPair)
			// }
		case "gid":
			return nil, errors.Errorf("not yet supported %s", kvPair)
			// if len(kvSplit) != 2 {
			// 	return nil, errors.Errorf("invalid mount arg %s", kvPair)
			// }
			// var err error
			// gid, err = strconv.ParseInt(kvSplit[1], 10, 64)
			// if err != nil {
			// 	return nil, errors.Errorf("invalid mount arg %s", kvPair)
			// }
		case "mode":
			if len(kvSplit) != 2 {
				return nil, errors.Errorf("invalid mount arg %s", kvPair)
			}
			var err error
			mountMode, err = parseMode(kvSplit[1])
			if err != nil {
				return nil, errors.Errorf("failed to parse mount mode %s", kvSplit[1])
			}
		case "sharing":
			if len(kvSplit) != 2 {
				return nil, errors.Errorf("invalid mount arg %s", kvPair)
			}
			switch kvSplit[1] {
			case "shared":
				sharingMode = llb.CacheMountShared
			case "private":
				sharingMode = llb.CacheMountPrivate
			case "locked":
				sharingMode = llb.CacheMountLocked
			default:
				return nil, errors.Errorf("invalid mount arg %s", kvPair)
			}
		case "from":
			return nil, errors.Errorf("not yet supported %s", kvPair)
		default:
			return nil, errors.Errorf("invalid mount arg %s", kvPair)
		}
	}
	if mountType == "" {
		return nil, errors.Errorf("mount type not specified")
	}
	if mountID == "" {
		mountID = path.Clean(mountTarget)
	}

	switch mountType {
	case "bind-experimental":
		if mountSource == "" {
			return nil, errors.Errorf("mount source not specified")
		}
		if mountTarget == "" {
			return nil, errors.Errorf("mount target not specified")
		}
		if mountMode != 0 {
			return nil, errors.Errorf("mode is not supported for type=bind-experimental")
		}
		mountOpts = append(mountOpts, llb.HostBind(), llb.SourcePath(mountSource))
		return []llb.RunOption{llb.AddMount(mountTarget, llb.Scratch(), mountOpts...)}, nil
	case "cache":
		if mountTarget == "" {
			return nil, errors.Errorf("mount target not specified")
		}
		if mountMode != 0 {
			return nil, errors.Errorf("mode is not supported for type=cache")
		}
		key, err := cacheKeyTargetInput(c.targetInputActiveOnly())
		if err != nil {
			return nil, err
		}
		cachePath := path.Join("/run/cache", key, mountID)
		mountOpts = append(mountOpts, llb.AsPersistentCacheDir(cachePath, sharingMode))
		state = c.cacheContext
		return []llb.RunOption{pllb.AddMount(mountTarget, state, mountOpts...)}, nil
	case "tmpfs":
		if mountTarget == "" {
			return nil, errors.Errorf("mount target not specified")
		}
		if mountMode != 0 {
			return nil, errors.Errorf("mode is not supported for type=tmpfs")
		}
		state = c.platr.Scratch()
		mountOpts = append(mountOpts, llb.Tmpfs())
		return []llb.RunOption{pllb.AddMount(mountTarget, state, mountOpts...)}, nil
	case "ssh-experimental":
		sshOpts := []llb.SSHOption{llb.SSHID(mountID)}
		if mountTarget != "" {
			sshOpts = append(sshOpts, llb.SSHSocketTarget(mountTarget))
		}
		if mountMode != 0 {
			return nil, errors.Errorf("mode is not supported for type=ssh-experimental")
		}
		return []llb.RunOption{llb.AddSSHSocket(sshOpts...)}, nil
	case "secret":
		if mountTarget == "" {
			return nil, errors.Errorf("mount target not specified")
		}
		if mountMode == 0 {
			// TODO: Perhaps this should just default to the current user automatically from
			//       buildkit side. Then we wouldn't need to open this up to everyone.
			mountMode = 0444
		}

		secretName := strings.TrimPrefix(mountID, "+secrets/")

		secretOpts := []llb.SecretOption{
			llb.SecretID(c.secretID(secretName)),
			llb.SecretFileOpt(0, 0, mountMode),
		}
		return []llb.RunOption{llb.AddSecret(mountTarget, secretOpts...)}, nil
	default:
		return nil, errors.Errorf("invalid mount type %s", mountType)
	}
}

var errInvalidOctal = errors.New("invalid octal")

func parseMode(s string) (int, error) {
	if len(s) == 0 || s[0] != '0' {
		return 0, errInvalidOctal
	}
	mode, err := strconv.ParseInt(s, 8, 64)
	return int(mode), err
}

func cacheKeyTargetInput(ti dedup.TargetInput) (string, error) {
	digest, err := ti.HashNoTag()
	if err != nil {
		return "", errors.Wrapf(err, "compute hash key for %s", ti.TargetCanonical)
	}
	return digest, nil
}
