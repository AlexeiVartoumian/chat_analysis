#!/bin/bash

#t0d0 check how keys.json is being creatd with print on indexes . 
#cehck listbucketv2 when ready . if keys are there then this will work . if not debug on  set of indexes to see why 
# keys.json does not containt he live and suspended roles / 


## update . looks to be good ! listbucketv2 works because of "w" flag where we create a fresh batch each time .
## therefore all we have to do here is follow the business rules of the db piping keys.json to this file . 
while IFS= read -r file; do

    echo "$file"

    type="${file%%-*}"
    type="${type^^}"


    if [ "$type" = "LIVE" ]; then 
        ./start insert $file JOB_LIFECYCLE_UPDATE
    elif [ "$type" = "EXPIRED" ]; then
        ./start insert $file JOB_LIFECYCLE_UPDATE
    
    elif [ "$type" = "SUSPENDED" ]; then
        ./start insert $file JOB_LIFECYCLE_UPDATE_SUSPENDED
    
    fi
done < <(jq -r '.[][][][]' keys.json)