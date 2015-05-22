VERSION=0.1.8

all: release

release: release-deps
	gox -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" .

release-deps:
	go get github.com/mitchellh/gox

.PNONY: all release release-deps
