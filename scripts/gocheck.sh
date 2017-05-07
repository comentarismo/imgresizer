#!/usr/bin/env bash

go list -f '{{join .Deps "\n"}}' |grep comentarismo
go list -f '{{join .DepsErrors "\n"}}' comentarismo/server