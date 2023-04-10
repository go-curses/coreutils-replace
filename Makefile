#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

# Copyright (c) 2023  The Go-Curses Authors
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

#: uncomment to echo instead of execute
#CMD=echo

-include .env
#export

.PHONY: all help run
.PHONY: clean distclean realclean
.PHONY: local unlocal tidy be-update
.PHONY: debug build build-all build-amd64 build-arm64
.PHONY: release release-all release-amd64 release-arm64
.PHONY: install install-autocomplete
.PHONY: install-shortlinks install-shortestlinks

BIN_NAME ?= rpl
UNTAGGED_VERSION ?= v0.2.3
UNTAGGED_COMMIT ?= 0000000000

SHELL := /bin/bash
LOG_LEVEL := debug

RUN_ARGS := --help

GO_ENJIN_PKG := nil
BE_LOCAL_PATH := nil

GOPKG_KEYS ?= CDK CTK CLD CLE CLN CLP

CDK_GO_PACKAGE ?= github.com/go-curses/cdk
CDK_LOCAL_PATH ?= ../cdk

CTK_GO_PACKAGE ?= github.com/go-curses/ctk
CTK_LOCAL_PATH ?= ../ctk

CLD_GO_PACKAGE ?= github.com/go-curses/corelibs/diff
CLD_LOCAL_PATH ?= ../corelibs/diff

CLE_GO_PACKAGE ?= github.com/go-curses/corelibs/errors
CLE_LOCAL_PATH ?= ../corelibs/errors

CLN_GO_PACKAGE ?= github.com/go-curses/corelibs/notify
CLN_LOCAL_PATH ?= ../corelibs/notify

CLP_GO_PACKAGE ?= github.com/go-curses/corelibs/path
CLP_LOCAL_PATH ?= ../corelibs/path

CLEAN_FILES     ?= ${BIN_NAME} ${BIN_NAME}.*.* coverage.out pprof.*
DISTCLEAN_FILES ?=
REALCLEAN_FILES ?=

BUILD_VERSION_VAR := main.APP_VERSION
BUILD_RELEASE_VAR := main.APP_RELEASE

include Golang.cmd.mk

help:
	@echo "usage: make [target]"
	@echo
	@echo "qa targets:"
	@echo "  vet         - run go vet command"
	@echo "  test        - perform all available tests"
	@echo "  cover       - perform all available tests with coverage report"
	@echo
	@echo "cleanup targets:"
	@echo "  clean       - cleans package and built files"
	@echo "  distclean   - clean and removes extraneous files"
	@echo
	@echo "build targets:"
	@echo "  debug       - build debug ${BUILD_NAME}"
	@echo "  build       - build clean ${BUILD_NAME}"
	@echo "  release     - build clean and compress ${BUILD_NAME}"
	@echo
	@echo "cross-build targets:"
	@echo "  build-all     - both build-arm64 and build-amd64"
	@echo "  build-arm64   - build clean ${BIN_NAME}.${BUILD_OS}.arm64"
	@echo "  build-amd64   - build clean ${BIN_NAME}.${BUILD_OS}.amd64"
	@echo "  release-all   - both release-arm64 and release-amd64"
	@echo "  release-arm64 - build clean and compress ${BIN_NAME}.${BUILD_OS}.arm64"
	@echo "  release-amd64 - build clean and compress ${BIN_NAME}.${BUILD_OS}.amd64"
	@echo
	@echo "install targets:"
	@echo "  install       - installs ${BUILD_NAME} to ${DESTDIR}${prefix}/bin/${BIN_NAME}"
	@echo "  install-arm64 - installs ${BIN_NAME}.${BUILD_OS}.arm64 to ${DESTDIR}${prefix}/bin/${BIN_NAME}"
	@echo "  install-amd64 - installs ${BIN_NAME}.${BUILD_OS}.amd64 to ${DESTDIR}${prefix}/bin/${BIN_NAME}"
	@echo
	@echo "run targets:"
	@echo "  run         - run the dev build (sanely handle crashes)"
	@echo "  profile.cpu - run the dev build and profile CPU"
	@echo "  profile.mem - run the dev build and profile memory"
	@echo
	@echo "go helpers:"
	@echo "  local       - add go.mod local GOPKG_KEYS replacements"
	@echo "  unlocal     - remove go.mod local GOPKG_KEYS replacements"
	@echo "  generate    - run go generate ./..."
	@echo "  be-update   - get latest GOPKG_KEYS dependencies"
	@echo
	@echo "Notes:"
	@echo "  GOPKG_KEYS are go packages managed by this Makefile."
	@echo "  The following are the available GOPKG_KEYS:" \
		$(if ${GOPKG_KEYS},$(foreach key,${GOPKG_KEYS},; echo "    $(key): $($(key)_GO_PACKAGE) ($($(key)_LOCAL_PATH))"))

