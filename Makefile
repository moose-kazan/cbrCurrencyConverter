.PHONY:

build:
	mkdir -p build
	fyne build --src ./cmd -o ../build/cbr

clean:
	rm -rf build

test:
	go test -v -covermode=count './...'
