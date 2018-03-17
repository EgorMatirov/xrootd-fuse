## What is it?
A fuse-based filesystem for [xrootd](http://xrootd.org/).

## Usage
Mounting a remote directory ```/tmp``` from server ```127.0.0.1:1094``` to the local directory ```~/mount``` with username ```gopher``` can be done by:
```console
   foo@bar:~$ xrootd-fuse -address 127.0.0.1:1094 -remoteDir /tmp -target ~/mount -username gopher
```
After that you can see inside ```~/mount``` content of remote ```/tmp``` folder.