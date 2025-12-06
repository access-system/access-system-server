#!/usr/bin/env bash

if [ "$1" = "help" ]; then
  echo "Usage: ./run.sh [profile]"
  echo "Profiles:"
  echo "  dev  - Run the application in development mode"
  echo "  test - Run the application in test mode"
  echo "  help - Show this help message"

  exit 0
fi

if [ "$1" = "dev" ]; then
  docker-compose --profile dev up --build

  exit 0
fi

if [ "$1" = "clean-dev" ]; then
  docker rm access-system-postgres
  docker volume rm access-system-server_pgdata
  docker-compose --profile dev up --build

  exit 0
fi

if [ "$1" = "test" ]; then
  # Check if internal/mocks exists
  if [ ! -d "internal/mocks" ]; then
    echo "internal/mocks not found, running go generate ./..."
    go generate ./...
    if [ $? -ne 0 ]; then
      echo "Error: go generate failed. Aborting tests."
      exit 4
    fi
  fi

  echo "Run tests..."
  docker-compose --profile test up --build -d
  if [ $? -ne 0 ]; then
    echo "Error: Failed to start test containers."
    docker-compose --profile test down
    exit 2
  fi

  # Wait for the test container to finish
  docker wait access-system-server-test
  WAIT_EXIT_CODE=$?
  if [ $WAIT_EXIT_CODE -ne 0 ]; then
    echo "Error: Test container did not start or failed unexpectedly."
    docker-compose --profile test down
    exit 3
  fi

  echo "Test container logs:"
  docker logs access-system-server-test

  echo "Tests finished, cleaning up..."
  docker-compose --profile test down
  if [ $? -ne 0 ]; then
    echo "Warning: Cleanup failed. Please check Docker resources manually."
    exit 4
  fi

  exit 0
fi

echo "No profile specified. Use 'help' for usage information."
exit 1
