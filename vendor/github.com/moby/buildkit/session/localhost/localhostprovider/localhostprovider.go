package localhostprovider

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/localhost"
	"github.com/pkg/errors"
	"github.com/tonistiigi/fsutil"
	"github.com/tonistiigi/fsutil/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewLocalhostProvider() (session.Attachable, error) {
	return &localhostProvider{
		m: &sync.Mutex{},
	}, nil
}

type localhostProvider struct {
	m *sync.Mutex
}

func (lp *localhostProvider) Register(server *grpc.Server) {
	localhost.RegisterLocalhostServer(server, lp)
}

func (lp *localhostProvider) Exec(stream localhost.Localhost_ExecServer) error {
	ctx := stream.Context()
	opts, _ := metadata.FromIncomingContext(ctx)
	_ = opts // opts aren't used for anything at the moment

	// first message must contain the command (and no stdin)
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	if len(msg.Command) == 0 {
		return fmt.Errorf("command is empty")
	}
	cmdStr := msg.Command[0]
	args := msg.Command[1:]
	workingDir := msg.Dir

	// it might be possible to run in parallel; but it hasn't been tested.
	lp.m.Lock()
	defer lp.m.Unlock()

	if workingDir != "" {
		err = os.MkdirAll(workingDir, 0755)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	cmd := exec.CommandContext(ctx, cmdStr, args...)
	cmd.Dir = workingDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	m := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(2)
	const readSize = 8196
	var readErr []error
	go func() {
		defer wg.Done()
		for buf := make([]byte, readSize); ; {
			n, err := stdout.Read(buf)
			if n > 0 {
				m.Lock()
				resp := localhost.OutputMessage{
					Stdout: buf[:n],
				}
				err := stream.Send(&resp)
				if err != nil {
					readErr = append(readErr, err)
					m.Unlock()
					return
				}
				m.Unlock()
			}
			if err != nil {
				m.Lock()
				if err != io.EOF {
					readErr = append(readErr, err)
				}
				m.Unlock()
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for buf := make([]byte, readSize); ; {
			n, err := stderr.Read(buf)
			if n > 0 {
				m.Lock()
				resp := localhost.OutputMessage{
					Stderr: buf[:n],
				}
				err := stream.Send(&resp)
				if err != nil {
					readErr = append(readErr, err)
					m.Unlock()
					return
				}
				m.Unlock()
			}
			if err != nil {
				m.Lock()
				if err != io.EOF {
					readErr = append(readErr, err)
				}
				m.Unlock()
				return
			}
		}
	}()

	wg.Wait()
	if len(readErr) != 0 {
		for _, err := range readErr {
			fmt.Fprintf(os.Stderr, "got error while reading from locally-run process: %v\n", err)
		}
		return readErr[0]
	}

	var exitCode int
	status := localhost.DONE
	err = cmd.Wait()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if waitStatus, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = waitStatus.ExitStatus()
			} else {
				status = localhost.KILLED
			}
		} else {
			status = localhost.KILLED
		}
	}

	resp := localhost.OutputMessage{
		ExitCode: int32(exitCode),
		Status:   status,
	}
	if err := stream.Send(&resp); err != nil {
		return err
	}

	return nil
}

func (lp *localhostProvider) Get(stream localhost.Localhost_GetServer) error {
	// first message must contain the path to fetch
	msg, err := stream.Recv()
	if err != nil {
		return err
	}
	path := string(msg.Data)

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return sendDir(stream, path)
	case mode.IsRegular():
		return sendFile(stream, path)
	default:
		panic(fmt.Sprintf("unhandled mode file %v in localhostProvider.Get", mode))
	}
}

func sendFile(stream localhost.Localhost_GetServer, path string) error {
	err := stream.Send(&localhost.BytesMessage{
		Data: []byte{'f', 0x00}, // 'f' denotes a file; 0 denotes version 0 of the send file protocol
	})
	if err != nil {
		return err
	}

	stat, err := fsutil.Stat(path)
	if err != nil {
		return errors.Wrapf(err, "failed to stat %s", path)
	}
	payload, err := stat.Marshal()
	if err != nil {
		return err
	}
	err = stream.Send(&localhost.BytesMessage{
		Data: payload,
	})
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", path)
	}
	defer f.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := f.Read(buf)
		if n > 0 {
			err = stream.Send(&localhost.BytesMessage{
				Data: buf[:n],
			})
		}
		switch err {
		case nil:
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}

func sendDir(stream localhost.Localhost_GetServer, path string) error {
	err := stream.Send(&localhost.BytesMessage{
		Data: []byte{'d', 0x00}, // 'd' denotes a file; 0 denotes version 0 of the send file protocol
	})
	if err != nil {
		return err
	}

	fs, err := fsutil.NewFS(path)
	if err != nil {
		return err
	}
	fs, err = fsutil.NewFilterFS(fs, &fsutil.FilterOpt{
		IncludePatterns: []string{"*"},
	})
	if err != nil {
		return err
	}
	err = fsutil.Send(stream.Context(), stream, fs, nil, nil)
	if err != nil {
		return errors.Wrap(err, "fsutil.Send failed")
	}
	return nil
}

func (lp *localhostProvider) Put(stream localhost.Localhost_PutServer) error {
	// first message must contain the path to fetch
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	mode := msg.Data[0]
	switch mode {
	case 'f':
		version := msg.Data[1]
		if version != 0 {
			panic(fmt.Sprintf("unhandled file version %v", version))
		}
		return receiveFile(stream)
	case 'd':
		version := msg.Data[1]
		if version != 0 {
			panic(fmt.Sprintf("unhandled dir version %v", version))
		}
		return receiveDir(stream)
	default:
		panic(fmt.Sprintf("unhandled mode %v", mode))
	}
}

func receiveFile(stream localhost.Localhost_PutServer) (err error) {
	// first message contains the path
	msg, err := stream.Recv()
	if err != nil {
		return errors.WithStack(err)
	}
	dstPath := string(msg.Data)

	dirPath := path.Dir(dstPath)
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}

	// second message contains the file stat
	msg, err = stream.Recv()
	if err != nil {
		return err
	}
	stat := types.Stat{}
	err = stat.Unmarshal(msg.Data)
	if err != nil {
		return errors.WithStack(err)
	}

	f, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(stat.Mode))
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

outer:
	for {
		msg, err := stream.Recv()
		switch err {
		case nil:
			// ignore
		case io.EOF:
			break outer
		default:
			return errors.WithStack(err)
		}
		_, err = f.Write(msg.Data)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	mtime := time.Unix(0, stat.ModTime)
	err = os.Chtimes(dstPath, mtime, mtime)
	if err != nil {
		return errors.Wrap(err, "failed to change file time")
	}

	// send an empty message to confirm we're done processing
	return stream.Send(&localhost.BytesMessage{})
}

func receiveDir(stream localhost.Localhost_PutServer) error {
	msg, err := stream.Recv()
	if err != nil {
		return errors.WithStack(err)
	}
	dest := string(msg.Data)
	err = os.MkdirAll(dest, 0775)
	if err != nil {
		return errors.WithStack(err)
	}
	uid := syscall.Getuid()
	gid := syscall.Getgid()
	ctx := stream.Context()
	return errors.WithStack(fsutil.Receive(ctx, stream, dest, fsutil.ReceiveOpt{
		Filter: func(p string, stat *types.Stat) bool {
			stat.Uid = uint32(uid)
			stat.Gid = uint32(gid)
			return true
		},
	}))
}
