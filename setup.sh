#!/bin/bash

set -e

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Go is not installed. Installing Go..."
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        GO_VERSION="1.20.13"
        wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
        rm go${GO_VERSION}.linux-amd64.tar.gz
        export PATH=$PATH:/usr/local/go/bin
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        GO_VERSION="1.20.13"
        curl -LO https://go.dev/dl/go${GO_VERSION}.darwin-amd64.pkg
        sudo installer -pkg go${GO_VERSION}.darwin-amd64.pkg -target /
        rm go${GO_VERSION}.darwin-amd64.pkg
        export PATH=$PATH:/usr/local/go/bin
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
    else
        echo "Please install Go manually: https://go.dev/dl/"
        exit 1
    fi
else
    echo "Go is already installed."
fi

# Clone the repository
if [ ! -d "cli" ]; then
    git clone https://github.com/deployaja/deployaja-cli.git
fi

cd cli

# Build the binary
go build -o aja main.go

# Copy binary to /usr/local/bin
sudo cp aja /usr/local/bin/


echo "DeployAja CLI installed successfully! Run 'deployaja --help' to get started."
