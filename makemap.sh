#!/bin/sh

set -e

rm -rf emitter
rm -rf worker
mkdir emitter
mkdir worker

find ./output -type f -name '*.json' > worker/jsons
CNT=0
while read line
do
        CNT=`expr $CNT + 1`
        cp $line emitter/$CNT.map.json
        echo "cp $line emitter/$CNT.map.json"
done < worker/jsons

find ./output -type f -name '*.png' > worker/pngs
CNT=0
while read line
do
        CNT=`expr $CNT + 1`
        cp $line emitter/$CNT.png
        echo "cp $line emitter/$CNT.png"
done < worker/pngs
