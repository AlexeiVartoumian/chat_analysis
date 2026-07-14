#!/bin/bash


while IFS= read -r file; do
    
    echo "$file"
    type="${file%%-*}"
    type="${type^^}"

    
    if [ "$type" = "ASHJOBSBYCOMPANY" ]; then 
        ./start insert $file DEPARTMENT_ASH
        ./start insert $file TEAM_ASH
        ./start insert $file JOBS_ASH
        ./start insert $file JOB_LIFECYCLE_ASH
    
   
    
    fi 

done < <(jq -r '.[][][][]' keys_ash.json)