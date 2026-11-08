#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
    if [ "$type" = "COMPANYBREZ" ]; then
        ./start insert $file COMPANY_BREZ

    elif [ "$type" = "PROCESSEDJOBSBREZ" ]; then
        ./start insert $file JOBS_BREZ
        ./start insert $file JOB_LIFECYCLE_BREZ

    elif [ "$type" = "JOBDESCRIPTIONBREZ" ]; then
        ./start insert $file JOB_DESCRIPTIONS_BREZ

    elif [ "$type" = "LIVEBREZ" ]; then
        ./start insert $file LIVE_BREZ

    elif [ "$type" = "DEADLINKS" ]; then
        ./start insert $file DEADLINKS_BREZ

    fi
done < <(jq -r '.[][][][]' keys_brez.json)