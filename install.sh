#!/bin/bash

# Determine the installation location
INSTALL_PATH="/usr/local/bin"

# Build the Go binary
go build -o gsnake main.go
if [ $? -ne 0 ]; then
    echo "Failed to build the Go binary. Please ensure you have Go installed and try again."
    exit 1
fi

# Copy the binary to the installation location
sudo cp gsnake "$INSTALL_PATH/"
if [ $? -ne 0 ]; then
    echo "Failed to copy the binary to the installation location. Please ensure you have sufficient permissions."
    exit 1
fi

# Update file permissions
sudo chmod +x "$INSTALL_PATH/gsnake"
if [ $? -ne 0 ]; then
    echo "Failed to update file permissions. Please ensure you have sufficient permissions."
    exit 1
fi

# Determine the shell in use
current_shell=$(basename "$SHELL")

# Check if PATH already contains INSTALL_PATH
if echo "$PATH" | grep -q "$INSTALL_PATH"; then
    echo "PATH already contains $INSTALL_PATH."
else
    # Update the system's PATH
    if [ "$current_shell" == "bash" ]; then
        echo "export PATH=\"$INSTALL_PATH:\$PATH\"" >> ~/.bashrc
        source ~/.bashrc
    elif [ "$current_shell" == "zsh" ]; then
        echo "export PATH=\"$INSTALL_PATH:\$PATH\"" >> ~/.zshrc
        source ~/.zshrc
    elif [ "$current_shell" == "fish" ]; then
        echo "set -gx PATH $INSTALL_PATH \$PATH" >> ~/.config/fish/config.fish
        source ~/.config/fish/config.fish
    else
        echo "Unrecognized shell: $current_shell"
        echo "Installation is complete, but we were unable to update your PATH due to an unrecognized shell."
        echo "To use 'gsnake', please add the following to your shell's configuration file manually:"
        echo "PATH=$INSTALL_PATH:\$PATH"
        exit 0
    fi

    if [ $? -ne 0 ]; then
        echo "Installation is complete, but we were unable to update your shell configuration file. Please ensure it is writable."
        echo "To use 'gsnake', please add the following to your shell's configuration file manually:"
        echo "PATH=$INSTALL_PATH:\$PATH"
        exit 1
    fi
fi

echo "Installation complete. You can now run 'gsnake' to start the game."

