package buildcontext

import (
	"context"
	"encoding/json"
	"fmt"

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
	defaultScoutImage = "docker/scout-cli:1.8.0"
	outputFileName    = "sbom.json"
)

type scoutResolver struct {
	imageCache *synccache.SyncCache
	console    conslogging.ConsoleLogger
}

type ResolvedScoutImage struct {
	Output *structpb.Struct
}

func (sr *scoutResolver) ResolveImage(ctx context.Context, gwClient gwclient.Client, platr *platutil.Resolver, imageName, digest string, opts ...llb.RunOption) (rsi *ResolvedScoutImage, finalErr error) {
	cacheKey := fmt.Sprintf("%s#%s", "image-digest", digest)
	rsiValue, err := sr.imageCache.Do(ctx, cacheKey, func(ctx context.Context, k interface{}) (interface{}, error) {
		vm := &vertexmeta.VertexMeta{
			TargetName: cacheKey,
			Internal:   true,
		}

		scoutImage := defaultScoutImage

		runOpts := []llb.RunOption{
			llb.Args([]string{
				"/docker-scout",
				"cves",
				"--format",
				"sbom",
				"--output",
				outputFileName,
				imageName,
			}),
			llb.WithCustomNamef("%s DOCKER SCOUT %s", vm.ToVertexPrefix(), digest),
		}

		runOpts = append(runOpts, opts...)
		opImg := pllb.Image(
			scoutImage, llb.MarkImageInternal, llb.ResolveModePreferLocal,
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
			Filename: outputFileName,
		})
		if err != nil {
			//return nil, errors.Wrapf(err, "failed to read docker scout report for image %s", imageName)
			rgp := &ResolvedScoutImage{
				//Output: fmt.Sprintf("failed to do stuff %v", err), // don't commit
			}
			return rgp, nil
		}

		var jsonMap map[string]interface{}
		err = json.Unmarshal(b, &jsonMap)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal docker scout report for image %s", imageName)
		}
		st, err := structpb.NewStruct(jsonMap)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to serialize docker scout report for image %s", imageName)
		}

		rgp := &ResolvedScoutImage{
			Output: st,
		}
		return rgp, nil
	})
	if err != nil {
		return nil, err
	}
	rsi = rsiValue.(*ResolvedScoutImage)
	return rsi, nil
}
