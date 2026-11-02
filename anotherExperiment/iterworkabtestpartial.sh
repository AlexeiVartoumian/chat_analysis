#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
    if [ "$type" = "OUTPUT" ]; then
        #./start insert $file COMPANY_WORK
        ./start insert $file JOBS_WORKAB
        ./start insert $file JOB_DESCRIPTIONS_WORKAB
        #./start insert $file JOB_LIFECYCLE_WORK
    
    elif [ "$type" = "COMPANIES" ]; then
        ./start insert $file COMPANY_WORKAB

    elif [ "$type" = "PROCESSEDJOBS" ]; then
        #./start insert $file COMPANY_WORKAB
        ./start insert $file JOBS_WORKAB_PARTIAL
        ./start insert $file JOB_LIFECYCLE_WORKAB
    
    elif [ "$type" = "DEADLINKS" ]; then
        ./start insert $file JOB_LIFECYCLE_WORKAB
    fi
done < <(jq -r '.[][][][]' keys_workab.json)
