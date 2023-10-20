package inputgraph

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/stretchr/testify/require"
)

func TestHashTargetWithDocker(t *testing.T) {
	r := require.New(t)
	target := domain.Target{
		LocalPath: "./testdata/with-docker",
		Target:    "with-docker-load",
	}

	ctx := context.Background()
	cons := conslogging.New(os.Stderr, &sync.Mutex{}, conslogging.NoColor, 0, conslogging.Info)

	org, project, hash, err := HashTarget(ctx, target, cons)
	r.NoError(err)
	r.Equal("earthly-technologies", org)
	r.Equal("core", project)

	hex := fmt.Sprintf("%x", hash)
	r.Equal("9d2903bc18c99831f4a299090abaf94d25d89321", hex)

	path := "./testdata/with-docker/Earthfile"

	tmpDir, err := os.MkdirTemp(os.TempDir(), "with-docker")
	r.NoError(err)

	tmpFile := filepath.Join(tmpDir, "Earthfile")
	defer func() {
		err = os.RemoveAll(tmpDir)
		r.NoError(err)
	}()

	err = copyFile(path, tmpFile)
	r.NoError(err)

	err = replaceInFile(tmpFile, "saved:latest", "other:latest")
	r.NoError(err)

	target = domain.Target{
		LocalPath: tmpDir,
		Target:    "with-docker-load",
	}

	_, _, hash, err = HashTarget(ctx, target, cons)
	r.NoError(err)

	hex = fmt.Sprintf("%x", hash)
	r.Equal("84b6f722421695a7ded144c1b72efb3b8f3339c6", hex)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}

func replaceInFile(path, find, replace string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	dataBytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	data := string(dataBytes)
	data = strings.ReplaceAll(data, find, replace)
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func TestHashTargetWithDockerNoAlias(t *testing.T) {
	r := require.New(t)
	target := domain.Target{
		LocalPath: "./testdata/with-docker",
		Target:    "with-docker-load-no-alias",
	}

	ctx := context.Background()
	cons := conslogging.New(os.Stderr, &sync.Mutex{}, conslogging.NoColor, 0, conslogging.Info)

	org, project, hash, err := HashTarget(ctx, target, cons)
	r.NoError(err)
	r.Equal("earthly-technologies", org)
	r.Equal("core", project)

	hex := fmt.Sprintf("%x", hash)
	r.Equal("d73e37689c7780cbff2cba2de1a23141618b7b14", hex)
}

func TestHashTargetWithDockerArgs(t *testing.T) {
	r := require.New(t)
	target := domain.Target{
		LocalPath: "./testdata/with-docker",
		Target:    "with-docker-load-args",
	}

	ctx := context.Background()
	cons := conslogging.New(os.Stderr, &sync.Mutex{}, conslogging.NoColor, 0, conslogging.Info)

	_, _, _, err := HashTarget(ctx, target, cons)
	r.Error(err)
}
