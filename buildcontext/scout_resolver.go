package buildcontext

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/syncutil/synccache"
	"github.com/earthly/earthly/util/vertexmeta"

	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	defaultScoutImage  = "docker/scout-cli:1.8.0"
	sbomOutputFileName = "sbom.json"
	spdxOutputFileName = "spdx.json"
)

type scoutResolver struct {
	imageCache *synccache.SyncCache
	console    conslogging.ConsoleLogger
}

type ResolvedScoutImage struct {
	Vulnerabilities *structpb.Struct
	Spdx            *structpb.Struct
}

func (sr *scoutResolver) getFile(ctx context.Context, gwClient gwclient.Client, platr *platutil.Resolver, imageName string, digest string, outputFile string, runOpts ...llb.RunOption) ([]byte, error) {
	opImg := pllb.Image(
		defaultScoutImage, llb.MarkImageInternal, llb.ResolveModePreferLocal,
		llb.Platform(platr.LLBNative()))

	scoutOp := opImg.Run(runOpts...)

	scoutSt := scoutOp.Root()

	scoutRef, err := llbutil.StateToRef(
		ctx, gwClient, scoutSt, false,
		platr.SubResolver(platutil.NativePlatform), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to docker scout for %s", imageName)
	}

	b, err := scoutRef.ReadFile(ctx, gwclient.ReadRequest{
		Filename: outputFile,
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (sr *scoutResolver) strToStructPB(b []byte, imageName string) (*structpb.Struct, error) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(b, &jsonMap)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal docker scout report for image %s", imageName)
	}
	st, err := structpb.NewStruct(jsonMap)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to serialize docker scout report for image %s", imageName)
	}
	return st, nil
}

func (sr *scoutResolver) ResolveImage(ctx context.Context, gwClient gwclient.Client, platr *platutil.Resolver, imageName, digest string, opts ...llb.RunOption) (rsi *ResolvedScoutImage, finalErr error) {
	cacheKey := fmt.Sprintf("%s#%s", "image-digest", digest)
	rsiValue, err := sr.imageCache.Do(ctx, cacheKey, func(ctx context.Context, k interface{}) (interface{}, error) {
		vm := &vertexmeta.VertexMeta{
			TargetName: cacheKey,
			Internal:   true,
		}

		results := make([]*structpb.Struct, 0, 2)
		for format, outputFilePath := range map[string]string{"sbom": sbomOutputFileName, "spdx": spdxOutputFileName} {
			runOpts := []llb.RunOption{
				llb.Args([]string{
					"/docker-scout",
					"cves",
					"--format",
					format,
					"--output",
					outputFilePath,
					imageName,
				}),
				llb.WithCustomNamef("%s DOCKER SCOUT %s", vm.ToVertexPrefix(), digest),
			}

			runOpts = append(runOpts, opts...)

			b, err := sr.getFile(ctx, gwClient, platr, imageName, digest, outputFilePath, runOpts...)
			if err != nil {
				if os.Getenv("IGNORE_SCOUT_ERROR") == "yes" {
					// hack for cases where no scout auth is given and we want to ignore this error
					rgp := &ResolvedScoutImage{}
					return rgp, nil
				}
				return nil, errors.Wrapf(err, "failed to read docker scout report for image %s", imageName)
			}
			st, err := sr.strToStructPB(b, imageName)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to serialize docker scout report for image %s", imageName)
			}
			results = append(results, st)
		}

		rgp := &ResolvedScoutImage{
			Vulnerabilities: results[0],
			Spdx:            results[1],
		}
		return rgp, nil
	})
	if err != nil {
		return nil, err
	}
	rsi = rsiValue.(*ResolvedScoutImage)
	return rsi, nil
}
