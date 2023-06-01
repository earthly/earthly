package features_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/earthly/earthly/features"
)

func TestFeaturesStringEnabled(t *testing.T) {
	fts := &features.Features{
		Major:              0,
		Minor:              5,
		ReferencedSaveOnly: true,
	}
	s := fts.String()
	Equal(t, "VERSION --referenced-save-only 0.5", s)
}

func TestFeaturesStringDisabled(t *testing.T) {
	fts := &features.Features{
		Major:              1,
		Minor:              1,
		ReferencedSaveOnly: false,
	}
	s := fts.String()
	Equal(t, "VERSION 1.1", s)
}

func TestApplyFlagOverrides(t *testing.T) {
	fts := &features.Features{}
	err := features.ApplyFlagOverrides(fts, "referenced-save-only")
	Nil(t, err)
	Equal(t, true, fts.ReferencedSaveOnly)
	Equal(t, false, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, false, fts.NoImplicitIgnore)
}

func TestApplyFlagOverridesWithDashDashPrefix(t *testing.T) {
	fts := &features.Features{}
	err := features.ApplyFlagOverrides(fts, "--referenced-save-only")
	Nil(t, err)
	Equal(t, true, fts.ReferencedSaveOnly)
	Equal(t, false, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, false, fts.NoImplicitIgnore)
}

func TestApplyFlagOverridesMultipleFlags(t *testing.T) {
	fts := &features.Features{}
	err := features.ApplyFlagOverrides(fts, "referenced-save-only,use-copy-include-patterns,no-implicit-ignore")
	Nil(t, err)
	Equal(t, true, fts.ReferencedSaveOnly)
	Equal(t, true, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, true, fts.NoImplicitIgnore)
}

func TestApplyFlagOverridesEmptyString(t *testing.T) {
	fts := &features.Features{}
	err := features.ApplyFlagOverrides(fts, "")
	Nil(t, err)
	Equal(t, false, fts.ReferencedSaveOnly)
	Equal(t, false, fts.UseCopyIncludePatterns)
	Equal(t, false, fts.ForIn)
	Equal(t, false, fts.RequireForceForUnsafeSaves)
	Equal(t, false, fts.NoImplicitIgnore)
}

func TestAvailableFlags(t *testing.T) {
	// This test feels like it may be overkill, but it's nice to know that if we
	// introduce a typo we have to introduce it twice for our tests to still
	// pass.
	for _, tt := range []struct {
		flag  string
		field string
	}{
		// 0.5
		{"exec-after-parallel", "ExecAfterParallel"},
		{"parallel-load", "ParallelLoad"},
		{"use-registry-for-with-docker", "UseRegistryForWithDocker"},

		// 0.6
		{"for-in", "ForIn"},
		{"no-implicit-ignore", "NoImplicitIgnore"},
		{"referenced-save-only", "ReferencedSaveOnly"},
		{"require-force-for-unsafe-saves", "RequireForceForUnsafeSaves"},
		{"use-copy-include-patterns", "UseCopyIncludePatterns"},

		// 0.7
		{"check-duplicate-images", "CheckDuplicateImages"},
		{"ci-arg", "EarthlyCIArg"},
		{"earthly-git-author-args", "EarthlyGitAuthorArgs"},
		{"earthly-locally-arg", "EarthlyLocallyArg"},
		{"earthly-version-arg", "EarthlyVersionArg"},
		{"explicit-global", "ExplicitGlobal"},
		{"git-commit-author-timestamp", "GitCommitAuthorTimestamp"},
		{"new-platform", "NewPlatform"},
		{"no-tar-build-output", "NoTarBuildOutput"},
		{"save-artifact-keep-own", "SaveArtifactKeepOwn"},
		{"shell-out-anywhere", "ShellOutAnywhere"},
		{"use-cache-command", "UseCacheCommand"},
		{"use-chmod", "UseChmod"},
		{"use-copy-link", "UseCopyLink"},
		{"use-host-command", "UseHostCommand"},
		{"use-no-manifest-list", "UseNoManifestList"},
		{"use-pipelines", "UsePipelines"},
		{"use-project-secrets", "UseProjectSecrets"},
		{"wait-block", "WaitBlock"},

		// unreleased
		{"no-use-registry-for-with-docker", "NoUseRegistryForWithDocker"},
		{"try", "TryFinally"},
		{"no-network", "NoNetwork"},
		{"arg-scope-and-set", "ArgScopeSet"},
		{"earthly-ci-runner-arg", "EarthlyCIRunnerArg"},
		{"use-docker-ignore", "UseDockerIgnore"},
	} {
		tt := tt
		t.Run(tt.flag, func(t *testing.T) {
			t.Parallel()

			var fts features.Features
			err := features.ApplyFlagOverrides(&fts, tt.flag)
			Nil(t, err)
			field := reflect.ValueOf(fts).FieldByName(tt.field)
			True(t, field.IsValid(), "field %v does not exist on %T", tt.field, fts)
			val, ok := field.Interface().(bool)
			True(t, ok, "field %v was not a boolean", tt.field)
			True(t, val, "expected field %v to be set to true by flag %v", tt.field, tt.flag)
		})
	}
}

func TestContext(t *testing.T) {

	fts := &features.Features{}

	t.Run("features can be set and retrieved from context", func(t *testing.T) {
		ctx := context.Background()
		newCtx, err := fts.WithContext(ctx)
		Equal(t, fts, features.FromContext(newCtx))
		NoError(t, err)
	})

	t.Run("context cannot be set more than once", func(t *testing.T) {
		ctx := context.Background()
		ctx2, err := fts.WithContext(ctx)
		NoError(t, err)
		ctx3, err := fts.WithContext(ctx2)
		Error(t, err)
		Equal(t, ctx2, ctx3)
	})

	t.Run("returns nil when not set in context", func(t *testing.T) {
		ctx := context.Background()
		Nil(t, features.FromContext(ctx))
	})
}
