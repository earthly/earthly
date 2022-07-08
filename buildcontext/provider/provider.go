// Package provider is heavily based on fsSyncProvider in github.com/moby/buildkit/session/filesync.
// The key difference between BuildContextProvider and fsSyncProvider is that in
// BuildContextProvider, the dirs can be added incrementally after the construction.
package provider

import (
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/earthly/earthly/conslogging"

	"github.com/dustin/go-humanize"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/filesync"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tonistiigi/fsutil"
	fstypes "github.com/tonistiigi/fsutil/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	keyOverrideExcludes   = "override-excludes"
	keyIncludePatterns    = "include-patterns"
	keyExcludePatterns    = "exclude-patterns"
	keyFollowPaths        = "followpaths"
	keyDirName            = "dir-name"
	keyExporterMetaPrefix = "exporter-md-"
)

var _ session.Attachable = (*BuildContextProvider)(nil)
var _ filesync.FileSyncServer = (*BuildContextProvider)(nil)

// BuildContextProvider is a BuildKit attachable which provides local files as part
// of the build context.
type BuildContextProvider struct {
	p      progressCb
	doneCh chan error

	mu   sync.Mutex
	dirs map[string]SyncedDir

	console conslogging.ConsoleLogger
}

// SyncedDir is a directory to be synced across.
type SyncedDir struct {
	Name     string
	Dir      string
	Excludes []string
	Map      func(string, *fstypes.Stat) fsutil.MapResult
}

// NewBuildContextProvider creates a new provider for sending build context files from client.
func NewBuildContextProvider(console conslogging.ConsoleLogger) *BuildContextProvider {
	return &BuildContextProvider{
		dirs:    map[string]SyncedDir{},
		console: console,
	}
}

// AddDirs adds local directories to the context.
func (bcp *BuildContextProvider) AddDirs(dirs map[string]string) {
	bcp.mu.Lock()
	defer bcp.mu.Unlock()
	for dirName, dir := range dirs {
		bcp.addDir(dirName, dir)
	}
}

// AddDir adds a single local directory to the context.
func (bcp *BuildContextProvider) AddDir(dirName, dir string) {
	bcp.mu.Lock()
	defer bcp.mu.Unlock()

	bcp.addDir(dirName, dir)
}

func (bcp *BuildContextProvider) addDir(dirName, dir string) {
	resetUIDAndGID := func(p string, st *fstypes.Stat) fsutil.MapResult {
		st.Uid = 0
		st.Gid = 0
		return fsutil.MapResultKeep
	}
	sd := SyncedDir{
		Name: dirName,
		Dir:  dir,
		Map:  resetUIDAndGID,
	}

	bcp.dirs[sd.Name] = sd
}

// Register registers the attachable.
func (bcp *BuildContextProvider) Register(server *grpc.Server) {
	filesync.RegisterFileSyncServer(server, bcp)
}

// DiffCopy implements the DiffCopy attachable.
func (bcp *BuildContextProvider) DiffCopy(stream filesync.FileSync_DiffCopyServer) error {
	return bcp.handle("diffcopy", stream)
}

// TarStream implements the DiffCopy attachable.
func (bcp *BuildContextProvider) TarStream(stream filesync.FileSync_TarStreamServer) error {
	return bcp.handle("tarstream", stream)
}

