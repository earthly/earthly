package autocomplete

import (
	"context"
	"os"
	"testing"

	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
)

func TestCachedCloudClient(topT *testing.T) {
	type testCtx struct {
		t      *testing.T
		expect expect.Expectation
		dir    string
		mclc   *mockCloudListClient
		ccc    cloudListClient
	}
	o := onpar.BeforeEach(onpar.New(topT), func(t *testing.T) testCtx {
		dir, err := os.MkdirTemp("", "earthlytest*")
		if err != nil {
			t.Fatalf("failed to create tmp dir: %v", err)
		}
		var mclc mockCloudListClient
		return testCtx{
			t:      t,
			expect: expect.New(t),
			dir:    dir,
			mclc:   &mclc,
			ccc:    NewCachedCloudClient(dir, &mclc),
		}
	})
	o.AfterEach(func(tc testCtx) {
		os.RemoveAll(tc.dir)
	})
	defer o.Run()

	o.Spec("caches orgs", func(tc testCtx) {
		for i := 0; i < 3; i++ {
			orgs, err := tc.ccc.ListOrgs(context.Background())
			tc.expect(err).To(not(haveOccurred()))
			tc.expect(len(orgs)).To(equal(2))
			tc.expect(orgs[0].Name).To(equal("abba"))
			tc.expect(orgs[1].Name).To(equal("abc"))
			tc.expect(tc.mclc.listOrgsCallCount).To(equal(1))
		}
	})

	o.Spec("caches projects", func(tc testCtx) {
		for i := 0; i < 3; i++ {
			projects, err := tc.ccc.ListProjects(context.Background(), "abba")
			tc.expect(err).To(not(haveOccurred()))
			tc.expect(len(projects)).To(equal(3))
			tc.expect(projects[0].Name).To(equal("Absolute ABBA"))
			tc.expect(projects[1].Name).To(equal("Arrival"))
			tc.expect(projects[2].Name).To(equal("Ring Ring"))
			tc.expect(tc.mclc.listProjectsCallCount).To(equal(1))
		}
		for i := 0; i < 3; i++ {
			projects, err := tc.ccc.ListProjects(context.Background(), "abc")
			tc.expect(err).To(not(haveOccurred()))
			tc.expect(len(projects)).To(equal(1))
			tc.expect(projects[0].Name).To(equal("def"))
			tc.expect(tc.mclc.listProjectsCallCount).To(equal(2))
		}
	})

	o.Spec("caches satellites", func(tc testCtx) {
		for i := 0; i < 3; i++ {
			projects, err := tc.ccc.ListSatellites(context.Background(), "abba")
			tc.expect(err).To(not(haveOccurred()))
			tc.expect(len(projects)).To(equal(2))
			tc.expect(projects[0].Name).To(equal("sat-one"))
			tc.expect(projects[1].Name).To(equal("sat-two"))
			tc.expect(tc.mclc.listSatellitesCallCount).To(equal(1))
		}
		for i := 0; i < 3; i++ {
			projects, err := tc.ccc.ListSatellites(context.Background(), "abc")
			tc.expect(err).To(not(haveOccurred()))
			tc.expect(len(projects)).To(equal(1))
			tc.expect(projects[0].Name).To(equal("xyz"))
			tc.expect(tc.mclc.listSatellitesCallCount).To(equal(2))
		}
	})

}
