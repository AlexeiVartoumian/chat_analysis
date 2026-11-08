#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
    if [ "$type" = "COMPANYRECRUIT" ]; then
        ./start insert $file COMPANY_RECRUIT

    elif [ "$type" = "PROCESSEDJOBSRECRUIT" ]; then
        ./start insert $file JOBS_RECRUIT
        ./start insert $file JOB_LIFECYCLE_RECRUIT

    elif [ "$type" = "JOBDESCRIPTIONRECRUIT" ]; then
        ./start insert $file JOB_DESCRIPTIONS_RECRUIT

    elif [ "$type" = "LIVERECRUIT" ]; then
        ./start insert $file LIVE_RECRUIT

    elif [ "$type" = "DEADLINKS" ]; then
        ./start insert $file DEADLINKS_RECRUIT

    fi
done < <(jq -r '.[][][][]' keys_recruit.json)