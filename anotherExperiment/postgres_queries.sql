-- reposted over the weekend
SELECT
    DATE(j.date_posted) as posted_date,
    TO_CHAR(j.date_posted, 'Day') as day_of_week,
    j.title,
    c.name as company_name,
    j.location,
    j.salary,
    jm.company_apply_url,
    jm.applicants_count,
    j.job_url
FROM jobs j
JOIN company c ON j.company_id = c.company_id
JOIN job_metadata jm ON j.job_id = jm.job_id
WHERE j.date_posted IS NOT NULL
  AND (EXTRACT(DOW FROM j.date_posted) IN (0, 6)  -- 0 = Sunday, 6 = Saturday
       OR jm.applicants_count LIKE '%Reposted%')
ORDER BY j.date_posted DESC, c.name ASC;


-- posted in last 3 days
SELECT
    j.date_posted,
    TRIM(TO_CHAR(j.date_posted, 'Day')) as day_posted,
    j.title,
    c.name as company_name,
    j.location,
    COALESCE(j.salary, 'Not Specified') as salary,
    jm.applicants_count,
    jm.company_apply_url,
    j.job_url
FROM jobs j
JOIN company c ON j.company_id = c.company_id
JOIN job_metadata jm ON j.job_id = jm.job_id
WHERE j.date_posted >= CURRENT_DATE - INTERVAL '3 days'
  AND jm.company_apply_url NOT LIKE 'https://www.linkedin%'
  AND jm.applicants_count NOT LIKE '%Reposted%'
ORDER BY j.date_posted DESC;


-- Most recent jobs categorized by date posted with company apply URL
-- Excluding jobs where applicants_count contains "Reposted"
SELECT 
    DATE(j.date_posted) as posted_date,
    j.title,
    c.name as company_name,
    j.location,
    j.salary,
    jm.company_apply_url,
    jm.applicants_count,
    j.job_url
FROM JOBS j
JOIN COMPANY c ON j.company_id = c.company_id
JOIN JOBS_METADATA jm ON j.job_id = jm.job_id
WHERE j.date_posted IS NOT NULL
  AND (jm.applicants_count NOT LIKE '%Reposted%' 
       OR jm.applicants_count IS NULL)
ORDER BY j.date_posted DESC, c.name ASC;

-- apply url not in linkedin
SELECT 
    j.date_posted,
    j.title,
    c.name as company_name,
    j.location,
    j.salary,
    jm.applicants_count,
    jm.company_apply_url,
    j.job_url
FROM JOBS j
JOIN COMPANY c ON j.company_id = c.company_id
LEFT JOIN JOB_METADATA jm ON j.job_id = jm.job_id
WHERE j.date_posted IS NOT NULL
  AND j.salary IS NOT NULL
  AND j.salary != 'Not specified'
  AND (jm.company_apply_url NOT LIKE 'https://www.linkedin%'
       OR jm.company_apply_url IS NULL)
ORDER BY j.date_posted DESC;

-- jobs with less then 100 people employed 

SELECT 
    j.date_posted,
    j.title,
    c.name as company_name,
    j.location,
    COALESCE(j.salary, 'Not Specified') as salary,
    cm.employee_count,
    cm.employee_count_range,
    cm.industry,
    jm.applicants_count,
    jm.company_apply_url,
    j.job_url
FROM JOBS j
JOIN COMPANY c ON j.company_id = c.company_id
LEFT JOIN COMPANY_METADATA cm ON c.company_id = cm.company_id
LEFT JOIN JOB_METADATA jm ON j.job_id = jm.job_id
WHERE j.date_posted IS NOT NULL
  AND (cm.employee_count < 100 
       OR cm.employee_count_range IN ('1-10', '11-50', '51-200'))
ORDER BY j.date_posted DESC
LIMIT 50;

-- Companies that repost jobs, sorted by employee count with applicant data
SELECT 
    c.company_id,
    c.name as company_name,
    cm.industry,
    cm.employee_count,
    cm.employee_count_range,
    COUNT(j.job_id) as total_jobs,
    SUM(CASE WHEN jm.applicants_count LIKE '%Reposted%' THEN 1 ELSE 0 END) as reposted_jobs,
    ROUND(SUM(CASE WHEN jm.applicants_count LIKE '%Reposted%' THEN 1 ELSE 0 END) * 100.0 / COUNT(j.job_id), 2) as repost_percentage,
    STRING_AGG(DISTINCT jm.applicants_count, ' | ') as all_applicant_counts
