#!/bin/sh
_=$(grep Version README.md | awk '{print "git rev-parse refs/tags/"$3" > /dev/null 2>&1;"}' | sh)
if [ "$(echo $?)" == 0 ]; then exit 1; fi;