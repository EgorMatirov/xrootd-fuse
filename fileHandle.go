package main

import (
	"context"
	"log"
	"math"

	"bazil.org/fuse"
	"github.com/EgorMatirov/xrootd"
	"github.com/pkg/errors"
)

type fileHandle struct {
	client *xrootd.Client
	handle [4]byte
}

func (f fileHandle) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) (err error) {
	if req.Size > math.MaxInt32 {
		return errors.Errorf("Requested too large size: %d bytes.", req.Size)
	}

	resp.Data, err = f.client.Read(ctx, f.handle, req.Offset, int32(req.Size))
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (f fileHandle) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	err := f.client.Write(ctx, f.handle, req.Offset, 0, req.Data)
	if err != nil {
		log.Fatal(err)
	}
	resp.Size = len(req.Data) // TODO: Use value from server if it exists
	return err
}

func (f fileHandle) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	err := f.client.Close(ctx, f.handle, 0)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (f fileHandle) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	err := f.client.Sync(ctx, f.handle)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
