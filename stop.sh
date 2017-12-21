#!/bin/bash


PROC=`pwd | xargs -i basename {}`

echo "stoping ${PROC}...!"

PIDS=`ps -fu${USER} |grep ${PROC} | grep -v grep |grep cfg| awk '{print $2}'`

KILL=0
for PID in $PIDS
do
    kill  $PID
    KILL=1
done
if [ $KILL = 1 ]
then
    echo "stop ${PROC} ok!"
else
    echo "Nothing stop!"
fi
