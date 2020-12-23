package earthfile2llb

import (
	"fmt"
	"path"
	"strings"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/states/dedup"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

func parseMounts(mounts []string, target domain.Target, ti dedup.TargetInput, cacheContext llb.State) ([]llb.RunOption, error) {
	var runOpts []llb.RunOption
	for _, mount := range mounts {
		mountRunOpts, err := parseMount(mount, target, ti, cacheContext)
		if err != nil {
			return nil, errors.Wrap(err, "parse mount")
		}
		runOpts = append(runOpts, mountRunOpts...)
	}
	return runOpts, nil
}

func parseMount(mount string, target domain.Target, ti dedup.TargetInput, cacheContext llb.State) ([]llb.RunOption, error) {
	var state llb.State
	var mountSource string
	var mountTarget string
	var mountID string
	var mountType string
	var mountOpts []llb.MountOption
	sharingMode := llb.CacheMountShared
	kvPairs := strings.Split(mount, ",")
	for _, kvPair := range kvPairs {
		kvSplit := strings.SplitN(kvPair, "=", 2)
		if len(kvSplit) == 0 {
			return nil, fmt.Errorf("invalid mount arg %s", kvPair)
		}
		switch kvSplit[0] {
		case "id":
			if len(kvSplit) != 2 {
				return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			}
			mountID = kvSplit[1]
		case "type":
			if len(kvSplit) != 2 {
				return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			}
			mountType = kvSplit[1]
		case "source":
			if len(kvSplit) != 2 {
				return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			}
			mountSource = kvSplit[1]
		case "target":
			if len(kvSplit) != 2 {
				return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			}
			mountTarget = kvSplit[1]
		case "ro", "readonly":
			if len(kvSplit) != 1 {
				return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			}
			mountOpts = append(mountOpts, llb.Readonly)
		case "uid":
			return nil, fmt.Errorf("not yet supported %s", kvPair)
			// if len(kvSplit) != 2 {
			// 	return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			// }
			// var err error
			// uid, err = strconv.ParseInt(kvSplit[1], 10, 64)
			// if err != nil {
			// 	return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			// }
		case "gid":
			return nil, fmt.Errorf("not yet supported %s", kvPair)
			// if len(kvSplit) != 2 {
			// 	return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			// }
			// var err error
			// gid, err = strconv.ParseInt(kvSplit[1], 10, 64)
			// if err != nil {
			// 	return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			// }
		case "mode":
			return nil, fmt.Errorf("not yet supported %s", kvPair)
			// if len(kvSplit) != 2 {
			// 	return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			// }
			// var err error
			// var mode64 int64
			// mode64, err = strconv.ParseInt(kvSplit[1], 8, 64)
			// if err != nil {
			// 	return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			// }
			// mode = int(mode64)
		case "sharing":
			if len(kvSplit) != 2 {
				return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			}
			switch kvSplit[1] {
			case "shared":
				sharingMode = llb.CacheMountShared
			case "private":
				sharingMode = llb.CacheMountPrivate
			case "locked":
				sharingMode = llb.CacheMountLocked
			default:
				return nil, fmt.Errorf("invalid mount arg %s", kvPair)
			}
		case "from":
			return nil, fmt.Errorf("not yet supported %s", kvPair)
		default:
			return nil, fmt.Errorf("invalid mount arg %s", kvPair)
		}
	}
	if mountType == "" {
		return nil, fmt.Errorf("mount type not specified")
	}
	if mountID == "" {
		mountID = path.Clean(mountTarget)
	}

	switch mountType {
	case "bind-experimental":
		if mountSource == "" {
			return nil, fmt.Errorf("mount source not specified")
		}
		if mountTarget == "" {
			return nil, fmt.Errorf("mount target not specified")
		}
		mountOpts = append(mountOpts, llb.HostBind(), llb.SourcePath(mountSource))
		return []llb.RunOption{llb.AddMount(mountTarget, llb.Scratch(), mountOpts...)}, nil
	case "cache":
		if mountTarget == "" {
			return nil, fmt.Errorf("mount target not specified")
		}
		key, err := cacheKeyTargetInput(ti)
		if err != nil {
			return nil, err
		}
		cachePath := path.Join("/run/cache", key, mountID)
		mountOpts = append(mountOpts, llb.AsPersistentCacheDir(cachePath, sharingMode))
		state = cacheContext
		return []llb.RunOption{llb.AddMount(mountTarget, state, mountOpts...)}, nil
	case "tmpfs":
		if mountTarget == "" {
			return nil, fmt.Errorf("mount target not specified")
		}
		state = llbutil.ScratchWithPlatform()
		mountOpts = append(mountOpts, llb.Tmpfs())
		return []llb.RunOption{llb.AddMount(mountTarget, state, mountOpts...)}, nil
	case "ssh-experimental":
		sshOpts := []llb.SSHOption{llb.SSHID(mountID)}
		if mountTarget != "" {
			sshOpts = append(sshOpts, llb.SSHSocketTarget(mountTarget))
		}
		return []llb.RunOption{llb.AddSSHSocket(sshOpts...)}, nil
	case "secret":
		if mountTarget == "" {
			return nil, fmt.Errorf("mount target not specified")
		}
		secretID := strings.TrimPrefix(mountID, "+secrets/")
		secretOpts := []llb.SecretOption{
			llb.SecretID(secretID),
			// TODO: Perhaps this should just default to the current user automatically from
			//       buildkit side. Then we wouldn't need to open this up to everyone.
			llb.SecretFileOpt(0, 0, 0444),
		}
		return []llb.RunOption{llb.AddSecret(mountTarget, secretOpts...)}, nil
	default:
		return nil, fmt.Errorf("invalid mount type %s", mountType)
	}
}

func cacheKeyTargetInput(ti dedup.TargetInput) (string, error) {
	digest, err := ti.HashNoTag()
	if err != nil {
		return "", errors.Wrapf(err, "compute hash key for %s", ti.TargetCanonical)
	}
	return digest, nil
}
