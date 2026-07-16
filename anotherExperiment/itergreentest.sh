#!/bin/bash


while IFS= read -r file; do
    
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    
    if [ "$type" = "COMPANYGREEN"]; then
        ./start insert $file COMPANY_GREEN
    if [ "$type" = "PROCESSEDJOBSGREEN" ]; then 
        ./start insert $file JOBS_GREEN
        ./start insert $file JOB_LIFECYCLE_GREEN
    
    elif [ "$type" = "JOBDESCRIPTIONGREEN" ]; then
        ./start insert $file JOB_DESCRIPTIONS_GREEN
    
    fi 

done < <(jq -r '.[][][][]' keys_green.json)