### Status
[![Master branch build status](https://travis-ci.com/gesundheitscloud/goSTM.svg?branch=master)](https://travis-ci.com/gesundheitscloud/goSTM)
[![Go Report Card](https://goreportcard.com/badge/github.com/gesundheitscloud/goSTM)](https://goreportcard.com/report/github.com/gesundheitscloud/goSTM)

# goSTM

This is an early MVP.

## Dev

### Install andlabs/ui
#### Mac
```
go get github.com/gesundheitscloud/goSTM
``` 
#### Linux
There is currently a bug with the `libui_linux_amd64.a` file, which makes the setup slightly more compicated. Get [ui](https://github.com/andlabs/ui) via `go get` and [libui](https://github.com/andlabs/libui) manually (without `go get`):
```
go get github.com/andlabs/ui ## This will throw an error, but thats ok
mkdir $GOPATH/src/github.com/andlabs
cd $GOPATH/src/github.com/andlabs
git clone https://github.com/andlabs/libui
cd libui
git checkout 3.5alpha ## There is currently a bug on latest master .. See https://github.com/andlabs/ui/issues/230#issuecomment-289231075
make ## Might throw some errors, but thats fine
cp out/libui.a $GOPATH/src/github.com/andlabs/ui/libui_linux_amd64.a
cd ${GOPATH}/src/github.com/andlabs/ui/
go build
```


### SSH config parser

The following golang ssh parser is used. [ssh_config](https://github.com/kevinburke/ssh_config).

Documentation can be found here [godoc](https://godoc.org/github.com/kevinburke/ssh_config).

### SSH Tunnel

The SSH tunnel in the ssh_tunnel.go file is inspired by this
[gist post](https://gist.github.com/svett/5d695dcc4cc6ad5dd275)

The SSH tunnel is implemented like this.

* open local port and wait for connections
* open fork for incoming connections on the local port (ssh config 'LocalForward' bind address and port)
* this fork opens a SSH remote connection to the ssh host (ssh config 'Hostname')
* open a remote connection on that SSH remote connection to the target host (ssh config 'LocalForward' host and hostport)

### TODOs
[Here](TODO.md)
