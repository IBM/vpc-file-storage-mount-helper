#
# Copyright 2024 IBM Corp.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

export LINT_VERSION="1.43.0"
GOPACKAGES=$(shell go list ./... | grep -v /vendor/ | grep -v /tests)

NAME := mount-helper-container
APP_VERSION := 0.1.1
MAINTAINER := "IKS Storage"
DESCRIPTION := "IBM Mount-helper-container service"
LICENSE := "IBM"
DEB_ARCH := all
RPM_ARCH := x86_64
RPM_RELEASE_NUM := 1

BUILD_DIR := $(NAME)-$(APP_VERSION)
BIN_DIR := bin
INSTALL_TAR_FILE := "$(NAME)-latest.tar.gz"
CHECKSUM_FILE := "$(INSTALL_TAR_FILE).sha256"

DEBIAN_CONTROL := $(BUILD_DIR)/DEBIAN/control
REDHAT_SPEC := $(BUILD_DIR)/red-hat.spec

deb-build:
deb-build: build-linux
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)/DEBIAN
	mkdir -p ${BUILD_DIR}/etc/systemd/system/ ${BUILD_DIR}/usr/local/bin/ $(BUILD_DIR)/etc/ibmcloud

	cp install/mount-helper-container.service $(BUILD_DIR)/etc/systemd/system/
	cp install/share.conf $(BUILD_DIR)/etc/ibmcloud/
	cp ${BIN_DIR}/mount-helper-container-server $(BUILD_DIR)/usr/local/bin/
	cp install/postinst.sh $(BUILD_DIR)/DEBIAN/postinst
	cp install/prerm.sh $(BUILD_DIR)/DEBIAN/prerm

	echo "Package: $(BUILD_DIR)" > $(DEBIAN_CONTROL)
	echo "Version: $(APP_VERSION)" >> $(DEBIAN_CONTROL)
	echo "Maintainer: $(MAINTAINER)" >> $(DEBIAN_CONTROL)
	echo "Architecture: $(DEB_ARCH)" >> $(DEBIAN_CONTROL)
	echo "Description: $(DESCRIPTION)" >> $(DEBIAN_CONTROL)
	echo "Depends: mount.ibmshare-0.1.7" >> $(DEBIAN_CONTROL)

	dpkg-deb --build $(BUILD_DIR)
	rm -rf $(BUILD_DIR)

rpm-build:
rpm-build: build-linux
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)/
	mkdir -p $(BUILD_DIR)/rpm/SOURCES/usr/local/bin $(BUILD_DIR)/rpm/SOURCES/etc/systemd/system $(BUILD_DIR)/rpm/SOURCES/etc/ibmcloud
	cp install/mount-helper-container.service $(BUILD_DIR)/rpm/SOURCES/etc/systemd/system
	cp install/share.conf $(BUILD_DIR)/rpm/SOURCES/etc/ibmcloud
	cp ${BIN_DIR}/mount-helper-container-server $(BUILD_DIR)/rpm/SOURCES/usr/local/bin

	echo "Name: $(NAME)" > $(REDHAT_SPEC)
	echo "Version: $(APP_VERSION)" >> $(REDHAT_SPEC)
	echo "Release: $(RPM_RELEASE_NUM)" >> $(REDHAT_SPEC)
	echo "Summary: $(DESCRIPTION)" >> $(REDHAT_SPEC)
	echo "License: $(LICENSE)" >>  $(REDHAT_SPEC)
	echo "BuildArch: $(RPM_ARCH)" >>  $(REDHAT_SPEC)
	echo "Requires: mount.ibmshare = 0.1.7" >> $(REDHAT_SPEC)
	echo "%define _rpmfilename $(NAME)-$(APP_VERSION).rpm" >> $(REDHAT_SPEC)
	echo "%build" >> $(REDHAT_SPEC)

	echo "%install" >> $(REDHAT_SPEC)
	echo "rm -rf %{buildroot}" >> $(REDHAT_SPEC)
	echo "mkdir -p %{buildroot}/etc/systemd/system/" >> $(REDHAT_SPEC)
	echo "mkdir -p %{buildroot}/usr/local/bin/" >> $(REDHAT_SPEC)
	echo "mkdir -p %{buildroot}/etc/ibmcloud/" >> $(REDHAT_SPEC)
	echo "cp -r %{_sourcedir}/* %{buildroot}/" >> $(REDHAT_SPEC)

	echo "%description" >> $(REDHAT_SPEC)
	echo "This package provides the IBM Mount Share helper container service."

	echo "%files" >> $(REDHAT_SPEC)
	echo "/etc/systemd/system/mount-helper-container.service" >> $(REDHAT_SPEC)
	echo "/usr/local/bin/mount-helper-container-server" >> $(REDHAT_SPEC)
	echo "/etc/ibmcloud/share.conf" >> $(REDHAT_SPEC)

	echo "%post" >> $(REDHAT_SPEC)
	echo "systemctl enable mount-helper-container.service" >> $(REDHAT_SPEC)
	echo "systemctl start mount-helper-container.service" >> $(REDHAT_SPEC)

	echo "%preun" >> $(REDHAT_SPEC)
	echo "systemctl stop mount-helper-container.service || true" >> $(REDHAT_SPEC)
	echo "systemctl disable mount-helper-container.service || true" >> $(REDHAT_SPEC)

	rpmbuild -ba --build-in-place --define "_topdir $(PWD)/$(BUILD_DIR)/rpm" $(REDHAT_SPEC)
	cp  $(BUILD_DIR)/rpm/RPMS/* ./
	rm -rf $(BUILD_DIR)

tar-package:
tar-package:
	mkdir packages
	cp ${NAME}-* packages/
	cd packages && tar -czvf  ../$(INSTALL_TAR_FILE)  *
	sha256sum ./$(INSTALL_TAR_FILE) > $(CHECKSUM_FILE)
	@printf "Production - Install package created ok: $(INSTALL_TAR_FILE)\n"

packages:
packages: deb-build rpm-build tar-package clean

build:
	CGO_ENABLED=0 GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) go build -mod mod -o ${BIN_DIR}/mount-helper-container-server ./server/server.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod mod -o ${BIN_DIR}/mount-helper-container-server -ldflags "-s -w" -a ./server/server.go

clean:
	rm -rf ${BIN_DIR}
	rm *.deb *.rpm
	rm -rf packages

startservice:
startservice: build-linux
	systemctl daemon-reload
	cp ./mount-helper-container.service /etc/systemd/system/
	cp ./${BIN_DIR}/mount-helper-container-server /usr/local/bin/
	service mount-helper-container start
	systemctl enable mount-helper-container
	systemctl start mount-helper-container

test:
	go test -v -timeout 1800s -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o=cover.html

ut-coverage: test
	@./scripts/coverage.sh
	rm cover.html cover.out
	go mod tidy
