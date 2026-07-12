#!/bin/bash


while IFS= read -r file; do
    
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"

    if [ "$type" = "PROCESSEDJOBSASH" ]; then 
        ./start insert $file COMPANY_ASH
        ./start insert $file JOBS_ASH
        ./start insert $file JOB_LIFECYCLE_ASH
    
    elif [ "$type" = "JOBDESCRIPTIONASH" ]; then
        ./start insert $file JOB_DESCRIPTIONS_ASH
    
    fi 

done < <(jq -r '.[][][][]' keys_ash.json)