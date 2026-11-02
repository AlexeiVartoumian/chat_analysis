#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
    if [ "$type" = "OUTPUT" ]; then
        ./start insert $file COMPANY_BAMBI
        ./start insert $file JOBS_BAMBI_PARTIAL
        ./start insert $file JOBS_BAMBI
        ./start insert $file JOB_DESCRIPTIONS_BAMBI
        ./start insert $file JOB_LIFECYCLE_BAMBI
    
   elif [[ "$type" == PROCESSEDJOBS* ]]; then
        ./start insert $file COMPANY_BAMBI
        ./start insert $file JOBS_BAMBI_PARTIAL
        ./start insert $file JOB_LIFECYCLE_BAMBI
    
    elif [ "$type" = "DEPARTMENTS" ]; then
        ./start insert $file JOBS_BAMBI_PARTIAL
    
    elif [ "$type" = "DEADLINKS" ]; then
        ./start insert $file JOB_LIFECYCLE_BAMBI
    fi
done < <(jq -r '.[][][]' keys_bambi.json)