# Remove build directory if it exists
if [ -d build-darwin-arm ]; then
    rm -rf build-darwin-arm
fi

# Create build directory
mkdir build-darwin-arm

# Copy files to build directory
cp -r targets/darwin/. build-darwin-arm

# Build agent
GOOS=darwin GOARCH=arm64 go build -o build-darwin-arm/APITeam\ Agent.app/Contents/MacOS/apiteam-agent -tags darwin_arm

# Recursively remove all gitkeep files
find . -name ".gitkeep" -type f -delete

# Use pckgbuild to create a package from the build directory
pkgbuild --component "build-darwin-arm/APITeam Agent.app" --version 0.1.0 --install-location /Applications build-darwin-arm/apiteam-agent.pkg

# Cleanup
rm -r build-darwin-arm/APITeam\ Agent.app