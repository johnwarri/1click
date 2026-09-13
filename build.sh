#!/usr/bin/env bash
# Builds one standalone, dependency-free binary per OS/arch into ./dist.
# Run this on ANY of the three OSes — Go cross-compiles the others for you.
set -euo pipefail

APP=crossrun          # <-- change this to your app's name
OUT=dist

rm -rf "$OUT"
mkdir -p "$OUT"

# CGO off => fully static binaries that don't depend on anything installed.
# (If your logic ever calls C libraries, cross-compiling stops working and
#  you'll need to build each OS on its own machine — see the GitHub Action.)
export CGO_ENABLED=0

build() {
  local goos=$1 goarch=$2 ext=${3:-}
  local name="${APP}-${goos}-${goarch}${ext}"
  echo "building ${name}"
  GOOS="$goos" GOARCH="$goarch" \
    go build -trimpath -ldflags "-s -w" -o "${OUT}/${name}" .
}

build windows amd64 .exe
build windows arm64 .exe
build darwin  amd64          # Intel Macs
build darwin  arm64          # Apple Silicon Macs
build linux   amd64
build linux   arm64

# Friendlier names: darwin -> macos
mv "${OUT}/${APP}-darwin-amd64" "${OUT}/${APP}-macos-amd64"
mv "${OUT}/${APP}-darwin-arm64" "${OUT}/${APP}-macos-arm64"

# If lipo is available (i.e. we're on a Mac), stitch the two Mac builds into
# ONE universal binary that runs on both Intel and Apple Silicon. This is what
# the landing page hands to Mac visitors, so they never pick the wrong chip.
if command -v lipo >/dev/null 2>&1; then
  echo "building ${APP}-macos-universal (Intel + Apple Silicon)"
  lipo -create -output "${OUT}/${APP}-macos-universal" \
    "${OUT}/${APP}-macos-amd64" "${OUT}/${APP}-macos-arm64"
fi

echo
echo "done -> ${OUT}/"
ls -la "${OUT}"
