VERSION=0.1.0
PATH_BUILD=build/
FILE_ARCH=linux_amd64
FILE_COMMAND=pretty-safe-backup
LN_FILE_COMMAND=psb

clean:
	@rm -rf ./build

build: clean
	@$(GOPATH)/bin/goxc \
	  -bc="linux,amd64 linux,arm64" \
	  -pv=$(VERSION) \
	  -d=$(PATH_BUILD) \
	  -build-ldflags "-X main.VERSION=$(VERSION)"

version:
	@echo $(VERSION)

install:
	install -d -m 0755 '/etc/xdg/psb/'
	install -d -m 0755 '/var/log/psb/'
	install $(PATH_BUILD)$(VERSION)/$(FILE_ARCH)/$(FILE_COMMAND) '/usr/bin/$(FILE_COMMAND)'
	ln -s '/usr/bin/$(FILE_COMMAND)' '/usr/bin/$(LN_FILE_COMMAND)'