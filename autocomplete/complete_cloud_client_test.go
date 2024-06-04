package autocomplete

import (
	"context"
	"testing"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/conslogging"

	"github.com/urfave/cli/v2"
)

func getAppWithEarthlyFlags() *cli.App {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "org",
		},
		&cli.StringFlag{
			Name: "satellite",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name: "secrets",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "org",
				},
				&cli.StringFlag{
					Name: "project",
				},
			},
		},
	}
	return app
}

type mockCloudListClient struct {
	listOrgsCallCount       int
	listProjectsCallCount   int
	listSatellitesCallCount int
}

func (mclc *mockCloudListClient) ListOrgs(ctx context.Context) ([]*cloud.OrgDetail, error) {
	mclc.listOrgsCallCount += 1
	return []*cloud.OrgDetail{
		{
			Name: "abba",
		},
		{
			Name: "abc",
		},
	}, nil
}

func (mclc *mockCloudListClient) ListProjects(ctx context.Context, orgName string) ([]*cloud.Project, error) {
	mclc.listProjectsCallCount += 1
	if orgName == "abba" {
		return []*cloud.Project{
			{
				Name: "Absolute ABBA",
			},
			{
				Name: "Arrival",
			},
			{
				Name: "Ring Ring",
			},
		}, nil
	}
	if orgName == "abc" {
		return []*cloud.Project{
			{
				Name: "def",
			},
		}, nil
	}
	return []*cloud.Project{}, nil
}

func (mclc *mockCloudListClient) ListSatellites(ctx context.Context, orgName string) ([]cloud.SatelliteInstance, error) {
	mclc.listSatellitesCallCount += 1
	if orgName == "abba" {
		return []cloud.SatelliteInstance{
			{
				Name: "sat-one",
			},
			{
				Name: "sat-two",
			},
		}, nil
	}
	if orgName == "abc" {
		return []cloud.SatelliteInstance{
			{
				Name: "xyz",
			},
		}, nil
	}
	return []cloud.SatelliteInstance{}, nil
}

func getPotentialsWithMockListClient(t *testing.T, cmd string) []string {
	logger := conslogging.Current(conslogging.NoColor, 0, conslogging.Info, false)
	gitLookup := buildcontext.NewGitLookup(logger, "")
	resolver := buildcontext.NewResolver(nil, gitLookup, logger, "", "", "", 0, "")
	mclc := mockCloudListClient{}

	potentials, err := GetPotentials(context.TODO(), resolver, nil, cmd, len(cmd), getAppWithEarthlyFlags(), &mclc)
	NoError(t, err, "GetPotentials failed")
	return potentials
}

func TestOrgCompletion(t *testing.T) {
	potentials := getPotentialsWithMockListClient(t, "earthly secrets --org ab")
	Equal(t, []string{"abba", "abc"}, potentials)
}

func TestProjectCompletionForAbba(t *testing.T) {
	potentials := getPotentialsWithMockListClient(t, "earthly secrets --org abba --project A")
	Equal(t, []string{"Absolute ABBA", "Arrival"}, potentials)
}

func TestProjectCompletionForAbc(t *testing.T) {
	potentials := getPotentialsWithMockListClient(t, "earthly secrets --org abc --project d")
	Equal(t, []string{"def"}, potentials)
}

// TestProjectCompletionWorksWithNoCloudClient tests that autocompletion still works if the cloud client fails to be initialized
// and instead a nil pointer is passed in
func TestProjectCompletionWorksWithNoCloudClient(t *testing.T) {
	logger := conslogging.Current(conslogging.NoColor, 0, conslogging.Info, false)
	gitLookup := buildcontext.NewGitLookup(logger, "")
	resolver := buildcontext.NewResolver(nil, gitLookup, logger, "", "", "", 0, "")

	for _, cmd := range []string{
		"earthly secrets --org mayb",
		"earthly secrets --org maybe --project call-m",
	} {
		potentials, err := GetPotentials(context.TODO(), resolver, nil, cmd, len(cmd), getAppWithEarthlyFlags(), nil)
		NoError(t, err, "GetPotentials failed")
		Equal(t, 0, len(potentials))
	}
}
