package localhost

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/pkg/idtools"
	"github.com/moby/buildkit/session"
	"github.com/pkg/errors"
	"github.com/tonistiigi/fsutil"
	"github.com/tonistiigi/fsutil/types"
)

// RunOnLocalHostMagicStr is a magic mount path which is used to signal the RUN command
// should be run locally whenever this mount path exists; the contents of the mounted path is irrelevant.
const RunOnLocalHostMagicStr = "271c67a1-94d9-4241-8bca-cbae334622ae"

// CopyFileMagicStr is a magic command that copies a file from the local system into a snapshot
// it's used as "CopyFileMagicStr src dest"
const CopyFileMagicStr = "39a51ba7-d8c6-43ac-b3aa-f987b2db1ced"

// SendFileMagicStr is a magic command that copies a file from a snapshot to the localhost
// it's used as "SendFileMagicStr src dest"
const SendFileMagicStr = "98325231-d2e6-931c-b12a-84273bca21db"

// Mountable is from buildkit/snapshot; however the snapshot package wont build on darwin
// so we must pull this in here to avoid pulling in linux-specific packages.
type Mountable interface {
	IdentityMapping() *idtools.IdentityMapping
}

// LocalhostExec is called by buildkitd; it connects to the user's client to request the client execute a command localy.
func LocalhostExec(ctx context.Context, c session.Caller, args []string, dir string, stdout, stderr io.Writer) error {
	// stdout and stderr get closed in execOp.Exec()

	client := NewLocalhostClient(c.Conn())
	stream, err := client.Exec(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer stream.CloseSend()

	req := InputMessage{
		Command: args,
		Dir:     dir,
	}
	if err := stream.SendMsg(&req); err != nil {
		return errors.WithStack(err)
	}

	var exitCodeSet bool
	var exitCode int
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.WithStack(err)
		}
		stdout.Write(msg.Stdout)
		stderr.Write(msg.Stderr)
		switch msg.Status {
		case RUNNING:
			//ignore
		case DONE:
			if exitCodeSet {
				panic("received multiple DONE messages (shouldn't happen)")
			}
			exitCode = int(msg.ExitCode)
			exitCodeSet = true
		default:
			return fmt.Errorf("unhandled exit status: %d", msg.Status)
		}
	}

	if exitCode != 0 {
		return fmt.Errorf("exit code: %d", exitCode)
	}

	return nil
}

