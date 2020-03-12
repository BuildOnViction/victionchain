#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
ethdir="$workspace/src/github.com/chancoin"
if [ ! -L "$ethdir/chancoin" ]; then
    mkdir -p "$ethdir"
    cd "$ethdir"
    ln -s ../../../../../. chancoin
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$ethdir/chancoin"
PWD="$ethdir/chancoin"

# Launch the arguments with the configured environment.
exec "$@"
