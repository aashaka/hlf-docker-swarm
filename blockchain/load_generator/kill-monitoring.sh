#!/bin/bash

# Aashaka: Note that this measures resource usage in the host, not the containers.

for i in peer{1..4} ;
do
    ssh $i 'pkill nmon; pkill tcpdump; pkill inotifywait; rm -rf *.nmon *.pcap *inotify.log data.zip' &
    rm -rf logs/$i
    mkdir -p logs/$i
done

wait
