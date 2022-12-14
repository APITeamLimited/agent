# Remove build directory if it exists
if [ -d build-darwin-intel ]; then
    rm -rf build-darwin-intel
fi

# Create build directory
mkdir build-darwin-intel

# Copy files to build directory
cp -r targets/darwin/. build-darwin-intel

# Build agent
GOOS=darwin GOARCH=amd64 go build -o build-darwin-intel/APITeam\ Agent.app/Contents/MacOS/apiteam-agent -tags darwin_intel

# Recursively remove all gitkeep files
find . -name ".gitkeep" -type f -delete

# Use pckgbuild to create a package from the build directory
pkgbuild --root build-darwin-intel --identifier cloud.apiteam.agent --version 0.1.0 --install-location /Applications/APITeam\ Agent.app build-darwin-intel/APITeam\ Agent.pkg