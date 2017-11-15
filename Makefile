VERSION=0.0.1
PATH_BUILD=build/
FILE_ARCH=linux_amd64
FILE_COMMAND=psb

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
	install -d -m 0755 '/usr/share/icons/hicolor/48x48/apps/'
	install -m 0644 icon/48x48/psb.png /usr/share/icons/hicolor/48x48/apps/psb.png
	install $(PATH_BUILD)$(VERSION)/$(FILE_ARCH)/$(FILE_COMMAND) '/usr/bin/$(FILE_COMMAND)'