// LocalhostGet fetches a file or directory located at src on the localhost; and copies it into the mounted snapshot under dest
func LocalhostGet(ctx context.Context, c session.Caller, src, dest string, mount Mountable) error {
	client := NewLocalhostClient(c.Conn())

	stream, err := client.Get(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer stream.CloseSend()

	req := BytesMessage{
		Data: []byte(src),
	}
	if err := stream.SendMsg(&req); err != nil {
		return errors.WithStack(err)
	}

	msg, err := stream.Recv()
	if err != nil {
		return errors.WithStack(err)
	}
	if len(msg.Data) == 0 {
		panic("received GetResponse contains no data; shouldn't happen")
	}
	mode := msg.Data[0]
	switch mode {
	case 'f':
		version := msg.Data[1]
		if version != 0 {
			panic(fmt.Sprintf("unhandled file version %v", version))
		}
		return receiveFile(stream, dest)
	case 'd':
		version := msg.Data[1]
		if version != 0 {
			panic(fmt.Sprintf("unhandled dir version %v", version))
		}
		return receiveDir(stream, dest, mount)
	default:
		panic(fmt.Sprintf("unhandled mode %v", mode))
	}
}

func receiveFile(stream Localhost_GetClient, dest string) (err error) {
	msg, err := stream.Recv()
	if err != nil {
		return errors.WithStack(err)
	}

	stat := types.Stat{}
	err = stat.Unmarshal(msg.Data)
	if err != nil {
		return errors.WithStack(err)
	}

	// change group and system-wide permissions to match user's permissions
	mode := stat.Mode
	umode := (mode & 0700) >> 6
	// ...???XXX?????? -> ...???000000000 -> ...???000000XXX -> ...???000XXXXXX -> ...???XXXXXXXXX
	mode = (mode ^ (mode & 0777)) | umode | (umode << 3) | (umode << 6)

	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, os.FileMode(mode))
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

outer:
	for {
		msg, err := stream.Recv()
		switch err {
		case nil:
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
	err = os.Chtimes(dest, mtime, mtime)
	if err != nil {
		return errors.Wrap(err, "failed to change file time")
	}
	return nil
}

func receiveDir(stream Localhost_GetClient, dest string, mount Mountable) error {
	err := os.MkdirAll(dest, 0775)
	if err != nil {
		return errors.WithStack(err)
	}

	ctx := stream.Context()
	return errors.WithStack(fsutil.Receive(ctx, stream, dest, fsutil.ReceiveOpt{
		Filter: func(p string, stat *types.Stat) bool {
			if idmap := mount.IdentityMapping(); idmap != nil {
				identity, err := idmap.ToHost(idtools.Identity{
					UID: int(stat.Uid),
					GID: int(stat.Gid),
				})
				if err != nil {
					return false
				}
				stat.Uid = uint32(identity.UID)
				stat.Gid = uint32(identity.GID)
			}
			// whatever permissions the user has, give them to group and others as well
			// this matches behavior of gitsource, given that umask is 0
			// ...??XXX?????? -> ...0000000XXX
			umode := (stat.Mode & 0700) >> 6
			// ...???XXX?????? -> ...???000000000 -> ...???000000XXX -> ...???000XXXXXX -> ...???XXXXXXXXX
			stat.Mode = (stat.Mode ^ (stat.Mode & 0777)) | umode | (umode << 3) | (umode << 6)
			return true
		},
	}))
}

// LocalhostPut is for uploading data from buildkit to the localhost
// src is a path on the buildkit host; dest is a path on the localhost
func LocalhostPut(ctx context.Context, c session.Caller, src, dst string) error {
	client := NewLocalhostClient(c.Conn())

	stream, err := client.Put(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer stream.CloseSend()

	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return localhostPutSendDir(stream, src, dst)
	case mode.IsRegular():
		return localhostPutSendFile(stream, src, dst)
	default:
		panic(fmt.Sprintf("unhandled mode file %v in localhostProvider.Get", mode))
	}
}

func localhostPutSendFile(stream Localhost_PutClient, src, dst string) error {
	// first tell localhost-provider server that we're sending a file
	err := stream.Send(&BytesMessage{
		Data: []byte{'f', 0x00}, // 'f' denotes a file; 0 denotes version 0 of the send file protocol
	})
	if err != nil {
		return errors.WithStack(err)
	}

	// next send a path of where we are uploading to
	req := BytesMessage{
		Data: []byte(dst),
	}
	if err := stream.SendMsg(&req); err != nil {
		return errors.WithStack(err)
	}

	stat, err := fsutil.Stat(src)
	if err != nil {
		return errors.Wrapf(err, "failed to stat %s", src)
	}
	payload, err := stat.Marshal()
	if err != nil {
		return err
	}
	err = stream.Send(&BytesMessage{
		Data: payload,
	})
	if err != nil {
		return err
	}

	f, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", src)
	}
	defer f.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrapf(err, "failed to read from %s", src)
		}
		if n > 0 {
			err = stream.Send(&BytesMessage{
				Data: buf[:n],
			})
			if err != nil {
				return err
			}
		}
	}

	err = stream.CloseSend()
	if err != nil {
		return errors.Wrap(err, "stream.CloseSend failed")
	}

	// wait for final message to confirm the client has finished receiving and writing the file
	_, err = stream.Recv()
	return err
}

func localhostPutSendDir(stream Localhost_GetClient, src, dst string) error {
	err := stream.Send(&BytesMessage{
		Data: []byte{'d', 0x00}, // 'd' denotes a file; 0 denotes version 0 of the send file protocol
	})
	if err != nil {
		return err
	}

	err = stream.Send(&BytesMessage{
		Data: []byte(dst),
	})
	if err != nil {
		return err
	}

	fs := fsutil.NewFS(src, &fsutil.WalkOpt{
		IncludePatterns: []string{"*"},
	})
	err = fsutil.Send(stream.Context(), stream, fs, nil, nil)
	if err != nil {
		return errors.Wrap(err, "fsutil.Send failed")
	}
	return nil
}
