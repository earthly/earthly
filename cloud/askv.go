package cloud

import (
	"context"
	"strings"

	"github.com/earthly/cloud-api/askv"
	"github.com/pkg/errors"
)

func (c *Client) AutoSkipExists(ctx context.Context, org string, hash []byte) (bool, error) {
	response, err := c.askv.Exists(c.withAuth(ctx), &askv.ExistsRequest{
		OrgName: org,
		Hash:    hash,
	})
	if err != nil {
		return false, errors.Wrap(err, "auto-skip existence check failed")
	}
	return response.Exists, nil
}

func (c *Client) AutoSkipAdd(ctx context.Context, org, path, target string, hash []byte) error {
	_, err := c.askv.Add(c.withAuth(ctx), &askv.AddRequest{
		OrgName:    org,
		TargetName: target,
		TargetPath: path,
		Hash:       hash,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add auto-skip hash")
	}
	return nil
}

func (c *Client) AutoSkipPrune(ctx context.Context, org, pathPrefix, target string, deep bool) (int, error) {
	target = strings.TrimPrefix(target, "+")

	req := &askv.PruneTargetRequest{
		OrgName:       org,
		TargetPath:    pathPrefix,
		TargetName:    target,
		UsePathPrefix: deep,
	}

	res, err := c.askv.PruneTarget(c.withAuth(ctx), req)
	if err != nil {
		return 0, errors.Wrap(err, "failed to prune auto-skip data")
	}

	return int(res.GetCount()), nil
}
