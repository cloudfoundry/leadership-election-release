#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
source $DIR/shared

echo "-- BOSH create-release"
prep_src
create_release release.tgz

boshdir=$(mktemp -d)
echo "-- BOSH extract and compile in ${boshdir}"
pushd ${boshdir}
  tar -xzf ${DIR}/../release.tgz

  pushd packages
    mkdir -p leadership-election
    tar -xzf leadership-election.tgz -C leadership-election

    pushd leadership-election
      go build -buildmode=c-shared -mod=vendor -o leadership-election-agent ./src/cmd/leadership-election-agent
    popd
  popd
popd

cleanup
rm -f ${DIR}/../release.tgz
rm -rf ${boshdir}