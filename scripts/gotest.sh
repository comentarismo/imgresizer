#!/usr/bin/env bash

rediscachedisabled=true godep go test -v $(go list ./limiter | grep -v /vendor/);