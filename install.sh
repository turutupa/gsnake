#!/bin/bash

# Determine the installation location
INSTALL_PATH="/usr/local/bin"

# Build the Go binary
go build -o gsnake main.go

# Copy the binary to the installation location
sudo cp gsnake "$INSTALL_PATH/"

# Update file permissions
sudo chmod +x "$INSTALL_PATH/gsnake"

# Update the system's PATH
echo "export PATH=\"$INSTALL_PATH:\$PATH\"" >> ~/.bashrc
source ~/.bashrc

echo "Installation complete. You can now run 'gsnake' to start the game."

