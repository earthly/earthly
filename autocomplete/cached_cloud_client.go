package autocomplete

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/util/cliutil"
)

type cloudListClient interface {
	ListOrgs(ctx context.Context) ([]*cloud.OrgDetail, error)
	ListProjects(ctx context.Context, orgName string) ([]*cloud.Project, error)
	ListSatellites(ctx context.Context, orgName string) ([]cloud.SatelliteInstance, error)
}

type cachedCloudClient struct {
	c                cloudListClient
	installationName string
}

func NewCachedCloudClient(installationName string, c cloudListClient) cloudListClient {
	if installationName == "" {
		installationName = "earthly"
	}
	return &cachedCloudClient{
		c:                c,
		installationName: installationName,
	}
}

func isStale(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return true
	}
	if time.Since(fi.ModTime()) > time.Second*10 {
		return true
	}
	return false
}

func getJSONPath(installationName, filename string) string {
	if strings.HasPrefix(installationName, "/") {
		// to allow tests to use temp dir instead of ~/.earthly/
		return filepath.Join(installationName, filename)
	}
	return filepath.Join(cliutil.GetEarthlyDir(installationName), filename)
}

var errCacheStale = fmt.Errorf("cache is stale")

func readJSON(installationName, filename string, v any) error {
	path := getJSONPath(installationName, filename)
	if isStale(path) {
		return errCacheStale
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
func saveJSON(installationName, filename string, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	f, err := os.Create(getJSONPath(installationName, filename))
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	_, err = f.Write(b)
	return err
}

func (ccc *cachedCloudClient) ListOrgs(ctx context.Context) ([]*cloud.OrgDetail, error) {
	cached := struct {
		Orgs []string `json:"orgs"`
	}{}
	filename := ".autocomplete.orgs"
	if err := readJSON(ccc.installationName, filename, &cached); err == nil {
		res := []*cloud.OrgDetail{}
		for _, s := range cached.Orgs {
			res = append(res, &cloud.OrgDetail{
				Name: s,
			})
		}
		return res, nil
	}
	orgs, err := ccc.c.ListOrgs(ctx)
	if err != nil {
		return nil, err
	}
	for _, org := range orgs {
		cached.Orgs = append(cached.Orgs, org.Name)
	}
	_ = saveJSON(ccc.installationName, filename, &cached)
	return orgs, nil
}

func (ccc *cachedCloudClient) ListProjects(ctx context.Context, orgName string) ([]*cloud.Project, error) {
	cached := struct {
		Org      string
		Projects []string `json:"project"`
	}{}
	filename := ".autocomplete.projects"
	if err := readJSON(ccc.installationName, filename, &cached); err == nil && cached.Org == orgName {
		res := []*cloud.Project{}
		for _, s := range cached.Projects {
			res = append(res, &cloud.Project{
				Name: s,
			})
		}
		return res, nil
	}
	cached.Projects = nil
	projects, err := ccc.c.ListProjects(ctx, orgName)
	if err != nil {
		return nil, err
	}
	for _, project := range projects {
		cached.Projects = append(cached.Projects, project.Name)
	}
	cached.Org = orgName
	_ = saveJSON(ccc.installationName, filename, &cached)
	return projects, nil
}

func (ccc *cachedCloudClient) ListSatellites(ctx context.Context, orgName string) ([]cloud.SatelliteInstance, error) {
	cached := struct {
		Org        string
		Satellites []string `json:"satellites"`
	}{}
	filename := ".autocomplete.satellites"
	if err := readJSON(ccc.installationName, filename, &cached); err == nil && cached.Org == orgName {
		res := []cloud.SatelliteInstance{}
		for _, s := range cached.Satellites {
			res = append(res, cloud.SatelliteInstance{
				Name: s,
			})
		}
		return res, nil
	}
	cached.Satellites = nil
	satellites, err := ccc.c.ListSatellites(ctx, orgName)
	if err != nil {
		return nil, err
	}
	for _, sat := range satellites {
		cached.Satellites = append(cached.Satellites, sat.Name)
	}
	cached.Org = orgName
	_ = saveJSON(ccc.installationName, filename, &cached)
	return satellites, nil
}
