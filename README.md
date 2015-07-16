Stash Backup
------------

A command-line backup tool in a similar vein to [Arq](https://www.arqbackup.com/),
but without a UI and with full support for linux. Supporting Amazon (S3, glacier)
and Google cloud storage (nearline and regular).

### To Use:

Install Go:

    brew install go

    or

    sudo apt-get install golang

Install `stash`:

    go get github.com/sparrc/stash/...

Add a backup destination (not working yet):

    stash destination add

View all backup destinations (not working yet):

    stash destination list

Add folder to an existing backup destination (not working yet):

    stash folder add

### To Contribute:

First [Setup your $GOPATH and go environment](https://golang.org/doc/code.html)

Get `stash`:

    go get github.com/sparrc/stash/...

Change your remote to push to Phabricator:

    cd $GOPATH/src/github.com/sparrc/stash
    git remote remove origin
    git remote add origin ssh://git@phabricator.csparr.net/diffusion/STASH/stash.git

Runn tests (from root dir)

    go test ./...

Re-compile and install binaries

    go install ./...

OR compile and run directly

    go run cmd/stashd/main.go
