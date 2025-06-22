#!/bin/sh

set -e

# Parse inputs from GitHub Action
COMMAND="$1"
CONFIG_FILE="$2"
ENVIRONMENT="$3"
PROJECT_NAME="$4"
ADDITIONAL_ARGS="$5"

echo "🚀 DeployAja CLI Action"
echo "Command: $COMMAND"
echo "Config File: $CONFIG_FILE"
echo "Environment: $ENVIRONMENT"
echo "Project: $PROJECT_NAME"

# Build the command
CMD_ARGS=""

if [ -n "$ENVIRONMENT" ] && [ "$ENVIRONMENT" != "null" ]; then
    CMD_ARGS="$CMD_ARGS --env $ENVIRONMENT"
fi

if [ -n "$PROJECT_NAME" ] && [ "$PROJECT_NAME" != "null" ]; then
    CMD_ARGS="$CMD_ARGS --project $PROJECT_NAME"
fi

if [ -n "$ADDITIONAL_ARGS" ] && [ "$ADDITIONAL_ARGS" != "null" ]; then
    CMD_ARGS="$CMD_ARGS $ADDITIONAL_ARGS"
fi

# Execute the command
echo "Executing: aja $COMMAND $CMD_ARGS"
exec aja $COMMAND $CMD_ARGS 