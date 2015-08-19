VERSION := $(shell sh -c 'git describe --always --tags')

build: deps
	go build -o stash -ldflags "-X main.Version $(VERSION)" cmd/stash/main.go
	go build -o stashd -ldflags "-X main.Version $(VERSION)" cmd/stashd/main.go

deps:
	go get -t -v ./...

update-deps:
	# stash any changes for `go get` to work
	git stash || return 0
	go get -t -v -u -f ./...
	# go get with `-u -f` checks out the master branch
	git checkout -
	# re-apply stashed changes
	git stash apply || return 0

