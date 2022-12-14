# Remove build directories if they exist
if [ -d pkg-build ]; then
    rm -rf pkg-build
fi

if [ -d build-agent-linux ]; then
    rm -rf build-agent-linux
fi

# Create build directory
mkdir build-agent-linux

# Build agent and package it for deb
GOOS=linux GOARCH=amd64 go build -o build-agent-linux/apiteam-agent -tags linux
go-bin-deb generate --output=build-agent-linux/apiteam-agent.deb --file=targets/linux/deb.json

# Add apiteam-agent to a tar.gz file
mkdir build-agent-linux/apiteam
cp build-agent-linux/apiteam-agent build-agent-linux/apiteam/apiteam-agent

# Create a tar.gz file
tar -czvf build-agent-linux/apiteam-agent.tar.gz -C build-agent-linux/apiteam .

# Clean up
rm -rf build-agent-linux/apiteam
rm build-agent-linux/apiteam-agent
