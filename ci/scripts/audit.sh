#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-net
  make audit
popd