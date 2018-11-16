#!/bin/bash

export ETCD_AUTO_COMPACTION_RETENTION="10m"

#IMAGE=quay.io/coreos/etcd:v3.3.9

etcd \
    --data-dir=/tmp/etcd-memory-test-n1 \
    --name=etcd-cluster-n1 \
    --initial-advertise-peer-urls=http://localhost:2380 \
    --listen-peer-urls=http://localhost:2380 \
    --listen-client-urls=http://localhost:2379 \
    --advertise-client-urls=http://localhost:2379 \
    --initial-cluster=etcd-cluster-n1=http://localhost:2380
#    --initial-cluster-state=existing
