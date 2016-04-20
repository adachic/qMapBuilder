#!/bin/sh

set -e

MAPDIR=/Users/adachic/qenemy/map
OUTDIR=/Users/adachic/qenemy/output
ENEMYDIR=/Users/adachic/qenemy

rm -rf $OUTDIR
rm -rf $MAPDIR
rm -rf worker
mkdir $OUTDIR
mkdir $MAPDIR
mkdir worker

find ./output -type f -name '*.json' > worker/jsons
CNT=0
while read line
do
        CNT=`expr $CNT + 1`
        echo "cp $line $MAPDIR/$CNT.map.json"
        cp $line $MAPDIR/$CNT.map.json
        cp $line $OUTDIR/$CNT.map.json
done < worker/jsons

find ./output -type f -name '*.png' > worker/pngs
CNT=0
while read line
do
        CNT=`expr $CNT + 1`
        echo "cp $line $OUTDIR/$CNT.png"
        cp $line $OUTDIR/$CNT.png
done < worker/pngs

echo "cd $ENEMYDIR"
cd $ENEMYDIR

echo "./run10.sh 1 $CNT"
./run10.sh 1 $CNT
