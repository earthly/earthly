package base

import (
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/util/cliutil"
)

func (cli *CLI) InitFrontend(cliCtx *cli.Context) error {
	// command line option overrides the config which overrides the default value
	if !cliCtx.IsSet("buildkit-image") && cli.Cfg().Global.BuildkitImage != "" {
		cli.Flags().BuildkitdImage = cli.Cfg().Global.BuildkitImage
	}

	if cli.Flags().UseTickTockBuildkitImage {
		if cliCtx.IsSet("buildkit-image") {
			return fmt.Errorf("the --buildkit-image and --ticktock flags are mutually exclusive")
		}
		if cli.Cfg().Global.BuildkitImage != "" {
			return fmt.Errorf("the --ticktock flags can not be used in combination with the buildkit_image config option")
		}
		cli.Flags().BuildkitdImage += "-ticktock"
	}

	bkURL, err := url.Parse(cli.Flags().BuildkitHost) // Not validated because we already did that when we calculated it.
	if err != nil {
		return errors.Wrap(err, "failed to parse generated buildkit URL")
	}

	if bkURL.Scheme == "tcp" && cli.Cfg().Global.TLSEnabled {
		cli.Flags().BuildkitdSettings.ClientTLSCert = cli.Cfg().Global.ClientTLSCert
		cli.Flags().BuildkitdSettings.ClientTLSKey = cli.Cfg().Global.ClientTLSKey
		cli.Flags().BuildkitdSettings.TLSCA = cli.Cfg().Global.TLSCACert
		cli.Flags().BuildkitdSettings.ServerTLSCert = cli.Cfg().Global.ServerTLSCert
		cli.Flags().BuildkitdSettings.ServerTLSKey = cli.Cfg().Global.ServerTLSKey
	}

	cli.Flags().BuildkitdSettings.AdditionalArgs = cli.Cfg().Global.BuildkitAdditionalArgs
	cli.Flags().BuildkitdSettings.AdditionalConfig = cli.Cfg().Global.BuildkitAdditionalConfig
	cli.Flags().BuildkitdSettings.Timeout = time.Duration(cli.Cfg().Global.BuildkitRestartTimeoutS) * time.Second
	cli.Flags().BuildkitdSettings.Debug = cli.Flags().Debug
	cli.Flags().BuildkitdSettings.BuildkitAddress = cli.Flags().BuildkitHost
	cli.Flags().BuildkitdSettings.LocalRegistryAddress = cli.Flags().LocalRegistryHost
	cli.Flags().BuildkitdSettings.UseTCP = bkURL.Scheme == "tcp"
	cli.Flags().BuildkitdSettings.UseTLS = cli.Cfg().Global.TLSEnabled
	cli.Flags().BuildkitdSettings.MaxParallelism = cli.Cfg().Global.BuildkitMaxParallelism
	cli.Flags().BuildkitdSettings.CacheSizeMb = cli.Cfg().Global.BuildkitCacheSizeMb
	cli.Flags().BuildkitdSettings.CacheSizePct = cli.Cfg().Global.BuildkitCacheSizePct
	cli.Flags().BuildkitdSettings.CacheKeepDuration = cli.Cfg().Global.BuildkitCacheKeepDurationS
	cli.Flags().BuildkitdSettings.EnableProfiler = cli.Flags().EnableProfiler
	cli.Flags().BuildkitdSettings.NoUpdate = cli.Flags().NoBuildkitUpdate

	// ensure the MTU is something allowable in IPv4, cap enforced by type. Zero is autodetect.
	if cli.Cfg().Global.CniMtu != 0 && cli.Cfg().Global.CniMtu < 68 {
		return errors.New("invalid overridden MTU size")
	}
	cli.Flags().BuildkitdSettings.CniMtu = cli.Cfg().Global.CniMtu

	if cli.Cfg().Global.IPTables != "" && cli.Cfg().Global.IPTables != "iptables-legacy" && cli.Cfg().Global.IPTables != "iptables-nft" {
		return errors.New(`invalid overridden iptables name. Valid values are "iptables-legacy" or "iptables-nft"`)
	}
	cli.Flags().BuildkitdSettings.IPTables = cli.Cfg().Global.IPTables
	earthlyDir, err := cliutil.GetOrCreateEarthlyDir(cli.Flags().InstallationName)
	if err != nil {
		return errors.Wrap(err, "failed to get earthly dir")
	}
	cli.Flags().BuildkitdSettings.StartUpLockPath = filepath.Join(earthlyDir, "buildkitd-startup.lock")

	return nil
}
