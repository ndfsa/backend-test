GOCMD=go
GOBUILD=$(GOCMD) build
GOFLAGS=-gcflags=all="-N -l"

LIBDIR := ./internal

APIDIR := ./cmd/api/
APIGOFILES := $(shell find $(APIDIR) $(LIBDIR) -name "*.go")
APIBIN := .build/apiSrv

AUTHDIR := ./cmd/auth/
AUTHGOFILES := $(shell find $(AUTHDIR) $(LIBDIR) -name "*.go")
AUTHBIN := .build/authSrv

all: $(AUTHBIN) $(APIBIN)

$(AUTHBIN): $(AUTHGOFILES)
	$(GOBUILD) $(GOFLAGS) -o $@ $(AUTHDIR)

$(APIBIN): $(APIGOFILES)
	$(GOBUILD) $(GOFLAGS) -o $@ $(APIDIR)

run: all
	mprocs ./$(AUTHBIN) ./$(APIBIN)

clean:
	rm ./.build/*
	rmdir ./.build

.PHONY: all clean run
