#!/bin/bash

vncserver -localhost no
ssh -L 5901:localhost:5901 -R 8080:localhost:1080 -D 1080 -i "thing.pem" ubuntu@ipaddr

DISPLAY=:1 xfce4-terminal & disown
xhost +local:
chromium-browser --no-sandbox --remote-debugging-port=9222 &

cd cookieextractor/
rm -rf cookies-grouped.json
npx tsx getCookies.ts

aws s3 cp cookies-grouped.json s3://destina.json


#dumps
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

sudo apt install -y postgresql-common
sudo apt install -y postgresql-client
pg_dump --version
pg_dump -h localhost -U postgres -d interview   | gzip   | aws s3 cp - s3://********/backups/dump_$(date +%Y%m%d_%H%M%S).sql.gz

#then on vm
#copy from s3 bucket to dest
gunzip /tmp/dump.sql.gz
#docker exec -i "containter-id" psql -U postgres -d interview < /tmp/dump.sql
