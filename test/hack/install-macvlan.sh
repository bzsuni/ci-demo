#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright Authors of Spider

# Copy 10-macvlan.conflist to kind-node
echo "CLUSTER_NAME: $CLUSTER_NAME"
NODES=($(docker ps | grep $CLUSTER_NAME | awk '{print $1}'))
for node in ${NODES[@]}
do
  echo "docker cp $(pwd)/../tools/yamls/10-macvlan.conflist $node:/etc/cni/net.d"
  docker cp ../yaml/macvlan/10-macvlan.conflist $node:/etc/cni/net.d
done

sleep 5

# Install macvlan + whereabouts
# Prepare Image
IMAGE_EXIST=$(docker images | grep ghcr.io/k8snetworkplumbingwg/whereabouts:latest-amd64)
if test -z "$IMAGE_EXIST"; then
  docker pull ghcr.io/k8snetworkplumbingwg/whereabouts:latest-amd64
fi
kind load docker-image ghcr.io/k8snetworkplumbingwg/whereabouts:latest-amd64 --name $CLUSTER_NAME

# Install whereabouts
export PATH=$PATH:/usr/local/git/bin
echo $PATH
git -C /tmp clone https://github.com/k8snetworkplumbingwg/whereabouts
if [ -d "/tmp/whereabouts" ]; then
  kubectl apply \
      -f /tmp/whereabouts/doc/crds/daemonset-install.yaml \
      -f /tmp/whereabouts/doc/crds/whereabouts.cni.cncf.io_ippools.yaml \
      -f /tmp/whereabouts/doc/crds/whereabouts.cni.cncf.io_overlappingrangeipreservations.yaml \
      -f /tmp/whereabouts/doc/crds/ip-reconciler-job.yaml
  rm -rf /tmp/whereabouts
else
  echo "Install whereabouts failed,Please try again"
fi
