Stash Backup
------------

A command-line backup tool in a similar vein to [Arq](https://www.arqbackup.com/),
but without a UI and with full support for linux. Supporting Amazon (S3, glacier)
and Google cloud storage (nearline and regular).

### To Use:

Install Go:

    brew install go

    or

    sudo apt-get install go

Install `stash`:

    go get github.com/cameronsparr/stash

Add a backup destination (not working yet):

    stash destination add

View all backup destinations (not working yet):

    stash destination list

Add folder to an existing backup destination (not working yet):

    stash folder add

### To Contribute:

To start contributing, get the code and copy the pre-commit hook:

    go get github.com/cameronsparr/stash
    cd $GOPATH/src/github.com/cameronsparr/stash
    cp ./pre-commit ./.git/hooks/pre-commit