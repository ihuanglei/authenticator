NOW = $(shell date -u '+%Y%m%d%I%M%S')
LDFLAGS += -X "github.com/ihuanglei/authenticator/pkg/build.BuildTime=$(shell date -u '+%Y-%m-%d %H:%M:%S' -d '+8 hour')"
BUILD_FLAGS = -v
OS := $(shell uname -s | tr 'A-Z' 'a-z')

export GO111MODULE=on

all: build-linux build-windows build-darwin

linux: build-linux

windows: build-windows

mac: build-darwin

install: install

doc: build-doc

build-linux: 
	mkdir -p dist/linux
	GOOS=linux go build $(BUILD_FLAGS) -a -ldflags '$(LDFLAGS)' -o dist/linux/authenticator
	cp -r auth.simple.yml README.md docker-compose.yml lib dist/linux
	
build-windows:
	mkdir -p dist/windows
	GOOS=windows go build $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -o dist/windows/authenticator.exe
	cp -r auth.simple.yml README.md docker-compose.yml lib dist/windows

build-darwin:
	mkdir -p dist/darwin
	GOOS=darwin go build $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -o dist/darwin/authenticator
	cp -r auth.simple.yml README.md docker-compose.yml lib dist/darwin

build-doc:
	swag init -g controller/api_doc.go
	rm -f ./docs/docs.go

install:
	mkdir -p /etc/authenticator
	cd dist/linux && cp auth.simple.yml /etc/authenticator/auth.conf.yml && cp -r lib /etc/authenticator &&	cp authenticator /usr/local/bin && mv authenticator.service /lib/systemd/system

clean:
	rm -rf dist
	rm -rf docs


