.PHONY:

ifeq ($(OS),Windows_NT)
        EXEEXT := .exe
		BUILD_FLAGS := -ldflags="-H windowsgui"
endif


build:
	mkdir -p build
	go build ${BUILD_FLAGS} -o build/cbr${EXEEXT} ./cmd/*.go

clean:
	rm -rf build

test:
	go test -v -covermode=count './...'
