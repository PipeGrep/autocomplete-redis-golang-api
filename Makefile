

APP=		autocomplete

SRC=		src/app.go				\
		src/config.go				\
		src/utils.go				\
		src/redis.go				\
		src/router.go				\
		src/middleware.go			\
		src/models.go				\
		src/index.go				\
		src/search.go				\
		src/get_object.go

PKG=		github.com/garyburd/redigo/redis	\
		github.com/ant0ine/go-json-rest/rest	\
		github.com/fatih/structs



GOPATH := ${PWD}/pkg:${GOPATH}
export GOPATH

default:build

build:
	go build -v -o ./bin/$(APP) $(SRC)

fmt:
	go fmt $(SRC)

run:	build
	./bin/$(APP)

vendor_clean:
	rm -dRf ./pkg/*

vendor_get:
	GOPATH=${PWD}/pkg go get -d -u -v $(PKG)

vendor_update: vendor_get
	rm -rf `find ./pkg/src -type d -name .git` \
	    && rm -rf `find ./pkg/src -type d -name .hg` \
	    && rm -rf `find ./pkg/src -type d -name .bzr` \
	    && rm -rf `find ./pkg/src -type d -name .svn`
