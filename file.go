package main

import (
	"context"
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/EgorMatirov/xrootd/requests/open"
)

type File struct {
	*node
}

func flagsToOptions(flags fuse.OpenFlags) (options open.Options) {
	if flags.IsReadOnly() {
		options |= open.OptionsOpenRead
	} else {
		options |= open.OptionsOpenUpdate
	}

	if flags == fuse.OpenCreate {
		options |= open.OptionsNew
	}

	return
}

func flagsToMode(flags fuse.OpenFlags) (mode open.Mode) {
	if flags.IsReadOnly() {
		mode |= open.ModeOwnerRead
	} else {
		mode |= open.ModeOwnerWrite
	}
	return
}

func (f File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	var options = flagsToOptions(req.Flags)
	var mode = flagsToMode(req.Flags)

	handle, err := f.client.Open(ctx, f.path, mode, options)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	resp.Flags |= fuse.OpenDirectIO
	return fileHandle{f.client, handle}, nil
}

func (f File) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	// TODO: Call Truncate, so if we do something like "echo 1 > ./file" it will override whole file
	return nil
}

func (f File) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	// TODO: need to have access to handle. See bug https://github.com/bazil/fuse/issues/92
	return nil
}
