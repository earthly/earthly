package conslogging

import (
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/hashicorp/go-multierror"
)

type BundleBuilder struct {
	RootPath string

	logmap map[string]*strings.Builder
}

func NewBundleBuilder(rootPath string) *BundleBuilder {
	return &BundleBuilder{
		RootPath: rootPath,
		logmap:   map[string]*strings.Builder{},
	}
}

func (bb *BundleBuilder) PrefixWriter(prefix string) io.Writer {
	if builder, ok := bb.logmap[prefix]; ok {
		return builder
	}

	writer := &strings.Builder{}
	bb.logmap[prefix] = writer
	return writer
}

func (bb *BundleBuilder) WriteToDisk() error {
	var err error
	for prefix, lines := range bb.logmap {
		tgtErr := ioutil.WriteFile(path.Join(bb.RootPath, strings.TrimSpace(prefix)), []byte(lines.String()), 0666)
		if err != nil {
			err = multierror.Append(err, tgtErr)
		}
	}

	return err
}