clean:
	@$(call __clean,${CLEAN_FILES})

distclean: clean
	@$(call __clean,${DISTCLEAN_FILES})

realclean: distclean
	@$(call __clean,${REALCLEAN_FILES})

debug: BUILD_VERSION=$(call __tag_ver)
debug: BUILD_RELEASE=$(call __rel_ver)
debug: TRIM_PATHS=$(call __go_trim_path)
debug: __golang
	@$(call __go_build_debug,"${BUILD_NAME}",${BUILD_OS},${BUILD_ARCH},./cmd/rpl)
	@${SHASUM_CMD} "${BUILD_NAME}"

build: BUILD_VERSION=$(call __tag_ver)
build: BUILD_RELEASE=$(call __rel_ver)
build: TRIM_PATHS=$(call __go_trim_path)
build: __golang
	@$(call __go_build_release,"${BUILD_NAME}",${BUILD_OS},${BUILD_ARCH},./cmd/rpl)
	@${SHASUM_CMD} "${BUILD_NAME}"

build-amd64: BUILD_VERSION=$(call __tag_ver)
build-amd64: BUILD_RELEASE=$(call __rel_ver)
build-amd64: TRIM_PATHS=$(call __go_trim_path)
build-amd64: __golang
	@$(call __go_build_release,"${BIN_NAME}.${BUILD_OS}.amd64",${BUILD_OS},amd64,./cmd/rpl)
	@${SHASUM_CMD} "${BIN_NAME}.${BUILD_OS}.amd64"

build-arm64: BUILD_VERSION=$(call __tag_ver)
build-arm64: BUILD_RELEASE=$(call __rel_ver)
build-arm64: TRIM_PATHS=$(call __go_trim_path)
build-arm64: __golang
	@$(call __go_build_release,"${BIN_NAME}.${BUILD_OS}.arm64",${BUILD_OS},arm64,./cmd/rpl)
	@${SHASUM_CMD} "${BIN_NAME}.${BUILD_OS}.arm64"

build-all: build-amd64 build-arm64

release: build
	@$(call __upx_build,"${BUILD_NAME}")

release-arm64: build-arm64
	@$(call __upx_build,"${BIN_NAME}.${BUILD_OS}.arm64")

release-amd64: build-amd64
	@$(call __upx_build,"${BIN_NAME}.${BUILD_OS}.amd64")

release-all: release-amd64 release-arm64

install:
	@if [ -f "${BUILD_NAME}" ]; then \
		echo "# ${BUILD_NAME} present"; \
		$(call __install_exe,"${BUILD_NAME}","${INSTALL_BIN_PATH}/${BIN_NAME}"); \
	else \
		echo "error: missing ${BUILD_NAME} binary" 1>&2; \
		false; \
	fi

install-arm64:
	@if [ -f "${BIN_NAME}.${BUILD_OS}.arm64" ]; then \
		echo "# ${BIN_NAME}.${BUILD_OS}.arm64 present"; \
		$(call __install_exe,"${BIN_NAME}.${BUILD_OS}.arm64","${INSTALL_BIN_PATH}/${BIN_NAME}"); \
	else \
		echo "error: missing ${BIN_NAME}.${BUILD_OS}.arm64 binary" 1>&2; \
		false; \
	fi

