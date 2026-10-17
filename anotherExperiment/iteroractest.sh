#!/bin/bash
while IFS= read -r file; do
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"
    #echo $type
    if [ "$type" = "COMPANYORAC" ]; then
        ./start insert $file COMPANY_ORAC
   
    elif [ "$type" = "PROCESSEDJOBSORAC" ]; then
        ./start insert $file JOBS_ORAC
        ./start insert $file JOB_LIFECYCLE_ORAC
    
    fi
done < <(jq -r '.[][][]' keys_orac.json)