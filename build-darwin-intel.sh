# Remove build directory if it exists
if [ -d build-darwin-intel ]; then
    rm -rf build-darwin-intel
fi

# Create build directory
mkdir build-darwin-intel

# Copy resources
mv apiteam-agent build-darwin-intel/apiteam-agent

# Copy files to build directory
cp -r targets/darwin build-darwin-intel

# Remove redis source code
rm -rf redis

# Build agent
GOOS=darwin GOARCH=amd64 go build -o build-darwin-intel/APITeam.app/Contents/MacOS/apiteam-agent -tags darwin-intel

# Recursively remove all gitkeep files
find . -name ".gitkeep" -type f -delete

# Clean up
rm build-darwin-intel/apiteam-agent