install-amd64:
	@if [ -f "${BIN_NAME}.${BUILD_OS}.amd64" ]; then \
		echo "# ${BIN_NAME}.${BUILD_OS}.amd64 present"; \
		$(call __install_exe,"${BIN_NAME}.${BUILD_OS}.amd64","${INSTALL_BIN_PATH}/${BIN_NAME}"); \
	else \
		echo "error: missing ${BIN_NAME}.${BUILD_OS}.amd64 binary" 1>&2; \
		false; \
	fi

run: export GO_CDK_LOG_FILE=./${BUILD_NAME}.cdk.log
run: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
run: export GO_CDK_LOG_FULL_PATHS=true
run:
	@if [ -f ${BUILD_NAME} ]; \
	then \
		echo "# running: ${BUILD_NAME} ${RUN_ARGS}"; \
		( ./${BUILD_NAME} ${RUN_ARGS} ) 2>> ${GO_CDK_LOG_FILE}; \
		if [ $$? -ne 0 ]; \
		then \
			stty sane; echo ""; \
			echo "# ${BIN_NAME} crashed, see: ./${BIN_NAME}.cdk.log"; \
			read -p "# Press <Enter> to reset terminal, <Ctrl+C> to cancel" RESP; \
			reset; \
			echo "# ${BIN_NAME} crashed, terminal reset, see: ./${BIN_NAME}.cdk.log"; \
		else \
			echo "# ${BIN_NAME} exited normally."; \
		fi; \
	fi

profile.cpu: export GO_CDK_LOG_FILE=./${BIN_NAME}.cdk.log
profile.cpu: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
profile.cpu: export GO_CDK_LOG_FULL_PATHS=true
profile.cpu: export GO_CDK_PROFILE_PATH=/tmp/${BIN_NAME}.cdk.pprof
profile.cpu: export GO_CDK_PROFILE=cpu
profile.cpu: debug
	@mkdir -v /tmp/${BIN_NAME}.cdk.pprof 2>/dev/null || true
	@if [ -f ${BIN_NAME} ]; \
		then \
			./${BIN_NAME} && \
			if [ -f /tmp/${BIN_NAME}.cdk.pprof/cpu.pprof ]; \
			then \
				read -p "# Press enter to open a pprof instance" JUNK \
				&& go tool pprof -http=:8080 /tmp/${BIN_NAME}.cdk.pprof/cpu.pprof ; \
			else \
				echo "# missing /tmp/${BIN_NAME}.cdk.pprof/cpu.pprof"; \
			fi ; \
		fi

profile.mem: export GO_CDK_LOG_FILE=./${BIN_NAME}.log
profile.mem: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
profile.mem: export GO_CDK_LOG_FULL_PATHS=true
profile.mem: export GO_CDK_PROFILE_PATH=/tmp/${BIN_NAME}.cdk.pprof
profile.mem: export GO_CDK_PROFILE=mem
profile.mem: debug
	@mkdir -v /tmp/${BIN_NAME}.cdk.pprof 2>/dev/null || true
	@if [ -f ${BIN_NAME} ]; \
		then \
			./${BIN_NAME} && \
			if [ -f /tmp/${BIN_NAME}.cdk.pprof/mem.pprof ]; \
			then \
				read -p "# Press enter to open a pprof instance" JUNK \
				&& go tool pprof -http=:8080 /tmp/${BIN_NAME}.cdk.pprof/mem.pprof; \
			else \
				echo "# missing /tmp/${BIN_NAME}.cdk.pprof/mem.pprof"; \
			fi ; \
		fi

tidy: __tidy

be-update: __be_update

local: __local

unlocal: __unlocal

vet: __vet

test: __test

cover: __cover

generate: __generate
