#!/bin/bash

while `true`; do
    ps -o command="",cpu="",rss="" -p $(pgrep etcd)
    sleep 10
done \
    | ts '%s' \
    | tee local-memory.log
