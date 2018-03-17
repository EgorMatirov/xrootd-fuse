package main

import (
	"context"
	"os"
	"time"

	"bazil.org/fuse"
	"github.com/EgorMatirov/xrootd"
	"github.com/EgorMatirov/xrootd/requests/stat"
	"log"
)

type node struct {
	client *xrootd.Client
	path   string
}

func (n node) Attr(ctx context.Context, a *fuse.Attr) error {
	st, err := n.client.Stat(ctx, n.path)
	if err != nil {
		log.Print(err)
		return err
	}

	a.Inode = uint64(st.ID)
	a.Valid = 0 // Do not cache attributes

	if st.Flags&stat.FlagIsDir != 0 {
		a.Mode |= os.ModeDir
	}
	if st.Flags&stat.FlagReadable != 0 {
		a.Mode |= 0444
	}
	if st.Flags&stat.FlagWritable != 0 {
		a.Mode |= 0222
	}
	if st.Flags&stat.FlagXset != 0 {
		a.Mode |= 0111
	}

	a.Size = uint64(st.Size)
	a.Mtime = time.Unix(st.ModificationTime, 0)

	return nil
}
