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

BIN_NAME := rpl
UNTAGGED_VERSION := v0.5.3
UNTAGGED_COMMIT := trunk

SHELL := /bin/bash
RUN_ARGS := --help
LOG_LEVEL := debug

GO_ENJIN_PKG := nil
BE_LOCAL_PATH := nil

GOPKG_KEYS ?= CDK CTK

CDK_GO_PACKAGE ?= github.com/go-curses/cdk
CDK_LOCAL_PATH ?= ../cdk

CTK_GO_PACKAGE ?= github.com/go-curses/ctk
CTK_LOCAL_PATH ?= ../ctk

GOPKG_KEYS += CL_RUN
CL_RUN_GO_PACKAGE ?= github.com/go-curses/corelibs/run
CL_RUN_LOCAL_PATH ?= ../corelibs/run

GOPKG_KEYS += CL_CHDIRS
CL_CHDIRS_GO_PACKAGE ?= github.com/go-curses/corelibs/chdirs
CL_CHDIRS_LOCAL_PATH ?= ../corelibs/chdirs

GOPKG_KEYS += CL_SPINNER
CL_SPINNER_GO_PACKAGE ?= github.com/go-curses/corelibs/spinner
CL_SPINNER_LOCAL_PATH ?= ../corelibs/spinner

GOPKG_KEYS += CL_MAPS
CL_MAPS_GO_PACKAGE ?= github.com/go-curses/corelibs/maps
CL_MAPS_LOCAL_PATH ?= ../corelibs/maps

GOPKG_KEYS += CL_STRINGS
CL_STRINGS_GO_PACKAGE ?= github.com/go-curses/corelibs/strings
CL_STRINGS_LOCAL_PATH ?= ../corelibs/strings

GOPKG_KEYS += CL_SLICES
CL_SLICES_GO_PACKAGE ?= github.com/go-curses/corelibs/slices
CL_SLICES_LOCAL_PATH ?= ../corelibs/slices

GOPKG_KEYS += CL_PATH
CL_PATH_GO_PACKAGE ?= github.com/go-curses/corelibs/path
CL_PATH_LOCAL_PATH ?= ../corelibs/path

GOPKG_KEYS += CL_REGEXPS
CL_REGEXPS_GO_PACKAGE ?= github.com/go-curses/corelibs/regexps
CL_REGEXPS_LOCAL_PATH ?= ../corelibs/regexps

GOPKG_KEYS += CL_MATHS
CL_MATHS_GO_PACKAGE ?= github.com/go-curses/corelibs/maths
CL_MATHS_LOCAL_PATH ?= ../corelibs/maths

GOPKG_KEYS += CL_CONVERT
CL_CONVERT_GO_PACKAGE ?= github.com/go-curses/corelibs/convert
CL_CONVERT_LOCAL_PATH ?= ../corelibs/convert

GOPKG_KEYS += CL_FILEWRITER
CL_FILEWRITER_GO_PACKAGE ?= github.com/go-curses/corelibs/filewriter
CL_FILEWRITER_LOCAL_PATH ?= ../corelibs/filewriter

CLEAN_FILES     ?= ${BIN_NAME} ${BIN_NAME}.*.* coverage.out pprof.*
DISTCLEAN_FILES ?=
REALCLEAN_FILES ?=

BUILD_VERSION_VAR := main.APP_VERSION
BUILD_RELEASE_VAR := main.APP_RELEASE

SRC_CMD_PATH := ./cmd/rpl

include Golang.cmd.mk
include Golang.def.mk
include Golang.cdk.mk
