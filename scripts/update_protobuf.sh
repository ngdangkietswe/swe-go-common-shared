#!/bin/bash

# update_protobuf.sh: Script to update the common protobuf dependency

set -euo pipefail

# Logging function
log() {
    local level="$1"
    local message="$2"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $message"
}

# Error handling function
handle_error() {
    log "ERROR" "$1"
    exit 1
}

# Check if a directory is writable
check_writable() {
    local dir="$1"
    if [ -d "$dir" ] && [ ! -w "$dir" ]; then
        handle_error "Directory $dir is not writable. Please check permissions."
    fi
}

# Update Go protobuf dependency
update_protobuf() {
    log "INFO" "Checking and updating protobuf dependency..."

    # Fetch the specific protobuf dependency
    GOPROXY=direct go get github.com/ngdangkietswe/swe-protobuf-shared || handle_error "Failed to update protobuf dependency"

    # Tidy modules
    log "INFO" "Tidying Go modules..."
    go mod tidy || handle_error "Failed to tidy Go modules"

    # Clear vendor directory if it exists
    if [ -d "vendor" ]; then
        log "INFO" "Clearing existing vendor directory..."
        check_writable "vendor"
        rm -rf vendor || handle_error "Failed to remove vendor directory"
    fi

    # Vendor dependencies
    log "INFO" "Vendoring dependencies..."
    go mod vendor -v || handle_error "Failed to vendor Go modules"

    log "INFO" "Protobuf dependency updated and vendored successfully!"
}

# Main execution
main() {
    log "INFO" "Starting protobuf update script..."

    # Check for required commands
    command -v go >/dev/null 2>&1 || handle_error "Go is not installed"

    update_protobuf

    log "INFO" "Protobuf update script completed successfully!"
}

main