FROM company c
JOIN jobs j ON c.company_id = j.company_id
JOIN job_metadata jm ON j.job_id = jm.job_id
LEFT JOIN company_metadata cm ON c.company_id = cm.company_id
GROUP BY c.company_id, c.name, cm.industry, cm.employee_count, cm.employee_count_range
HAVING SUM(CASE WHEN jm.applicants_count LIKE '%Reposted%' THEN 1 ELSE 0 END) > 0
ORDER BY 
    cm.employee_count ASC,
    CASE cm.employee_count_range
        WHEN '1-10' THEN 1
        WHEN '11-50' THEN 2
        WHEN '51-200' THEN 3
        WHEN '201-500' THEN 4
        WHEN '501-1000' THEN 5
        WHEN '1001-5000' THEN 6
        WHEN '5001-10000' THEN 7
        WHEN '10001+' THEN 8
        ELSE 9
    END,
    reposted_jobs DESC;


-- companies lest then 100 people apply 
SELECT 
    j.date_posted,
    j.title,
    c.name as company_name,
    j.location,
    COALESCE(j.salary, 'Not Specified') as salary,
    cm.employee_count_range,
    cm.industry,
    jm.applicants_count,
    jm.company_apply_url,
    j.job_url
FROM jobs j
JOIN company c ON j.company_id = c.company_id
LEFT JOIN company_metadata cm ON c.company_id = cm.company_id
JOIN job_metadata jm ON j.job_id = jm.job_id
WHERE j.date_posted >= CURRENT_DATE - INTERVAL '14 days'
  AND jm.applicants_count NOT LIKE '%Reposted%'
  AND jm.applicants_count NOT LIKE '%Over 100%'
  AND (cm.employee_count < 500 OR cm.employee_count_range IN ('1-10', '11-50', '51-200', '201-500'))
  AND (REGEXP_MATCH(jm.applicants_count, '\s(\d+)\s*(applicants|people clicked apply)'))[1]::INTEGER < 50
ORDER BY 
    (REGEXP_MATCH(jm.applicants_count, '\s(\d+)\s*(applicants|people clicked apply)'))[1]::INTEGER ASC NULLS LAST;


-- Companies with less than 51 employees and their company apply URLs
SELECT 
    c.company_id,
    c.name as company_name,
    cm.industry,
    cm.employee_count,
    cm.employee_count_range,
    cm.description,
    COUNT(j.job_id) as total_jobs,
    STRING_AGG(DISTINCT j.title, ' | ') as job_titles,
    STRING_AGG(DISTINCT j.location, ' | ') as locations,
    STRING_AGG(DISTINCT jm.company_apply_url, ' | ') as company_apply_urls
FROM company c
LEFT JOIN company_metadata cm ON c.company_id = cm.company_id
JOIN jobs j ON c.company_id = j.company_id
JOIN job_metadata jm ON j.job_id = jm.job_id
WHERE (cm.employee_count < 51 
       OR cm.employee_count_range IN ('1-10', '11-50'))
  AND jm.company_apply_url IS NOT NULL
GROUP BY c.company_id, c.name, cm.industry, cm.employee_count, cm.employee_count_range, cm.description
ORDER BY cm.employee_count ASC, total_jobs DESC;


SELECT 
    c.company_id,
    c.name as company_name,
    cm.industry,
    cm.employee_count,
    cm.employee_count_range,
    cm.description,
    COUNT(j.job_id) as total_jobs,
    STRING_AGG(DISTINCT j.title, ' | ') as job_titles,
    STRING_AGG(DISTINCT j.location, ' | ') as locations,
    STRING_AGG(DISTINCT jm.company_apply_url, ' | ') as company_apply_urls
FROM company c
LEFT JOIN company_metadata cm ON c.company_id = cm.company_id
JOIN jobs j ON c.company_id = j.company_id
JOIN job_metadata jm ON j.job_id = jm.job_id
WHERE (cm.employee_count >= 500 
       OR cm.employee_count_range IN ('501-1000', '1001-5000', '5001-10000', '10001+'))
  AND jm.company_apply_url IS NOT NULL
