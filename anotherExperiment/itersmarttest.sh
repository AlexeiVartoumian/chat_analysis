#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
   
    if [ "$type" = "PROCESSEDJOBS" ]; then
        ./start insert $file COMPANY_SMART
        ./start insert $file JOBS_SMART
        ./start insert $file JOB_LIFECYCLE_SMART
    
    fi
done < <(jq -r '.[][][][]' keys_smart.json)