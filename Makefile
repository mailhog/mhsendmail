VERSION=0.1.9

all: deps release

deps:
	go get github.com/spf13/pflag

release: release-deps
	gox -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" .

release-deps:
	go get github.com/mitchellh/gox

.PNONY: all release release-deps
