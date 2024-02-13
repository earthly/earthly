package domain

import (
	"fmt"
	"strings"

	"github.com/earthly/earthly/conslogging"

	"github.com/pkg/errors"
)

// ImportTrackerVal is used to resolve imports
type ImportTrackerVal struct {
	fullPath        string
	allowPrivileged bool
}

// ImportTracker is a resolver which also takes into account imports.
type ImportTracker struct {
	local   map[string]ImportTrackerVal // local name -> import details
	global  map[string]ImportTrackerVal // local name -> import details
	console conslogging.ConsoleLogger
}

// NewImportTracker creates a new import resolver.
func NewImportTracker(console conslogging.ConsoleLogger, global map[string]ImportTrackerVal) *ImportTracker {
	gi := make(map[string]ImportTrackerVal)
	for k, v := range global {
		gi[k] = v
	}
	return &ImportTracker{
		local:   make(map[string]ImportTrackerVal),
		global:  gi,
		console: console,
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
	aTarget := fmt.Sprintf("%s+none", importStr) // form a fictional target for parasing purposes
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
