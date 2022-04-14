#!/usr/bin/env bash

function msg() {
   if [[ $# -ne 1 ]]; then echo "[func msg] one arg needed"; exit 1; fi
    echo -e "\033[35m $1 \033[0m"
}
function err() {
   if [[ $# -ne 1 ]]; then echo "[func err] one arg needed"; exit 1; fi
   echo -e "\033[31m $1 \033[0m"
}
# install go
msg "## install go"
if ! $(which go > /dev/null 2>&1); then
  sudo wget "https://golang.google.cn/dl/go1.18.linux-amd64.tar.gz"
  sudo tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz
  sudo rm -rf go1.18.linux-amd64.tar.gz

  # add go path
  if ! $(cat $HOME/.bashrc | grep GOROOT > /dev/null 2>&1); then
    sudo echo -e "export GOROOT=/usr/local/go" >> $HOME/.bashrc
  fi
  if ! $(cat $HOME/.bashrc | grep GOPATH > /dev/null 2>&1); then
    sudo echo -e "export GOPATH=$HOME/go" >> $HOME/.bashrc
  fi
  if ! $(cat $HOME/.bashrc | grep GOBIN > /dev/null 2>&1); then
    sudo echo -e "export GOBIN=$HOME/go/bin" >> $HOME/.bashrc
  fi
  if ! $(cat $HOME/.bashrc | grep -w PATH | grep GOPATH > /dev/null 2>&1); then
    sudo echo -e 'export PATH=$GOPATH:$GOBIN:$GOROOT/bin:$PATH' >> $HOME/.bashrc
  fi
st
  # set go env
  sudo go env -w GOPROXY=https://goproxy.cn,direct
  sudo go env -w GO111MODULE=on
fi

if ! $(go version); then err "err install go"; exit 1; fi

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