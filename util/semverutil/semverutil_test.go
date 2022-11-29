package semverutil

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Version
		wantErr bool
	}{
		{"v1.2.3", args{"v1.2.3"}, Version{1, 2, 3, ""}, false},
		{"1.2.3", args{"1.2.3"}, Version{1, 2, 3, ""}, false},
		{"123.234.345", args{"123.234.345"}, Version{123, 234, 345, ""}, false},
		{"v1.2.3-alpha", args{"v1.2.3-alpha"}, Version{1, 2, 3, "-alpha"}, false},
		{"1.2.3-alpha", args{"1.2.3-alpha"}, Version{1, 2, 3, "-alpha"}, false},
		{"v1.2.3-alpha.1", args{"v1.2.3-alpha.1"}, Version{1, 2, 3, "-alpha.1"}, false},
		{"1.2.3-alpha.1", args{"1.2.3-alpha.1"}, Version{1, 2, 3, "-alpha.1"}, false},
		{"v1.2.3-alpha.1+001", args{"v1.2.3-alpha.1+001"}, Version{1, 2, 3, "-alpha.1+001"}, false},
		{"1.2.3-alpha.1+001", args{"1.2.3-alpha.1+001"}, Version{1, 2, 3, "-alpha.1+001"}, false},
		{"v1.2.3+001", args{"v1.2.3+001"}, Version{1, 2, 3, "+001"}, false},
		{"1.2.3+001", args{"1.2.3+001"}, Version{1, 2, 3, "+001"}, false},
		{"not-a-version", args{"not-a-version"}, Version{}, true},
		{"v.1.2", args{"v.1.2"}, Version{}, true},
		{"v1.2", args{"v1.2"}, Version{}, true},
		{"1.2", args{"1.2"}, Version{}, true},
		{"1", args{"1"}, Version{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
