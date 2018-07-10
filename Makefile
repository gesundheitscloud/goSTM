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

# TODO: install like ui deps like this: https://github.com/andlabs/ui/issues/230#issuecomment-289231075
deps: build-env
	#go get gopkg.in/alecthomas/kingpin.v2
	#go get github.com/fatih/color
	#go get gopkg.in/yaml.v2
	go get github.com/mitchellh/gox
	go get github.com/mitchellh/go-homedir

verify:
	echo TODO

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
