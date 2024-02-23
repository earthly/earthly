package subcmd

import (
	"context"
	"testing"

	"github.com/earthly/earthly/cloud"
	"github.com/stretchr/testify/require"
)

type testOrgLister struct {
	results []*cloud.OrgDetail
	err     error
}

func (l *testOrgLister) ListOrgs(ctx context.Context) ([]*cloud.OrgDetail, error) {
	return l.results, l.err
}

func Test_getOrgAndProject(t *testing.T) {

	tests := []struct {
		desc                 string
		orgLister            orgLister
		orgFlag, projectFlag string
		path                 string
		wantOrg, wantProject string
		wantPersonal         bool
		errString            string
	}{
		{
			desc:         "org & project provided (not personal)",
			orgFlag:      "my-org",
			projectFlag:  "my-project",
			wantOrg:      "my-org",
			wantProject:  "my-project",
			wantPersonal: false,
			orgLister: &testOrgLister{
				results: []*cloud.OrgDetail{
					{
						Name:     "my-org",
						Personal: false,
					},
				},
			},
		},
		{
			desc:         "org & project provided (org not found)",
			orgFlag:      "my-org",
			projectFlag:  "my-project",
			wantOrg:      "my-org",
			wantProject:  "my-project",
			wantPersonal: false,
			errString:    "not a member",
			orgLister:    &testOrgLister{},
		},
		{
			desc:         "org & project provided (/user/ prefix)",
			orgFlag:      "my-org",
			projectFlag:  "my-project",
			wantOrg:      "personal-org",
			wantProject:  "",
			wantPersonal: true,
			path:         "/user/",
			orgLister: &testOrgLister{
				results: []*cloud.OrgDetail{
					{
						Name:     "personal-org",
						Personal: true,
					},
				},
			},
		},
		{
			desc:         "org & project provided (is personal)",
			orgFlag:      "my-org",
			projectFlag:  "my-project",
			wantOrg:      "my-org",
			wantProject:  "my-project",
			wantPersonal: true,
			orgLister: &testOrgLister{
				results: []*cloud.OrgDetail{
					{
						Name:     "my-org",
						Personal: true,
					},
				},
			},
		},
		{
			desc:         "org provided but not project",
			orgFlag:      "my-org",
			projectFlag:  "",
			wantOrg:      "my-org",
			wantProject:  "",
			wantPersonal: false,
			errString:    "--project flag is required",
			orgLister: &testOrgLister{
				results: []*cloud.OrgDetail{
					{
						Name:     "my-org",
						Personal: false,
					},
				},
			},
		},
		{
			desc:         "no org/project & personal not found",
			orgFlag:      "",
			projectFlag:  "",
			wantOrg:      "",
			wantProject:  "",
			wantPersonal: false,
			errString:    "provide an org",
			orgLister: &testOrgLister{
				results: []*cloud.OrgDetail{
					{
						Name:     "my-org",
						Personal: false,
					},
				},
			},
		},
		{
			desc:         "no org/project & personal found",
			orgFlag:      "",
			projectFlag:  "",
			wantOrg:      "personal-org",
			wantProject:  "",
			wantPersonal: true,
			orgLister: &testOrgLister{
				results: []*cloud.OrgDetail{
					{
						Name:     "personal-org",
						Personal: true,
					},
				},
			},
		},
	}

	for _, test := range tests {
		ctx := context.Background()
		t.Run(test.desc, func(t *testing.T) {
			r := require.New(t)
			org, project, isPersonal, err := getOrgAndProject(ctx, test.orgFlag, test.projectFlag, test.orgLister, test.path)
			r.Equal(test.wantOrg, org, "org does not match")
			r.Equal(test.wantProject, project, "project does not match")
			r.Equal(test.wantPersonal, isPersonal, "personal value does not match")
			if test.errString != "" {
				r.ErrorContains(err, test.errString)
			} else {
				r.NoError(err, "no error expected")
			}
		})
	}
}
