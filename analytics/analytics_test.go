package analytics

import (
	"testing"

	"github.com/earthly/earthly/domain"
)

func TestGetTarget(t *testing.T) {
	type args struct {
		repo   string
		target domain.Target
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Local target with separate repo",
			args: args{
				repo:   "github.com/foo/bar",
				target: domain.Target{Target: "foo"},
			},
			want: "github.com/foo/bar+foo",
		},
		{
			name: "Local target with separate repo and local path",
			args: args{
				repo:   "github.com/foo/bar",
				target: domain.Target{Target: "foo", LocalPath: "baz"},
			},
			want: "github.com/foo/bar/baz+foo",
		},
		{
			name: "Remote target",
			args: args{
				target: domain.Target{Target: "foo", GitURL: "github.com/foo/bar"},
			},
			want: "github.com/foo/bar+foo",
		},
		{
			name: "Remote target with path",
			args: args{
				target: domain.Target{Target: "foo", GitURL: "github.com/foo/bar/baz"},
			},
			want: "github.com/foo/bar/baz+foo",
		},
		{
			name: "Empty target, no repo",
			want: "",
		},
		{
			name: "Unknown repo",
			args: args{
				repo:   "unknown",
				target: domain.Target{Target: "foo"},
			},
			want: "",
		},
		{
			name: "Remote target with different repo",
			args: args{
				repo:   "github.com/foo/bar",
				target: domain.Target{Target: "foo", GitURL: "github.com/foo2/bar2"},
			},
			want: "github.com/foo2/bar2+foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTarget(tt.args.repo, tt.args.target); got != tt.want {
				t.Errorf("getTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRepo(t *testing.T) {
	type args struct {
		repo   string
		target domain.Target
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Local target with separate repo",
			args: args{
				repo:   "github.com/foo/bar",
				target: domain.Target{Target: "foo"},
			},
			want: "github.com/foo/bar",
		},
		{
			name: "Local target with separate repo and local path",
			args: args{
				repo:   "github.com/foo/bar",
				target: domain.Target{Target: "foo", LocalPath: "baz"},
			},
			want: "github.com/foo/bar",
		},
		{
			name: "Remote target",
			args: args{
				target: domain.Target{Target: "foo", GitURL: "github.com/foo/bar"},
			},
			want: "github.com/foo/bar",
		},
		{
			name: "Remote target with path",
			args: args{
				target: domain.Target{Target: "foo", GitURL: "github.com/foo/bar/baz"},
			},
			want: "github.com/foo/bar",
		},
		{
			name: "Empty target, no repo",
			want: "",
		},
		{
			name: "Unknown repo",
			args: args{
				repo:   "unknown",
				target: domain.Target{Target: "foo"},
			},
			want: "unknown",
		},
		{
			name: "Remote target with different repo",
			args: args{
				repo:   "github.com/foo/bar",
				target: domain.Target{Target: "foo", GitURL: "github.com/foo2/bar2"},
			},
			want: "github.com/foo2/bar2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRepo(tt.args.repo, tt.args.target); got != tt.want {
				t.Errorf("getRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}
