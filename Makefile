GOCMD=go
GOBUILD=$(GOCMD) build
GOFLAGS=-gcflags=all="-N -l"

LIBDIR=./internal

BACKDIR=./cmd/back/
BACKGOFILES = $(shell find $(BACKDIR) $(LIBDIR) -name "*.go")
BACKBIN = .build/backSrv

AUTHDIR=./cmd/auth/
AUTHGOFILES = $(shell find $(AUTHDIR) $(LIBDIR) -name "*.go")
AUTHBIN = .build/authSrv

all: $(AUTHBIN) $(BACKBIN)

$(AUTHBIN): $(AUTHGOFILES)
	$(GOBUILD) $(GOFLAGS) -o $@ $(AUTHDIR)

$(BACKBIN): $(BACKGOFILES)
	$(GOBUILD) $(GOFLAGS) -o $@ $(BACKDIR)

run: all
	mprocs ./$(AUTHBIN) ./$(BACKBIN)

clean:
	rm ./.build/*
	rmdir ./.build

.PHONY: all clean run
