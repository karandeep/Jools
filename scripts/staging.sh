#!/bin/bash

case $1 in
  start)
	cd /home/ubuntu/jools/staging/;
    echo $$ > /home/ubuntu/pids/staging.pid;
    exec ./staging -port=8082 2>> /tmp/jools_staging
    ;;
  stop)
    kill `cat /home/ubuntu/pids/staging.pid`;;
  *)
  echo "usage: ./staging.sh {start|stop}" ;;
esac
exit 0

