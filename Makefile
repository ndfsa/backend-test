GOCMD=go
GOBUILD=$(GOCMD) build
GOFLAGS=-gcflags '-N -l'

BACKDIR=./cmd/back/
BACKGOFILES = $(shell find $(BACKDIR) -name "*.go")
BACKBIN = .build/backSrv

AUTHDIR=./cmd/auth/
AUTHGOFILES = $(shell find $(AUTHDIR) -name "*.go")
AUTHBIN = .build/authSrv

all: $(AUTHBIN) $(BACKBIN)

$(AUTHBIN): $(AUTHGOFILES)
	$(GOBUILD) $(GOFLAGS) -o $@ $(AUTHDIR)

$(BACKBIN): $(BACKGOFILES)
	$(GOBUILD) $(GOFLAGS) -o $@ $(BACKDIR)

run: all
	mprocs ./$(AUTHBIN) ./$(BACKBIN) 'pgrep .*Srv --list-name'

clean:
	rm ./.build/*
	rmdir ./.build

.PHONY: all clean run
