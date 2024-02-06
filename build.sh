#!/bin/bash

set -eo pipefail

cd front/post-app
./build.sh
cd ../..

cd back/post-service
./build.sh
cd ../..
