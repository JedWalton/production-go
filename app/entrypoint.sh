#!/bin/sh

# Check if RUN_TESTS is set to true
if [ "$RUN_TESTS" = "true" ]; then
    echo "Running integration tests..."
else
    echo "Starting application..."
    exec /go/bin/production-go
fi

