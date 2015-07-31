[![Circle CI](https://circleci.com/gh/sparrc/stash.png?circle-token=:circle-token)](https://circleci.com/gh/sparrc/stash)

Stash Backup
------------

A command-line backup tool in a similar vein to
[Arq](https://www.arqbackup.com/), but aiming to be controlled completely
from a CLI rather than a UI, and with full support for linux & mac.

Currently planned to support Amazon (S3, glacier) and Google cloud storage.

### To Use:

Install Go:

    brew install go

    or

    sudo apt-get install golang

Install `stash`:

    go get github.com/sparrc/stash/...

Add a backup destination:

    stash destination add

View all backup destinations:

    stash destination list

Add folder to an existing backup destination (not working yet):

    stash folder add

### To Contribute:

First [Setup your $GOPATH and go environment](https://golang.org/doc/code.html)

Get `stash`:

    go get github.com/sparrc/stash/...

Go to source directory:

    cd $GOPATH/src/github.com/sparrc/stash

Run tests (from root dir)

    go test ./...

Re-compile and install binaries

    go install ./...

OR compile and run directly

    go run cmd/stashd/main.go
