#!/bin/sh

for ((i=0 ; i<100000; i++))
do
  curl http://localhost:30001/v1/health
  echo ""
  sleep 0.1s
done
