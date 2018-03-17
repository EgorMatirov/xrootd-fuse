package main

import (
	"context"
	"flag"
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/EgorMatirov/xrootd"
)

var target = flag.String("target", "", "Mount point")
var address = flag.String("address", "", "Address of xrootd server")
var username = flag.String("username", "", "Username which should be used in login process")
var remoteDir = flag.String("remoteDir", "", "Remote directory which should be mounted")

func main() {
	flag.Parse()

	const xrootdFuse = "xrootd-fuse"
	c, err := fuse.Mount(
		*target,
		fuse.FSName(xrootdFuse),
		fuse.Subtype(xrootdFuse),
		fuse.VolumeName(xrootdFuse),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	client, err := xrootd.New(context.Background(), *address)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = client.Login(context.Background(), *username); err != nil {
		log.Fatal(err)
	}

	if err = fs.Serve(c, FS{*remoteDir, client}); err != nil {
		log.Fatal(err)
	}

	<-c.Ready

	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}

type FS struct {
	root   string
	client *xrootd.Client
}

func (fs FS) Root() (fs.Node, error) {
	n := &node{fs.client, fs.root}
	return Dir{n}, nil
}
