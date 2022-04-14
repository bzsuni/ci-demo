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

