#!/bin/bash

# Define the list of combinations
platforms=(
  "aix/ppc64"
  "darwin/amd64"
  "darwin/arm64"
  "dragonfly/amd64"
  "freebsd/386"
  "freebsd/amd64"
  "freebsd/arm"
  "freebsd/arm64"
  "illumos/amd64"
  "linux/386"
  "linux/amd64"
  "linux/arm"
  "linux/arm64"
  "linux/mips"
  "linux/mips64"
  "linux/mips64le"
  "linux/mipsle"
  "linux/ppc64"
  "linux/ppc64le"
  "linux/riscv64"
  "linux/s390x"
  "netbsd/386"
  "netbsd/amd64"
  "netbsd/arm"
  "netbsd/arm64"
  "openbsd/386"
  "openbsd/amd64"
  "openbsd/arm"
  "openbsd/arm64"
  "openbsd/mips64"
  "plan9/386"
  "plan9/amd64"
  "plan9/arm"
  "solaris/amd64"
  "windows/386"
  "windows/amd64"
  "windows/arm"
  "windows/arm64"
)

# Create the "dist" directory if it doesn't exist
mkdir -p dist

# Loop through each combination and build the binary
for platform in "${platforms[@]}"; do
  # Extract GOOS and GOARCH from the combination
  IFS='/' read -ra parts <<< "$platform"
  GOOS="${parts[0]}"
  GOARCH="${parts[1]}"

  # Set environment variables for the current platform
  export GOOS GOARCH

  # Build the binary
  go build .

  # Determine the binary extension based on the operating system
  if [ "$GOOS" = "windows" ]; then
    binary_ext=".exe"
  else
    binary_ext=""
  fi

  # Move the binary to the "dist" directory with the appropriate name
  mv "tmsnake$binary_ext" "dist/tmsnake-$GOOS-$GOARCH$binary_ext"
done

echo "Build completed successfully!"