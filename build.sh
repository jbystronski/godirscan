#!/bin/bash

SRC=cmd/main.go

CMD="go build -o"
TARGET_DIR=bin

mkdir -p bin/linux_amd64
env GOOS=linux GOARCH=amd64 $CMD $TARGET_DIR/linux_amd64/godirscan $SRC
mkdir -p bin/macos_amd64
env GOOS=darwin GOARCH=amd64 $CMD $TARGET_DIR/macos_amd64/godirscan $SRC
mkdir -p bin/linux_arm64
env GOOS=linux GOARCH=arm64 $CMD $TARGET_DIR/linux_arm64/godirscan $SRC
mkdir -p bin/win_amd64
env GOOS=windows GOARCH=amd64 $CMD $TARGET_DIR/win_amd64/godirscan.exe $SRC
