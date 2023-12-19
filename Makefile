.PHONY:

ifeq ($(OS),Windows_NT)
        EXEEXT := .exe
		BUILD_FLAGS := -ldflags="-H windowsgui"
endif


build:
	mkdir -p build
	cd build && \
		rsrc -ico ../Icon.ico && \
		go build ${BUILD_FLAGS} -o cbr${EXEEXT} ../cmd/*.go

clean:
	rm -rf build

test:
	go test -v -covermode=count './...'

depends:
	go get bitbucket.org/rj/goey@latest
	go get golang.org/x/net@latest
ifeq ($(OS),Windows_NT)
	go install github.com/akavel/rsrc@latest
endif
