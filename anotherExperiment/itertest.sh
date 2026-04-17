#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
    if [ "$type" = "PROCESSEDJOBS" ]; then
        ./start insert $file COMPANY
        ./start insert $file JOBS

    elif [ "$type" = "COMPANY_DATA" ]; then
        ./start insert $file COMPANY_METADATA
    elif [ "$type" = "JOB_METADATA" ]; then
        ./start insert $file JOB_METADATA
        ./start insert $file JOB_LIFECYCLE
    else
        ./start insert $file $type

    fi
done < <(jq -r '.[][]' keys.json)