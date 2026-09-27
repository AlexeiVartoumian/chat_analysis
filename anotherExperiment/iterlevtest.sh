#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
    if [ "$type" = "PROCESSEDJOBSLEVER" ]; then
        ./start insert $file COMPANY_LEV
        ./start insert $file JOBS_LEV
        ./start insert $file JOB_LIFECYCLE_LEV
    elif [ "$type" = "JOBDESCRIPTIONSLEVER" ]; then
        ./start insert $file JOB_DESCRIPTIONS_LEV
        
    fi
done < <(jq -r '.[][][][]' keys_lev.json)



