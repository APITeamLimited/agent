# Remove build directory if it exists
if [ -d build-agent-linux ]; then
    rm -rf build-agent-linux
fi

# Create build directory
mkdir build-agent-linux
cd build-agent-linux

# Build agent
GOOS=linux GOARCH=amd64 go build -o build-agent-linux/apiteam-agent

# Copy files to build directory
cp targets/linux/redis-server build-agent-linux/redis-server
cp targets/linux/apiteam-agent.desktop build-agent-linux/apiteam-agent.desktop
cp targets/linux/snapcraft.yaml build-agent-linux/snapcraft.yaml
cp targets/linux/run.sh build-agent-linux/run.sh
cp apiteam-logo.png build-agent-linux/apiteam-logo.png

# Build snap
cd build-agent-linux
snapcraft

# Clean up
cd ..
#rm build-agent-linux/apiteam-agent
#rm build-agent-linux/redis-server
#rm build-agent-linux/apiteam-agent.desktop
#rm build-agent-linux/snapcraft.yaml
#rm build-agent-linux/run.sh
#rm build-agent-linux/apiteam-logo.png
