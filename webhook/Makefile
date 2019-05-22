TARGET      =   k8s_webhook_calm
SRC_DIR     =   $(notdir $(shell pwd))
SRCS        =   $(wildcard *.go)
#ALLSRCS    =   $(wildcard *.go ./forwardsvr/*.go ../TYSCProtocol/*.go)
VERSION     =   0.0.1
BUILDTIME   =   `date +%FT%T%z`
LDFLAGS     =   -ldflags "-X main.version=${VERSION} -X main.buildtime=${BUILDTIME}"
GCFLAGS     =   -gcflags "-N -l -m -m"
#GCFLAGS    =   
BINDIR      =   ./bin
GOTAGS      :=  -tags static -tags deadlock

.PHONY: clean all vet fmt

all : $(TARGET)

$(TARGET) : % : $(BINDIR)/%

$(BINDIR)/% : $(SRCS)
	gofmt -l -w $(SRCS)
	CGO_ENABLED=1 go build -race -mod vendor $(GOTAGS) -v -x $(GCFLAGS) $(LDFLAGS) -o $@ $(SRCS) 

clean :
	go clean -i 
	-rm $(BINDIR)/$(TARGET)

vet :
	@echo "Running vet..."
	@go vet -v $(go list ./...| grep -v /test/)	

lint :
	@echo "Running golint..."
	@golint ./forwardsvr/...

fmt:
	gofmt -l -w $(SRCS)