# Name
BINARY=wapty

# Variables
VERSION=0.2.0
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

.DEFAULT_GOAL: ${BINARY}

# Just build the wapty
# TODO call gopherjs
${BINARY}: rebind
	# Building the executable.
	go build ${LDFLAGS} -o ${BINARY}

fast:
	-rm ui/rice-box.go >& /dev/null #This will make rice use data that is on disk, creates a lighter executable
	go run ${LDFLAGS} *.go

test:
	go test -v ./...

rebind:
	# Cleaning and re-embedding assets
	cd ui && rm rice-box.go >& /dev/null; rice embed-go

install: installdeps rebind
	# Installing the executable
	go install ${LDFLAGS}

installdeps:
	# Installing dependencies to embed assets
	go get -u github.com/GeertJohan/go.rice/...
	# Installing dependencies to build JS
	go get -u github.com/gopherjs/jquery
	go get -u github.com/gopherjs/jsbuiltin
	go get -u github.com/gopherjs/websocket/...
	go get -u github.com/gopherjs/gopherjs

clean:
	-rm ui/rice-box.go
	go clean
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
