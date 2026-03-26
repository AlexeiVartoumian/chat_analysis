from __future__ import print_function



TABLES = {}

# As such the business rules are:

# 1. A company can post multiple job advertisements
# 2. A Candidate can apply for multiple job advertisements

# 3. A job has attributes: , jobid,  title , company , location , salary . posted date, job url , status.
# 4. A job may be split into addtional tables or separate tables as features could be added such as profession relation.

# 5. A company has attributes : id,  title  logo , industry/about?

# 6. a peson has attributes : id , email , username , role , password

# 8. an application has attribute : application id job id user id amd status


"job_id", "industry",          "company_name",     "company_about", "Employee count"
"job_id", "Applicants count", "Company Apply URL", "job_state"

TABLES['COMPANY'] = (
    "CREATE TABLE IF NOT EXISTS `COMPANY` ("
    " `company_id` BIGINT(11) NOT NULL," # extract company id
    " `name` varchar(64) NOT NULL,"
    " `logo` varchar(512),"
    " PRIMARY KEY (`company_id`)"
    ") ENGINE=InnoDB"
)

TABLES['COMPANY_METADATA'] = (
    "CREATE TABLE IF NOT EXISTS `COMPANY_METADATA` ("
    " `company_id` BIGINT(11) NOT NULL ," # extract company id
    " `name` varchar(64) NOT NULL,"
    " `industry` varchar(64) NOT NULL,"
    " `description` TEXT," 
    " `employee_count` int(11),"
    " `employee_count_range` varchar(64) NOT NULL," 
     "  CONSTRAINT `COMPANY_METADATA_ibfk_1` FOREIGN KEY (`company_id`)"
    "   REFERENCES `COMPANY` (`company_id`) ON DELETE CASCADE,"
    " PRIMARY KEY (`company_id`)"
    ") ENGINE=InnoDB"
)


TABLES['JOBS'] = (
    "CREATE TABLE IF NOT EXISTS `JOBS` ("
    " `job_id` BIGINT(11) NOT NULL," #extract job id as id
    " `title` varchar(256) NOT NULL,"
    " `location` varchar(64) NOT NULL,"
    " `salary` varchar(64),"
    " `date_posted` date,"
    " `job_url` varchar(256) NOT NULL,"
    " `easy_apply` varchar(256),"
    " `promoted` varchar(256),"
    " `expiry_date` date,"
    " `company_id` BIGINT(11) NOT NULL,"
    " CONSTRAINT `JOBS_ibfk_1` FOREIGN KEY (`company_id`)"
    "   REFERENCES `COMPANY` (`company_id`) ON DELETE CASCADE,"
    " PRIMARY KEY (`job_id`)"
    ") ENGINE=InnoDB"
)

TABLES['JOB_METADATA'] = (
    "CREATE TABLE IF NOT EXISTS `JOBS_METADATA` ("
    " `job_id` BIGINT(11) NOT NULL AUTO_INCREMENT," # extract company id
    " `applicants_count` varchar(512) NOT NULL,"
    " `company_apply_url` varchar(512) NOT NULL,"
    " `job_state` varchar(256) NOT NULL,"
    "  CONSTRAINT `JOBS_METADATA_ibfk_1` FOREIGN KEY (`job_id`)"
    "   REFERENCES `JOBS` (`job_id`) ON DELETE CASCADE,"
    " PRIMARY KEY (`job_id`)"
    ") ENGINE=InnoDB"
)

TABLES['USER'] = (
    "CREATE TABLE IF NOT EXISTS `USER` ("
    " `user_id` int(11) NOT NULL AUTO_INCREMENT,"
    " `first_name` varchar(64) NOT NULL,"  
    " `last_name` varchar(64) NOT NULL,"   
    " `email` varchar(128) UNIQUE,"
    " `password_hash` VARCHAR(255),"
    " `age` int(11) NOT NULL,"
    " PRIMARY KEY (`user_id`)"
    ") ENGINE=InnoDB"
)


TABLES['APPLICATION'] = (
    "CREATE TABLE IF NOT EXISTS `APPLICATION` ("
    " `application_id` int(11) NOT NULL AUTO_INCREMENT,"
    " `status` enum('live', 'expired'),"
    " `user_id` int(11) NOT NULL,"
    " `job_id` BIGINT(11) NOT NULL,"
    " `age` int(11),"
    " CONSTRAINT `APPLICATION_ibfk_1` FOREIGN KEY (`user_id`)"
    "   REFERENCES `USER` (`user_id`) ON DELETE CASCADE,"
    " CONSTRAINT `APPLICATION_ibfk_2` FOREIGN KEY (`job_id`)"
    "   REFERENCES `JOBS` (`job_id`) ON DELETE CASCADE,"
    " PRIMARY KEY (`application_id`)"
    ") ENGINE=InnoDB"
)

TABLES['JOB_DESCRIPTION'] = (
    "CREATE TABLE IF NOT EXISTS `JOBS_DESCRIPTION` ("
    " `job_id` BIGINT(11) NOT NULL AUTO_INCREMENT," # extract company id
    " `job_state` varchar(24000) NOT NULL,"
    "  CONSTRAINT `JOBS_DESCRIPTION_ibfk_1` FOREIGN KEY (`job_id`)"
    "   REFERENCES `JOBS` (`job_id`) ON DELETE CASCADE,"
    " PRIMARY KEY (`job_id`)"
    ") ENGINE=InnoDB"
)