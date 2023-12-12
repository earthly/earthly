package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/hint"

	"github.com/pkg/errors"
)

type pathResult int

const ( // iota is reset to 0
	notExist pathResult = iota // c0 == 0
	file                = iota // c1 == 1
	dir                 = iota // c2 == 2
)

// used in unit tests
type pathResultFunc func(path string) pathResult

// ImportTrackerVal is used to resolve imports
type ImportTrackerVal struct {
	fullPath        string
	allowPrivileged bool
}

// ImportTracker is a resolver which also takes into account imports.
type ImportTracker struct {
	local          map[string]ImportTrackerVal // local name -> import details
	global         map[string]ImportTrackerVal // local name -> import details
	console        conslogging.ConsoleLogger
	pathResultFunc pathResultFunc // to help with unit testing
}

// NewImportTracker creates a new import resolver.
func NewImportTracker(console conslogging.ConsoleLogger, global map[string]ImportTrackerVal) *ImportTracker {
	gi := make(map[string]ImportTrackerVal)
	for k, v := range global {
		gi[k] = v
	}
	return &ImportTracker{
		local:          make(map[string]ImportTrackerVal),
		global:         gi,
		console:        console,
		pathResultFunc: getPathResult,
	}
}

// Global returns the internal map of global imports.
func (ir *ImportTracker) Global() map[string]ImportTrackerVal {
	return ir.global
}

// SetGlobal sets the global import map.
func (ir *ImportTracker) SetGlobal(gi map[string]ImportTrackerVal) {
	ir.global = make(map[string]ImportTrackerVal)
	for k, v := range gi {
		ir.global[k] = v
	}
}

// Add adds an import to the resolver.
func (ir *ImportTracker) Add(importStr string, as string, global, currentlyPrivileged, allowPrivilegedFlag bool) error {
	if importStr == "" {
		return errors.New("IMPORTing empty string not supported")
	}
	aTarget := fmt.Sprintf("%s+none", importStr) // form a fictional target for parsing purposes
	parsedImport, err := ParseTarget(aTarget)
	if err != nil {
		return errors.Wrapf(err, "could not parse IMPORT %s", importStr)
	}
	importStr = parsedImport.ProjectCanonical() // normalize
	var path string
	allowPrivileged := currentlyPrivileged
	if parsedImport.IsImportReference() {
		return errors.Errorf("IMPORT %s not supported", importStr)
	} else if parsedImport.IsRemote() {
		path = parsedImport.GetGitURL()
		allowPrivileged = allowPrivileged && allowPrivilegedFlag
	} else if parsedImport.IsLocalExternal() {
		path = parsedImport.GetLocalPath()
		if pathErr := verifyLocalPath(path, ir.pathResultFunc); pathErr != nil {
			return pathErr
		}
		if allowPrivilegedFlag {
			ir.console.Printf("the --allow-privileged flag has no effect when referencing a local target\n")
		}
	} else {
		return errors.Errorf("IMPORT %s not supported", importStr)
	}
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 1 {
		return errors.Errorf("IMPORT %s not supported", importStr)
	}
	defaultAs := pathParts[len(pathParts)-1]
	if defaultAs == "" {
		return errors.Errorf("IMPORT %s not supported", importStr)
	}
	if (defaultAs == "." || defaultAs == "..") && as == "" {
		return errors.New("IMPORT requires AS if the import path ends with \".\" or \"..\"")
	}
	if as == "" {
		as = defaultAs
	}
	if strings.ContainsAny(as, "/:") {
		return errors.Errorf("invalid IMPORT AS %s", as)
	}

	if global {
		_, exists := ir.global[as]
		if exists {
			return errors.Errorf("import ref %s already exists in this scope", as)
		}
		ir.global[as] = ImportTrackerVal{
			fullPath:        importStr,
			allowPrivileged: allowPrivileged,
		}
	} else {
		_, exists := ir.local[as]
		if exists {
			return errors.Errorf("import ref %s already exists in this scope", as)
		}
		ir.local[as] = ImportTrackerVal{
			fullPath:        importStr,
			allowPrivileged: allowPrivileged,
		}
	}
	return nil
}

// Deref resolves the import (if any) and returns a reference with the full path.
func (ir *ImportTracker) Deref(ref Reference) (resolvedRef Reference, allowPrivileged bool, allowPrivilegedSet bool, err error) {
	if ref.IsImportReference() {
		resolvedImport, ok := ir.local[ref.GetImportRef()]
		if !ok {
			resolvedImport, ok = ir.global[ref.GetImportRef()]
			if !ok {
				return nil, false, false, errors.Errorf("import reference %s could not be resolved", ref.GetImportRef())
			}
		}
		var resolvedRef Reference
		resolvedRefStr := fmt.Sprintf("%s+%s", resolvedImport.fullPath, ref.GetName())
		switch ref.(type) {
		case Target:
			ref2, err := ParseTarget(resolvedRefStr)
			if err != nil {
				return nil, false, false, err
			}
			resolvedRef = Target{
				GitURL:    ref2.GitURL,
				Tag:       ref2.Tag,
				LocalPath: ref2.LocalPath,
				Target:    ref2.Target,
				ImportRef: ref.GetImportRef(), // set import ref too
			}
		case Command:
			ref2, err := ParseCommand(resolvedRefStr)
			if err != nil {
				return nil, false, false, err
			}
			resolvedRef = Command{
				GitURL:    ref2.GitURL,
				Tag:       ref2.Tag,
				LocalPath: ref2.LocalPath,
				Command:   ref2.Command,
				ImportRef: ref.GetImportRef(), // set import ref too
			}
		default:
			return nil, false, false, errors.New("ref resolve not supported for this type")
		}
		return resolvedRef, resolvedImport.allowPrivileged, true, nil
	}
	return ref, false, false, nil
}

func verifyLocalPath(path string, pathResultF pathResultFunc) error {
	earthlyFileName := "Earthfile"
	res := pathResultF(path)
	if res == notExist {
		return hint.Wrapf(errors.Errorf("path %q does not exist", path), "Verify the path %q exists", path)
	}
	if res == file {
		if filepath.Base(path) == earthlyFileName {
			return hint.Wrapf(errors.Errorf("path %q is not a directory", path), "Did you mean to import %q?", filepath.Dir(path))
		}
		return hint.Wrap(errors.Errorf("path %q is not a directory", path), "Please use a directory when using a local IMPORT path")
	}
	res = pathResultF(filepath.Join(path, earthlyFileName))
	if res == notExist {
		if filepath.Base(path) == earthlyFileName {
			return hint.Wrapf(errors.Errorf("path %q does not contain an %s", path, earthlyFileName), "The path %q ends with an %q which is a directory.\nDid you mean to create an %q as a file instead?", path, earthlyFileName, earthlyFileName)
		}
		return hint.Wrapf(errors.Errorf("path %q does not contain an %s", path, earthlyFileName), "Verify the path %q contains an %s", path, earthlyFileName)
	}
	if res == dir {
		return hint.Wrapf(errors.Errorf("path %q does contains an %s which is not a file", path, earthlyFileName), "The local IMPORT path %q contains an %q directory and not a file", path, earthlyFileName)
	}
	return nil
}

// getPathResult checks whether the given path does not exist, a directory or a file
func getPathResult(path string) pathResult {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return notExist
	}
	if info.IsDir() {
		return dir
	}
	return file
}
