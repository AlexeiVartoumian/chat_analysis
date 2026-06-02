#!/bin/bash


while IFS= read -r file; do
    
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"

    if [ "$type" = "PROCESSEDJOBSIND" ]; then 
        ./start insert $file COMPANY_DEED
        ./start insert $file JOBS_DEED
    
    elif [ "$type" = "JOBDESCRIPTIONIND" ]; then
        ./start insert $file JOB_DESCRIPTION_DEED
    fi 

done < <(jq -r '.[][][]' keys_deed.json)