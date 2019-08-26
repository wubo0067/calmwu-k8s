TARGET      =   grpc_test
SRCS        =   $(wildcard *.go)
ALLSRCS     =   $(wildcard *.go ./user/*.go)
VERSION     =   1.0.0
BUILDTIME   =   `date +%FT%T%z`
LDFLAGS     =   -ldflags "-X main.version=${VERSION} -X main.buildtime=${BUILDTIME}"
GCFLAGS     =   -gcflags "-N -l"
#GCFLAGS     =   
GOTAGS      :=  

.PHONY: clean all

all : $(TARGET)

$(TARGET) : $(ALLSRCS)
	make -C user
	go build $(GOTAGS) -mod vendor -v -x $(GCFLAGS) $(LDFLAGS) -o $@ $(SRCS)

clean :
	go clean -i 
	-rm $(BINDIR)/$(TARGET)