func (bcp *BuildContextProvider) handle(method string, stream grpc.ServerStream) (retErr error) {
	var pr *protocol
	for _, p := range supportedProtocols {
		if method == p.name && isProtoSupported(p.name) {
			pr = &p
			break
		}
	}
	if pr == nil {
		return errors.New("failed to negotiate protocol")
	}

	opts, _ := metadata.FromIncomingContext(stream.Context()) // if no metadata continue with empty object

	dirName := ""
	name, ok := opts[keyDirName]
	if ok && len(name) > 0 {
		dirName = name[0]
	}

	dir, err := bcp.getDir(dirName)
	if err != nil {
		return err
	}

	excludes := opts[keyExcludePatterns]
	if len(dir.Excludes) != 0 && (len(opts[keyOverrideExcludes]) == 0 || opts[keyOverrideExcludes][0] != "true") {
		excludes = dir.Excludes
	}
	includes := opts[keyIncludePatterns]

	followPaths := opts[keyFollowPaths]

	var mutex sync.Mutex
	console := bcp.console.WithPrefixAndSalt("context", dir.Dir)
	numStats := 0
	numSends := 0
	verboseProgressCB := func(relPath string, status fsutil.VerboseProgressStatus, numBytes int) {
		mutex.Lock()
		defer mutex.Unlock()
		fullPath := path.Join(dir.Dir, relPath)
		switch status {
		case fsutil.StatusStat:
			numStats++
			//console.VerbosePrintf("sent file stat for %s\n", fullPath) ignored as it is too verbose. TODO add different verbose levels to support ExtraVerbosePrintf
		case fsutil.StatusSent:
			console.VerbosePrintf("sent data for %s (%s)\n", fullPath, humanize.Bytes(uint64(numBytes)))
			numSends++
		case fsutil.StatusFailed:
			console.VerbosePrintf("sent data for %s failed\n", fullPath)
		case fsutil.StatusSkipped:
			console.VerbosePrintf("ignoring %s\n", fullPath)
		default:
			console.Warnf("unhandled progress status %v (path=%s, numBytes=%d)\n", status, fullPath, numBytes)
		}
	}

	progress := func(numBytes int, last bool) {
		mutex.Lock()
		defer mutex.Unlock()
		if last {
			console.Printf("transferred %d file(s) for context %s (%s, %d file/dir stats)", numSends, dir.Dir, humanize.Bytes(uint64(numBytes)), numStats)
		}
	}

	var doneCh chan error
	if bcp.doneCh != nil {
		doneCh = bcp.doneCh
		bcp.doneCh = nil
	}
	err = pr.sendFn(stream, fsutil.NewFS(dir.Dir, &fsutil.WalkOpt{
		ExcludePatterns:   excludes,
		IncludePatterns:   includes,
		FollowPaths:       followPaths,
		Map:               dir.Map,
		VerboseProgressCB: verboseProgressCB,
	}), progress, verboseProgressCB)
	if doneCh != nil {
		if err != nil {
			doneCh <- err
		}
		close(doneCh)
	}
	return err
}

func (bcp *BuildContextProvider) getDir(dirName string) (SyncedDir, error) {
	bcp.mu.Lock()
	defer bcp.mu.Unlock()
	dir, ok := bcp.dirs[dirName]
	if !ok {
		return SyncedDir{}, status.Errorf(codes.NotFound, "no access allowed to dir %q", dirName)
	}
	return dir, nil
}

// SetNextProgressCallback sets the progress callback function.
func (bcp *BuildContextProvider) SetNextProgressCallback(f func(int, bool), doneCh chan error) {
	bcp.p = f
	bcp.doneCh = doneCh
}

type progressCb func(int, bool)

type protocol struct {
	name   string
	sendFn func(stream filesync.Stream, fs fsutil.FS, progress progressCb, verboseProgress fsutil.VerboseProgressCB) error
	recvFn func(stream grpc.ClientStream, destDir string, cu filesync.CacheUpdater, progress progressCb, mapFunc func(string, *fstypes.Stat) bool) error
}

func isProtoSupported(p string) bool {
	// TODO: this should be removed after testing if stability is confirmed
	if override := os.Getenv("BUILD_STREAM_PROTOCOL"); override != "" {
		return strings.EqualFold(p, override)
	}
	return true
}

var supportedProtocols = []protocol{
	{
		name:   "diffcopy",
		sendFn: sendDiffCopy,
		recvFn: recvDiffCopy,
	},
}

func sendDiffCopy(stream filesync.Stream, fs fsutil.FS, progress progressCb, verboseProgress fsutil.VerboseProgressCB) error {
	return errors.WithStack(fsutil.Send(stream.Context(), stream, fs, progress, verboseProgress))
}

func recvDiffCopy(ds grpc.ClientStream, dest string, cu filesync.CacheUpdater, progress progressCb, filter func(string, *fstypes.Stat) bool) error {
	st := time.Now()
	defer func() {
		logrus.Debugf("diffcopy took: %v", time.Since(st))
	}()
	var cf fsutil.ChangeFunc
	var ch fsutil.ContentHasher
	if cu != nil {
		cu.MarkSupported(true)
		cf = cu.HandleChange
		ch = cu.ContentHasher()
	}
	return errors.WithStack(fsutil.Receive(ds.Context(), ds, dest, fsutil.ReceiveOpt{
		NotifyHashed:  cf,
		ContentHasher: ch,
		ProgressCb:    progress,
		Filter:        fsutil.FilterFunc(filter),
	}))
}
