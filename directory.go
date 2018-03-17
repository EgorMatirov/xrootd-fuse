package main

import (
	"context"
	"path"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/EgorMatirov/xrootd/requests/open"
	"github.com/EgorMatirov/xrootd/requests/stat"
)

type Dir struct {
	*node
}

func (d Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	st, err := d.client.Stat(ctx, path.Join(d.path, name))
	if err != nil {
		return nil, fuse.ENOENT // assume that for now. TODO: check error code
	}

	n := &node{d.client, path.Join(d.path, name)}

	if st.Flags&stat.FlagIsDir > 0 {
		return Dir{n}, nil
	}

	return File{n}, nil
}

func (d Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	dirNames, err := d.client.Dirlist(ctx, d.path)
	if err != nil {
		return nil, err
	}

	var dirs = make([]fuse.Dirent, len(dirNames))
	for i := 0; i < len(dirNames); i++ {
		dirs[i] = fuse.Dirent{Name: dirNames[i]}
	}

	return dirs, nil
}

func (d Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	var options = open.OptionsNew | flagsToOptions(req.Flags)

	filePath := path.Join(d.path, req.Name)
	handle, err := d.client.Open(ctx, filePath, open.ModeOwnerWrite|open.ModeOwnerRead, options)

	if err != nil {
		return nil, nil, err
	}

	return File{&node{d.client, filePath}}, fileHandle{d.client, handle}, nil
}
