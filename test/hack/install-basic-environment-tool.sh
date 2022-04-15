#!/usr/bin/env bash

function msg() {
   if [[ $# -ne 1 ]]; then echo "[func msg] one arg needed"; exit 1; fi
    echo -e "\033[35m $1 \033[0m"
}
function err() {
   if [[ $# -ne 1 ]]; then echo "[func err] one arg needed"; exit 1; fi
   echo -e "\033[31m $1 \033[0m"
}
function succ() {
   if [[ $# -ne 1 ]]; then echo "[func succ] one arg needed"; exit 1; fi
   echo -e "\033[32m $1 \033[0m"
}

# install go
msg "## install go"
if ! $(go version > /dev/null 2>&1); then
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
     sudo ln -s /usr/local/go/bin/go /usr/local/bin/go
  fi
  if ! $(go version > /dev/null 2>&1); then err "failed install go"; exit 1; fi
  # set go env
  go env -w GOPROXY=https://goproxy.cn,direct
  go env -w GO111MODULE=on
else
  msg "go has already been installed"
fi

# install docker
msg "## install docker"

if ! $(which docker > /dev/null 2>&1); then
  sudo yum -y install docker
  sudo systemctl start docker
  if ! $(which docker > /dev/null 2>&1); then
    err "failed install docker";
    exit 1;
  else
    succ "install docker succeed"
  fi
  sudo systemctl enable docker

  sudo groupadd docker
  sudo gpasswd -a $USER docker
  sudo newgrp docker
else
  msg "docker has already been installed"
fi
# resolve runner run docker permission denied
sudo chmod a+rwx /var/run/docker.sock

# install git nmap jq
needs="git nmap jq"
for need in $needs; do
  msg "## install $need"
  if ! $($need --version > /dev/null 2>&1); then
    sudo yum -y install $need
    if ! $($need --version > /dev/null 2>&1); then
      err "failed install $need"; exit 1;
    else
      succ "install $need succeed"
    fi
  else
    msg "$need has already been installed"
  fi
done