#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
source $DIR/shared

echo "-- BOSH create-release"
prep_src
create_release release.tgz
cleanup