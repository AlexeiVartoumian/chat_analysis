#!/bin/bash

#t0d0 check how keys.json is being creatd with print on indexes . 
#cehck listbucketv2 when ready . if keys are there then this will work . if not debug on  set of indexes to see why 
# keys.json does not containt he live and suspended roles / 
WHILE IFS= read -r file; do

    echo "$file"

    type="${file%%-*}"
    type="${type^^}"

    if [ "$type" =  ]