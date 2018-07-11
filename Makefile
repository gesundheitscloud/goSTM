release_platforms := linux darwin
branch_version := $(shell echo `git rev-parse --abbrev-ref HEAD`-SNAPSHOT)
tag_version := $(shell git tag -l --points-at HEAD)
bin_name := goSTM

build-env:
	CGO_ENABLED=1 
	CC=gcc 
	CXX=g++

build-snapshot: build-env
	go build -ldflags "-X main.toolVersion=$(branch_version)" -o bin/$(bin_name) .
	chmod +x bin/$(bin_name)

build-release: build-env
	go build -ldflags "-X main.toolVersion=$(tag_version)" -o bin/$(bin_name) .
	chmod +x bin/$(bin_name)

deps: build-env
	go get github.com/mitchellh/gox
	go get github.com/mitchellh/go-homedir
	# ssh config parser
	go get github.com/kevinburke/ssh_config

travis-test: build-snapshot
	sudo apt-get install openssh-server
	echo "TODO: try to connect to the ssh server with goSTM"

release: deps build-release verify
	echo -n "" > SHA256SUMS
	${GOPATH}/bin/gox -ldflags="-X main.toolVersion=$(version)" -osarch="linux/amd64" -osarch="darwin/amd64" -output="$(bin_name)_$(version)_{{.OS}}_{{.Arch}}" .
	for platform in $(release_platforms); do\
		mv $(bin_name)_$(version)_$${platform}_amd64 $(bin_name); \
		zip $(bin_name)_$(version)_$${platform}_amd64.zip $(bin_name); \
		sha256sum $(bin_name)_$(version)_$${platform}_amd64.zip >> SHA256SUMS; \
		rm $(bin_name); \
	done

clean:
	echo TODO
