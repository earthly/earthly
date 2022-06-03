package satelliteutil

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/earthly/earthly/conslogging"
)

type SatelliteClient struct {
	dataRoot string
}

type SatelliteClientConfig struct {
	ConfigPath string
}

func NewSatelliteClient(ctx context.Context, console conslogging.ConsoleLogger, cfg SatelliteClientConfig) (*SatelliteClient, error) {
	// Configure on-disk storage
	satelliteDataRoot := filepath.Join(cfg.ConfigPath, "satellite")

	console.WithPrefix("satellite").Printf("Ensuring satellite data root exists: %s", satelliteDataRoot)
	err := os.MkdirAll(satelliteDataRoot, 0755)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create satellite config root %s", satelliteDataRoot)
	}

	// Build satellite client
	console.WithPrefix("satellite").Printf("Client configuration loaded!")
	return &SatelliteClient{
		dataRoot: satelliteDataRoot,
	}, nil
}

func (s *SatelliteClient) CreateSatellite(ctx context.Context, console conslogging.ConsoleLogger, name string) (*SatelliteInfo, error) {
	//console.WithPrefix("satellite").Printf("Ensuring data dir: %s", s.dataRoot)
	//err := s.ensureSatelliteDataDir()
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not ensure satellite data directory")
	//}
	//
	//console.WithPrefix("satellite").Printf("Creating satellite %s", name)
	//
	//// cloud.LaunchSatellite
	//
	//console.WithPrefix("satellite").Printf("Satellite creation complete!")
	//console.WithPrefix("satellite").Printf("Public DNS is %s", info.PublicDNS)
	//
	//return info, nil
	return nil, nil
}

func (s *SatelliteClient) ListSatellites(ctx context.Context) ([]string, error) {
	//
	//// Make the listing stable for CLI convenience
	//sort.Strings(satellites)
	//
	//return satellites, nil
	return nil, nil
}

type SatelliteInfo struct {
	PublicIP      string
	PublicDNS     string
	InstanceID    string
	SatelliteName string
}

func (s *SatelliteClient) GetSatellite(ctx context.Context, satelliteName string) (*SatelliteInfo, error) {
	return nil, nil
}

func (s *SatelliteClient) DeleteSatellite(ctx context.Context, console conslogging.ConsoleLogger, satelliteName string) error {
	//console.WithPrefix("satellite").Printf("Terminating instance with id %s", parsedARN.Resource)
	return nil
}

func (s *SatelliteClient) ensureSatelliteDataDir() error {
	err := os.MkdirAll(s.dataRoot, 0755)
	if err != nil {
		return errors.Wrapf(err, "failed to create data dir: %s", s.dataRoot)
	}

	return nil
}
