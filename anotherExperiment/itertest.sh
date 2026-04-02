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
    else
        ./start insert $file $type

    fi
done < <(jq -r '.[][]' keys.json)