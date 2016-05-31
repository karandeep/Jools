#!/bin/bash

case $1 in
  start)
	cd /home/ubuntu/jools/;
    echo $$ > /home/ubuntu/pids/jools.pid;
    exec ./jools -port=8081 2>> /tmp/jools_prod
    ;;
  stop)
    kill `cat /home/ubuntu/pids/jools.pid`;;
  *)
  echo "usage: ./jools.sh {start|stop}" ;;
esac
exit 0
