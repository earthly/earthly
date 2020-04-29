package earthfile2llb

import (
	"fmt"
	"path"
	"strings"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb/dedup"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

func parseMounts(mounts []string, target domain.Target, ti dedup.TargetInput, cacheContext llb.State) ([]llb.RunOption, error) {
	var runOpts []llb.RunOption
	for _, mount := range mounts {
		target, st, mountOpts, err := parseMount(mount, target, ti, cacheContext)
		if err != nil {
			return nil, errors.Wrap(err, "parse mount")
		}
		runOpts = append(runOpts, llb.AddMount(target, st, mountOpts...))
	}
	return runOpts, nil
}

func parseMount(mount string, target domain.Target, ti dedup.TargetInput, cacheContext llb.State) (string, llb.State, []llb.MountOption, error) {
	var state llb.State
	var mountTarget string
	var mountID string
	var mountType string
	var mountOpts []llb.MountOption
	sharingMode := llb.CacheMountShared
	kvPairs := strings.Split(mount, ",")
	for _, kvPair := range kvPairs {
		kvSplit := strings.SplitN(kvPair, "=", 2)
		if len(kvSplit) == 0 {
			return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
		}
		switch kvSplit[0] {
		case "id":
			if len(kvSplit) != 2 {
				return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			}
			mountID = kvSplit[1]
		case "type":
			if len(kvSplit) != 2 {
				return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			}
			mountType = kvSplit[1]
		case "target":
			if len(kvSplit) != 2 {
				return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			}
			mountTarget = kvSplit[1]
		case "ro", "readonly":
			if len(kvSplit) != 1 {
				return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			}
			mountOpts = append(mountOpts, llb.Readonly)
		case "uid":
			return "", llb.State{}, nil, fmt.Errorf("Not yet supported %s", kvPair)
			// if len(kvSplit) != 2 {
			// 	return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			// }
			// var err error
			// uid, err = strconv.ParseInt(kvSplit[1], 10, 64)
			// if err != nil {
			// 	return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			// }
		case "gid":
			return "", llb.State{}, nil, fmt.Errorf("Not yet supported %s", kvPair)
			// if len(kvSplit) != 2 {
			// 	return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			// }
			// var err error
			// gid, err = strconv.ParseInt(kvSplit[1], 10, 64)
			// if err != nil {
			// 	return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			// }
		case "mode":
			return "", llb.State{}, nil, fmt.Errorf("Not yet supported %s", kvPair)
			// if len(kvSplit) != 2 {
			// 	return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			// }
			// var err error
			// var mode64 int64
			// mode64, err = strconv.ParseInt(kvSplit[1], 8, 64)
			// if err != nil {
			// 	return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			// }
			// mode = int(mode64)
		case "sharing":
			if len(kvSplit) != 2 {
				return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			}
			switch kvSplit[1] {
			case "shared":
				sharingMode = llb.CacheMountShared
			case "private":
				sharingMode = llb.CacheMountPrivate
			case "locked":
				sharingMode = llb.CacheMountLocked
			default:
				return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
			}
		case "from":
			return "", llb.State{}, nil, fmt.Errorf("Not yet supported %s", kvPair)
		case "source":
			return "", llb.State{}, nil, fmt.Errorf("Not yet supported %s", kvPair)
		default:
			return "", llb.State{}, nil, fmt.Errorf("Invalid mount arg %s", kvPair)
		}
	}
	if mountType == "" {
		return "", llb.State{}, nil, fmt.Errorf("Mount type not specified")
	}
	if mountTarget == "" {
		return "", llb.State{}, nil, fmt.Errorf("Mount target not specified")
	}

	switch mountType {
	case "cache":
		if mountID == "" {
			mountID = path.Clean(mountTarget)
		}
		key, err := cacheKeyTargetInput(ti)
		if err != nil {
			return "", llb.State{}, nil, err
		}
		cachePath := path.Join("/run/cache", key, mountID)
		mountOpts = append(mountOpts, llb.AsPersistentCacheDir(cachePath, sharingMode))
		state = cacheContext
	default:
		return "", llb.State{}, nil, fmt.Errorf("Invalid mount type %s", mountType)
	}
	return mountTarget, state, mountOpts, nil
}

func cacheKeyTargetInput(ti dedup.TargetInput) (string, error) {
	digest, err := ti.HashNoTag()
	if err != nil {
		return "", errors.Wrapf(err, "compute hash key for %s", ti.TargetCanonical)
	}
	return digest, nil
}
