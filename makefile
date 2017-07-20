# Name
BINARY=wapty

# Variables
VERSION=0.2.0
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Commit=${BUILD}"

.DEFAULT_GOAL: ${BINARY}

# Just build the wapty
# TODO call gopherjs
${BINARY}: buildjs rebind
	# Building the executable.
	go build ${LDFLAGS} -o ${BINARY}

run:
	# This will make rice use data that is on disk, creates a lighter executable
	# and it is faster to build
	-rm ui/rice-box.go >& /dev/null
	# Generating JS
	cd ui/gopherjs/ && gopherjs build -o ../static/gopherjs.js
	# Done generating JS, launching wapty
	go run ${LDFLAGS} wapty.go

fast: run

test: buildjs rebind
	go test ${LDFLAGS} ./...

testv: buildjs rebind
	go test -v -x -race ${LDFLAGS} ./...

buildjs:
	# Regenerating minified js
	cd ui/gopherjs/ && gopherjs build -m -o ../static/gopherjs.js 
	# Remove mappings
	rm ui/static/gopherjs.js.map

rebind:
	# Cleaning and re-embedding assets
	cd ui && rm rice-box.go 1>/dev/null 2>/dev/null; rice embed-go

install: buildjs rebind
	# Installing the executable
	go install ${LDFLAGS}

installdeps:
	# Installing dependencies to embed assets
	go get -u github.com/GeertJohan/go.rice/...
	# Installing dependencies to build JS
	go get -u github.com/gopherjs/gopherjs
	go get -u github.com/gopherjs/websocket/...
	# Installing Decode dependencies
	go get -u github.com/fatih/color
	go get -u github.com/pmezard/go-difflib/difflib

clean:
	# Cleaning all generated files
	-rm ui/rice-box.go
	-rm ui/static/gopherjs.js*
	go clean
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
