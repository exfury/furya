#!/bin/bash

set -e

GIT_TAG=$(git describe --tags)

echo "> Building $GIT_TAG..."

docker build . -t exfury/furya:$GIT_TAG