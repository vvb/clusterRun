#!/bin/bash

for i in $(seq 1 4);
do
	for k in $(seq 1 250);
	do
		echo "cluster-node$i:docker run -itd --name=cont"$i"_"$k" --net=net1 alpine /bin/sh"
	done
done;
