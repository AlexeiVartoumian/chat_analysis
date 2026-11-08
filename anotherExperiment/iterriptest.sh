#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
    if [ "$type" = "COMPANYRIP" ]; then
        ./start insert $file COMPANY_RIP

    elif [ "$type" = "DEPARTMENTRIP" ]; then
        ./start insert $file DEPARTMENT_RIP

    elif [ "$type" = "PROCESSEDJOBSRIP" ]; then
        ./start insert $file JOBS_RIP
        ./start insert $file JOB_LIFECYCLE_RIP

    elif [ "$type" = "JOBDESCRIPTIONRIP" ]; then
        ./start insert $file JOB_DESCRIPTIONS_RIP

    elif [ "$type" = "LIVERIP" ]; then
        ./start insert $file LIVE_RIP

    elif [ "$type" = "DEADLINKS" ]; then
        ./start insert $file DEADLINKS_RIP

    fi
done < <(jq -r '.[][][][]' keys_rip.json)