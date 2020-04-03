ifeq ($(OS),Windows_NT)
  EXECUTABLE_EXTENSION := .exe
else
  EXECUTABLE_EXTENSION :=
endif

GO_FILES = $(shell find . -type f -name '*.go')

all: lzr

lzr: $(GO_FILES)
	cd cmd/lzr && go build && cd ../..
	rm -f lzr
	ln -s cmd/lzr/lzr$(EXECUTABLE_EXTENSION) lzr

clean:
	cd cmd/lzr && go clean
	rm -f lzr
