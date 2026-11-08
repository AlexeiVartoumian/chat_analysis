#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
    if [ "$type" = "COMPANYTAILOR" ]; then
        ./start insert $file COMPANY_TAILOR

    elif [ "$type" = "PROCESSEDJOBSTAILOR" ]; then
        ./start insert $file JOBS_TAILOR
        ./start insert $file JOB_LIFECYCLE_TAILOR

    elif [ "$type" = "JOBDESCRIPTIONTAILOR" ]; then
        ./start insert $file JOB_DESCRIPTIONS_TAILOR

    elif [ "$type" = "LIVETAILOR" ]; then
        ./start insert $file LIVE_TAILOR

    elif [ "$type" = "DEADLINKS" ]; then
        ./start insert $file DEADLINKS_TAILOR

    fi
done < <(jq -r '.[][][][]' keys_tailor.json)