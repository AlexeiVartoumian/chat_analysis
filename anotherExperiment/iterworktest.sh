#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
    if [ "$type" = "OUTPUT" ]; then
        #./start insert $file COMPANY_WORK
        ./start insert $file JOBS_WORK
        ./start insert $file JOB_DESCRIPTIONS_WORK
        #./start insert $file JOB_LIFECYCLE_WORK
    
    elif [ "$type" = "PROCESSEDJOBS" ]; then
        ./start insert $file COMPANY_WORK
        ./start insert $file JOBS_WORK_PARTIAL
        ./start insert $file JOB_LIFECYCLE_WORK
    
    elif [ "$type" = "DEADLINKSWORKDAY" ]; then
        ./start insert $file JOB_LIFECYCLE_WORK
    fi
done < <(jq -r '.[][][][]' keys_work.json)
