PACKAGE = github.com/neckhair/owntracks-eventr
COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`
VERSION=`git describe --tags`
LDFLAGS = -ldflags "-X ${PACKAGE}/cmd.CommitHash=${COMMIT_HASH} -X ${PACKAGE}/cmd.BuildDate=${BUILD_DATE} -X ${PACKAGE}/cmd.Version=${VERSION}"
NOGI_LDFLAGS = -ldflags "-X ${PACKAGE}/eventr.BuildDate=${BUILD_DATE}"


.PHONY: build clean package
.DEFAULT_GOAL := test

test:
	go get -t ./...
	go test ./...

build:
	mkdir -p build
	go get -t ./...
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o build/owntracks-eventr-linux-amd64 ${PACKAGE}
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o build/owntracks-eventr-darwin-amd64 ${PACKAGE}

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
	rm -rf build/usr build/etc build/var
	mkdir -p build/usr/local/bin
	mkdir -p build/etc/systemd/system
	mkdir -p build/var/log
	mkdir -p build/var/lib/owntracks-eventr
	cp build/owntracks-eventr-linux-amd64 build/usr/local/bin/owntracks-eventr
	cp debian/owntracks-eventr.service build/etc/systemd/system/
	cp debian/owntracks-eventr.yml build/etc/
	touch build/var/log/owntracks-eventr.log
	touch build/var/lib/owntracks-eventr/events.log
	sudo chown -R root: build/usr
	sudo chown -R root: build/etc
	sudo chown -R root: build/var
	tar cvzf build/data.tar.gz -C build usr etc var
	tar cvzf build/control.tar.gz -C build control
	cd build && ar rc owntracks-eventr.deb debian-binary control.tar.gz data.tar.gz && cd ..

clean:
	rm -rf build
