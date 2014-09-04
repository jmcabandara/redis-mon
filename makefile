# Copyright (C) 2014 Zenoss, Inc
#
# redis-mon is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 2 of the License, or
# (at your option) any later version.
#
# redis-mon is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with Foobar. If not, see <http://www.gnu.org/licenses/>.

## setup all environment stuff
FULL_NAME = redis-mon
VERSION := $(shell cat ./VERSION)
DATE := $(shell date -u)
GIT_COMMIT ?= $(shell ./hack/gitstatus.sh)
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
# jenkins default, jenkins-${JOB_NAME}-${BUILD_NUMBER}
BUILD_TAG ?= 0
LDFLAGS = -ldflags "-X main.Version $(VERSION) -X main.Gitcommit '$(GIT_COMMIT)' -X main.Gitbranch '$(GIT_BRANCH)' -X main.Date '$(DATE)' -X main.Buildtag '$(BUILD_TAG)'"
MAINTAINER = dev@zenoss.com
LICENSE = GPLv2
VENDOR = Zenoss
URL = https://github.com/zenoss/redis-mon
PKGROOT=/tmp/pkgroot
DUID ?= $(shell id -u)
DGID ?= $(shell id -g)
DESCRIPTION := A redis utility that reposts statistics to zenoss metric consumer
GODEPS_FILES := $(shell find Godeps/)

redis-mon: VERSION *.go hack/* makefile $(GODEPS_FILES)
	godep go build ${LDFLAGS}
	chown $(DUID):$(DGID) $(FULL_NAME)

docker-deb: redis-mon-build
	docker run -rm -v `pwd`:/go/src/github.com/zenoss/redis-mon -e DUID=$(DUID) -e DGID=$(DGID) zenoss/redis-mon-build:$(VERSION) make deb

docker-rpm: redis-mon-build
	docker run -rm -v `pwd`:/go/src/github.com/zenoss/redis-mon -e DUID=$(DUID) -e DGID=$(DGID) zenoss/redis-mon-build:$(VERSION) make rpm

# actual work
.PHONY: redis-mon-build
redis-mon-build:
	docker build -t zenoss/redis-mon-build:$(VERSION) hack


stage_pkg: redis-mon
	mkdir -p /tmp/pkgroot/usr/bin
	cp -v redis-mon /tmp/pkgroot/usr/bin

tgz: stage_pkg
	tar cvfz /tmp/$(FULL_NAME)-$(GIT_COMMIT).tgz -C /tmp/pkgroot/usr .
	chown $(DUID):$(DGID) /tmp/$(FULL_NAME)-$(GIT_COMMIT).tgz
	cp -p /tmp/$(FULL_NAME)-$(GIT_COMMIT).tgz .

deb: stage_pkg
	fpm \
		-n $(FULL_NAME) \
		-v $(VERSION) \
		-s dir \
		-t deb \
		-a x86_64 \
		-C $(PKGROOT) \
		-m $(MAINTAINER) \
		--description "$(DESCRIPTION)" \
		--deb-user root \
		--deb-group root \
		--license $(LICENSE) \
		--vendor $(VENDOR) \
		--url $(URL) \
		-f -p /tmp \
		.
	chown $(DUID):$(DGID) /tmp/*.deb
	cp -p /tmp/*.deb .


# Make an RPM
rpm: stage_pkg
	fpm \
		-n $(FULL_NAME) \
		-v $(VERSION) \
		-s dir \
		-t rpm \
		-a x86_64 \
		-C $(PKGROOT) \
		-m $(MAINTAINER) \
		--description "$(DESCRIPTION)" \
		--rpm-user root \
		--rpm-group root \
		--license $(LICENSE) \
		--vendor $(VENDOR) \
		--url $(URL) \
		-f -p /tmp \
		.
	chown $(DUID):$(DGID) /tmp/*.rpm
	cp -p /tmp/*.rpm .

clean:
	go clean
	rm -f *.deb
	rm -f *.rpm
	rm -f *.tgz
