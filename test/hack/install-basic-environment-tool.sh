#!/usr/bin/env bash

function msg() {
   if [[ $# -ne 1 ]]; then echo "[func msg] one arg needed"; exit 1; fi
    echo -e "\033[35m $1 \033[0m"
}
# install go
msg "## install go"
if ! $(which go > /dev/null 2>&1); then
  echo need to install go
fi

# install docker
msg "## install docker"
if ! $(which docker > /dev/null 2>&1); then
  echo need to install go
fi
# install git
msg "## install git"
if ! $(which git > /dev/null 2>&1); then
  echo need to install go
fi