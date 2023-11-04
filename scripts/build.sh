#!/bin/bash

set -e

GIT_TAG=$(git describe --tags)

echo "> Building $GIT_TAG..."

docker build . -t cosmoscontracts/furya:$GIT_TAG