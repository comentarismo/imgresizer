#!/usr/bin/env bash

default: gofmt
	make start
	make log


gofmtvalidate:
	scripts/gofmt_validate.sh;
	scripts/gotest.sh;

gofmt:
	scripts/gofmt_perform.sh;

start-ci:
	scripts/godep-ci.sh
	scripts/start.sh

stop-ci:
	scripts/stop.sh

start: stop
	echo "will compile app";
	godep go build -o imgresizer main.go;
	echo "will start app with dev credentials";
	PORT=3666 rediscachedisabled=true \
	nohup ./imgresizer &
	make log

stop:
	echo "will stop dev app"
	pkill imgresizer | true

killall:
	lsof -i tcp:3000 | awk 'NR!=1 {print $2}' | xargs kill | true;
	echo "will sleep 5 secs ";

status:
	ps -ef |grep imgresizer

log:
	tail -f ./nohup.out

permission:
	chmod +x scripts/godep-ci.sh;
	chmod +x scripts/gofmt_perform.sh;
	chmod +x scripts/gofmt_validate.sh;
	chmod +x scripts/gotest.sh;
	chmod +x scripts/start.sh;
	chmod +x scripts/stop.sh;
	chmod +x scripts/gorethinkdb.sh;
	chmod +x scripts/gocheck.sh;

test:
	scripts/gotest.sh;

check:
	scripts/gocheck.sh;

rethinkdb:
	scripts/gorethinkdb.sh;

.PHONY: all test rethinkdb check permission
