#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
    if [ "$type" = "OUTPUT" ]; then
        ./start insert $file COMPANY_WORK
        ./start insert $file JOBS_WORK
        ./start insert $file JOB_DESCRIPTION
        ./start insert $file JOB_LIFECYCLE_WORK
        
    fi
done < <(jq -r '.[][][][]' keys_work.json)