GROUP BY c.company_id, c.name, cm.industry, cm.employee_count, cm.employee_count_range, cm.description
ORDER BY total_jobs DESC, cm.employee_count DESC;

-- search term performance
SELECT 
    st.term,
    sw.run_at,
    sw.total_jobs_found,
    sw.net_new_jobs,
    (sw.total_jobs_found - sw.net_new_jobs) as duplicate_jobs,
    ROUND(sw.net_new_jobs * 100.0 / NULLIF(sw.total_jobs_found, 0), 2) as fresh_pct
FROM SEARCH_WORKFLOW sw
JOIN SEARCH_TERM st ON sw.search_term_id = st.search_term_id
ORDER BY sw.run_at DESC;

--dupes on j's
SELECT 
    j.job_id,
    j.title,
    c.name as company_name,
    j.date_posted,
    COUNT(DISTINCT jst.workflow_id) as times_surfaced,
    STRING_AGG(DISTINCT st.term, ' | ') as matched_search_terms
FROM JOB_SEARCH_TERM jst
JOIN JOBS j ON jst.job_id = j.job_id
JOIN COMPANY c ON j.company_id = c.company_id
JOIN SEARCH_WORKFLOW sw ON jst.workflow_id = sw.workflow_id
JOIN SEARCH_TERM st ON sw.search_term_id = st.search_term_id
GROUP BY j.job_id, j.title, c.name, j.date_posted
HAVING COUNT(DISTINCT jst.workflow_id) > 1
ORDER BY times_surfaced DESC;

--runs
SELECT DISTINCT ON (st.term)
    st.term,
    sw.run_at,
    sw.total_jobs_found,
    sw.net_new_jobs,
    (sw.total_jobs_found - sw.net_new_jobs) as duplicates
FROM SEARCH_TERM st
JOIN SEARCH_WORKFLOW sw ON st.search_term_id = sw.search_term_id
ORDER BY st.term, sw.run_at DESC;

SELECT 
    st.term,
    COUNT(sw.workflow_id) as total_runs,
    SUM(sw.total_jobs_found) as total_found,
    SUM(sw.net_new_jobs) as total_fresh,
    ROUND(AVG(sw.net_new_jobs * 100.0 / NULLIF(sw.total_jobs_found, 0)), 2) as avg_fresh_pct,
    MAX(sw.run_at) as last_run
FROM SEARCH_TERM st
JOIN SEARCH_WORKFLOW sw ON st.search_term_id = sw.search_term_id
GROUP BY st.term
ORDER BY avg_fresh_pct ASC;


SELECT c.company_id , c.name,count(j.job_id) as total_jobs from COMPANY c JOIN JOBS j on c.company_id = j.company_id group by c.company_id ORDER BY total_jobs DESC;

SELECT Count(name) , Industry FROM COMPANY_METADATA GROUP BY Industry ORDER BY Count(Name) DESC;


-- jobs closing in three days
SELECT jl.job_id  ,jl.job_state ,jl.first_seen_at  ,jl.last_seen_listed_at  ,jl.first_seen_closed_at  , jl.next_scan_at ,  
jl.suspended_count ,JOBS.date_posted FROM JOB_LIFECYCLE jl , JOBS 
WHERE (jl.first_seen_closed_at::date - jl.first_seen_at::date) < 4 
AND jl.first_seen_closed_at != '2026-04-24 18:46:03+00' --todo remove bad dates
AND JOBS.job_id = jl.job_id; 


SELECT jl.job_id  ,jl.job_state ,jl.first_seen_at  ,jl.last_seen_listed_at  ,jl.first_seen_closed_at  , jl.next_scan_at , jl.suspended_count , JOBS.job_url ,JOBS.date_posted ,c.name , c.company_id FROM JOB_LIFECYCLE jl , JOBS , COMPANY_METADATA c WHERE (jl.first_seen_closed_at::date - jl.first_seen_at::date) < 4 AND jl.first_seen_closed_at != '2026-04-24 18:46:03+00'AND JOBS.job_id = jl.job_id and JOBS.company_id = c.company_id ORDER BY JOBS.date_posted DESC;