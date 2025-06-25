#!/bin/sh

set -e

# Parse inputs from GitHub Action
COMMAND="$1"
CONFIG_FILE="$2"
ADDITIONAL_ARGS="$3"

echo "🚀 DeployAja CLI Action"
echo "Command: $COMMAND"
echo "Config File: $CONFIG_FILE"
echo "Additional Args: $ADDITIONAL_ARGS"

# Build the command
CMD_ARGS=""

if [ -n "$CONFIG_FILE" ] && [ "$CONFIG_FILE" != "null" ]; then
    CMD_ARGS="$CMD_ARGS -f $CONFIG_FILE"
fi

if [ -n "$ADDITIONAL_ARGS" ] && [ "$ADDITIONAL_ARGS" != "null" ]; then
    CMD_ARGS="$CMD_ARGS $ADDITIONAL_ARGS"
fi

# Execute the command
echo "Executing: aja $COMMAND $CMD_ARGS"
exec aja $COMMAND $CMD_ARGS 