.PHONY:

ifeq ($(OS),Windows_NT)
        EXEEXT := .exe
		BUILD_FLAGS := -ldflags="-H windowsgui"
endif


build:
ifeq ($(OS),Windows_NT)
	rsrc -manifest ./cbr.exe.manifest -ico ./Icon.ico -o ./cmd/rsrc.syso
endif
	mkdir -p build
	cd build && \
		go build ${BUILD_FLAGS} -o cbr${EXEEXT} ../cmd/

clean:
	rm -rf build
	rm -f cmd/rsrc.syso

test:
	go test -v -covermode=count './...'

depends:
	go get bitbucket.org/rj/goey@latest
	go get golang.org/x/net@latest
ifeq ($(OS),Windows_NT)
	go install github.com/akavel/rsrc@latest
endif
