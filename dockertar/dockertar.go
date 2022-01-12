package dockertar

import (
	"archive/tar"
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// GetID returns the docker sha256 ID of the image stored within the given .tar file.
func GetID(tarFilePath string) (string, error) {
	tarFile, err := os.Open(tarFilePath)
	if err != nil {
		return "", errors.Wrapf(err, "open file %s for reading", tarFilePath)
	}
	defer tarFile.Close()
	bufTarFile := bufio.NewReader(tarFile)
	tarR := tar.NewReader(bufTarFile)
	for {
		header, err := tarR.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", errors.Wrapf(err, "reading tar %s", tarFilePath)
		}
		if header.Name == "manifest.json" && !header.FileInfo().IsDir() {
			dt, err := io.ReadAll(tarR)
			if err != nil {
				return "", errors.Wrapf(err, "read manifest.json from tar %s", tarFilePath)
			}
			var jsonData []struct {
				Config string `json:"Config"`
			}
			err = json.Unmarshal(dt, &jsonData)
			if err != nil {
				return "", errors.Wrapf(err, "unmarshal json tar manifest for %s", tarFilePath)
			}
			if len(jsonData) != 1 {
				return "", errors.Errorf("unexpected len != 1 docker manifest in %s", tarFilePath)
			}
			return strings.TrimSuffix(jsonData[0].Config, ".json"), nil
		}
	}
	return "", errors.Errorf("docker tar manifest.json not found in tar %s", tarFilePath)
}
