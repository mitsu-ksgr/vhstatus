#!/usr/bin/env sh
#------------------------------------------------------------------------------
#
# Build on docker.
#
#------------------------------------------------------------------------------
set -u
umask 0022
export LC_ALL=C
readonly SCRIPT_NAME=$(basename $0)

readonly TAG="vhstatus_vhstatus:latest"

cd $(dirname $0)/../

docker build -t $TAG .
docker run --rm -v $(pwd):/go/src/github.com/mitsu-ksgr/vhstatus $TAG

echo "build succeeded: ./vhstatus-server"

