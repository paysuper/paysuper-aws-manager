#!/usr/bin/env sh

if [ -n "$1" ] && [ ${0:0:4} = "/bin" ]; then
  ROOT_DIR=$1/..
else
  ROOT_DIR="$( cd "$( dirname "$0" )" && pwd )/.."
fi

mockery -name=AwsManagerInterface -dir=${ROOT_DIR} -output ${ROOT_DIR}/pkg/mocks