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
    log "INFO" "Cleaning Go mod cache to avoid stale data..."
    go clean -modcache

    log "INFO" "Force updating protobuf dependency..."
    GOPROXY=direct go get -u github.com/ngdangkietswe/swe-protobuf-shared@latest || handle_error "Failed to update protobuf dependency"

    log "INFO" "Running go mod tidy..."
    go mod tidy || handle_error "Failed to tidy Go modules"

    if [ -d "vendor" ]; then
        log "INFO" "Clearing existing vendor directory..."
        check_writable "vendor"
        rm -rf vendor || handle_error "Failed to remove vendor directory"
    fi

    log "INFO" "Vendoring dependencies..."
    go mod vendor -v || handle_error "Failed to vendor Go modules"

    log "INFO" "Verifying if protobuf module is vendored..."
    if ! find vendor/github.com/ngdangkietswe/swe-protobuf-shared -type f > /dev/null 2>&1; then
        log "WARNING" "Protobuf module not found in vendor. Running go mod why..."
        go mod why github.com/ngdangkietswe/swe-protobuf-shared || true
    else
        log "INFO" "Protobuf module vendored successfully."
    fi
}

# Main execution
main() {
    log "INFO" "Starting protobuf update script..."

    command -v go >/dev/null 2>&1 || handle_error "Go is not installed"

    update_protobuf

    log "INFO" "Protobuf update script completed successfully!"
}

main
