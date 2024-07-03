#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

# Variables
BASEDIR=$(dirname $(realpath "$0"))
echo "Base Dir: $BASEDIR"
TARGET="bon-voyage-agent"
SRC_DIR="$BASEDIR/src"
PLUGIN_SRC_DIR="$BASEDIR/src/plugins"
BUILD_DIR="$BASEDIR/build"
PLUGIN_BUILD_DIR="$BUILD_DIR/plugins"

# Check if Go is installed
if ! [ -x "$(command -v go)" ]; then
  echo "Error: Go is not installed." >&2
  exit 1
fi

# Determine the operating system
OS="$(uname -s)"
case "$OS" in
    Linux*)     GOOS="linux";;
    Darwin*)    GOOS="darwin";;
    *)          echo "Unsupported OS: $OS"; exit 1;;
esac

# Clean previous builds
echo "Cleaning previous builds..."
rm -f $BUILD_DIR/$TARGET

# Build the agent
echo "Building agent ..."
cd $SRC_DIR

GIT_HASH=$(git rev-parse HEAD)
go build -o $BUILD_DIR/$TARGET -ldflags="-X main.Commit=$GIT_HASH"
cd $BASEDIR

# Build the plugins
cd $PLUGIN_SRC_DIR
./build_plugins.sh $PLUGIN_BUILD_DIR
cd $BASEDIR
echo "Build finished. See README.md for usage."
