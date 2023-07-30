GOCMD=go
GOBUILD=$(GOCMD) build
GOFLAGS=-gcflags '-N -l'
RM=/bin/rm

BACKGOFILES = $(shell find ./testAPI/ -name "*.go")
BACKBIN = .build/backSrv

AUTHGOFILES = $(shell find ./auth/ -name "*.go")
AUTHBIN = .build/authSrv

all: $(AUTHBIN) $(BACKBIN)

$(AUTHBIN): $(AUTHGOFILES)
	$(GOBUILD) $(GOFLAGS) -o $@ ./auth/srv/

$(BACKBIN): $(BACKGOFILES)
	$(GOBUILD) $(GOFLAGS) -o $@ ./testAPI/srv/

run: all
	mprocs ./$(AUTHBIN) ./$(BACKBIN) 'pgrep .*Srv --list-name'

clean:
	$(RM) -rf ./build

.PHONY: all clean run
