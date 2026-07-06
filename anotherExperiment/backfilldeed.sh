#!/bin/bash


while IFS= read -r file; do
    
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"

    if [ "$type" = "LIVE" ]; then 
        ./start insert $file JOB_LIFECYCLE_DEED
    
    elif [ "$type" = "EXPIRED" ]; then
        ./start insert $file JOB_LIFECYCLE_DEED
    
    fi 

done < <(jq -r '.[][][][]' keys_deed.json)