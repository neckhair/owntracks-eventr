PACKAGE = github.com/neckhair/owntracks-eventr
COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`
LDFLAGS = -ldflags "-X ${PACKAGE}/eventr.CommitHash=${COMMIT_HASH} -X ${PACKAGE}/eventr.BuildDate=${BUILD_DATE}"
NOGI_LDFLAGS = -ldflags "-X ${PACKAGE}/eventr.BuildDate=${BUILD_DATE}"

VERSION=`git describe --tags`

.PHONY: build clean package
.DEFAULT_GOAL := test

test:
	go get -t ./...
	go test ./...

build:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o build/owntracks-eventr-linux-amd64 ${PACKAGE}

package: build
	echo 2.0 > build/debian-binary
	echo "Package: owntracks-eventr" > build/control
	echo "Version:" ${VERSION} >> build/control
	echo "Architecture: amd64" >> build/control
	echo "Section: net" >> build/control
	echo "Maintainer: Philippe Hässig <[phil@neckhair.ch]>" >> build/control
	echo "Priority: optional" >> build/control
	echo "Description: [Track Owntracks events and write them into a file]" >> build/control
	echo " Built" `date`
	rm -rf build/usr
	mkdir -p build/usr/local/bin
	mkdir -p build/etc/init
	cp build/owntracks-eventr-linux-amd64 build/usr/local/bin/owntracks-eventr
	# cp upstart/*.conf build/etc/init
	sudo chown -R root: build/usr
	sudo chown -R root: build/etc
	tar cvzf build/data.tar.gz -C build usr etc
	tar cvzf build/control.tar.gz -C build control
	cd build && ar rc owntracks-eventr.deb debian-binary control.tar.gz data.tar.gz && cd ..

clean:
	rm -rf build
