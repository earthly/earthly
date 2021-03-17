package buildcontext

import (
	"context"
	"fmt"
	"strings"

	"github.com/earthly/earthly/domain"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

// ImportResolver is a resolver which also takes into account imports.
type ImportResolver struct {
	resolver *Resolver

	localImports  map[string]string // local name -> import full path
	globalImports map[string]string // local name -> import full path
}

// NewImportResolver creates a new import resolver.
func NewImportResolver(r *Resolver, globalImports map[string]string) *ImportResolver {
	li := make(map[string]string)
	gi := make(map[string]string)
	for k, v := range globalImports {
		gi[k] = v
		li[k] = v
	}
	return &ImportResolver{
		resolver:      r,
		localImports:  li,
		globalImports: gi,
	}
}

// GlobalImports returns the internal map of global imports.
func (ir *ImportResolver) GlobalImports() map[string]string {
	return ir.globalImports
}

// AddImport adds an import to the resolver.
func (ir *ImportResolver) AddImport(importStr string, as string, global bool) error {
	if importStr == "" {
		return errors.New("IMPORTing empty string not supported")
	}
	aTarget := fmt.Sprintf("%s+none", importStr) // form a fictional target for parasing purposes
	parsedImport, err := domain.ParseTarget(aTarget)
	if err != nil {
		return errors.Wrapf(err, "could not parse IMPORT %s", importStr)
	}
	importStr = parsedImport.ProjectCanonical() // normalize
	var path string
	if parsedImport.IsImportReference() {
		return errors.Errorf("IMPORT %s not supported", importStr)
	} else if parsedImport.IsRemote() {
		path = parsedImport.GetGitURL()
	} else if parsedImport.IsLocalExternal() {
		path = parsedImport.GetLocalPath()
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
		return errors.New("IMPORT requires AS if the import path ends with . or ..")
	}
	if as == "" {
		as = defaultAs
	}
	if strings.ContainsAny(as, "/:") {
		return errors.Errorf("invalid IMPORT AS %s", as)
	}

	ir.localImports[as] = importStr
	if global {
		ir.globalImports[as] = importStr
	}
	return nil
}

// DerefImport resolves the import (if any) and returns a reference with the full path.
func (ir *ImportResolver) DerefImport(ref domain.Reference) (domain.Reference, error) {
	if ref.IsImportReference() {
		fullPath, ok := ir.localImports[ref.GetImportRef()]
		if !ok {
			return nil, errors.Errorf("import reference %s could not be resolved", ref.GetImportRef())
		}
		var resolvedRef domain.Reference
		resolvedRefStr := fmt.Sprintf("%s+%s", fullPath, ref.GetName())
		switch ref.(type) {
		case domain.Target:
			ref2, err := domain.ParseTarget(resolvedRefStr)
			if err != nil {
				return nil, err
			}
			resolvedRef = &domain.Target{
				GitURL:    ref2.GitURL,
				Tag:       ref2.Tag,
				LocalPath: ref2.LocalPath,
				Target:    ref2.Target,
				ImportRef: ref.GetImportRef(), // set import ref too
			}
		case domain.Command:
			ref2, err := domain.ParseCommand(resolvedRefStr)
			if err != nil {
				return nil, err
			}
			resolvedRef = &domain.Command{
				GitURL:    ref2.GitURL,
				Tag:       ref2.Tag,
				LocalPath: ref2.LocalPath,
				Command:   ref2.Command,
				ImportRef: ref.GetImportRef(), // set import ref too
			}
		default:
			return nil, errors.New("ref resolve not supported for this type")
		}
		return resolvedRef, nil
	}
	return ref, nil
}

// Resolve is similar to Resolver.Resolve, except that is also takes into account imports.
func (ir *ImportResolver) Resolve(ctx context.Context, gwClient gwclient.Client, ref domain.Reference) (*Data, error) {
	resolvedRef, err := ir.DerefImport(ref)
	if err != nil {
		return nil, err
	}
	return ir.resolver.Resolve(ctx, gwClient, resolvedRef)
}
