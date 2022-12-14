# Remove build directories if they exist
if [ -d pkg-build ]; then
    rm -rf pkg-build
fi

if [ -d build-linux ]; then
    rm -rf build-linux
fi

# Create build directory
mkdir build-linux

# Build agent and package it for deb
GOOS=linux GOARCH=amd64 go build -o build-linux/apiteam-agent -tags linux
go-bin-deb generate --output=build-linux/apiteam-agent.deb --file=targets/linux/deb.json

# Add apiteam-agent to a tar.gz file
mkdir build-linux/apiteam
cp build-linux/apiteam-agent build-linux/apiteam/apiteam-agent

# Create a tar.gz file
tar -czvf build-linux/apiteam-agent.tar.gz -C build-linux/apiteam .

# Clean up
rm -rf build-linux/apiteam
rm build-linux/apiteam-agent
