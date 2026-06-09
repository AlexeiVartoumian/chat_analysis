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



#memory management
#check var usage on logs vacuuim
sudo du -sh /var/* | sort -rh | head -10
sudo journalctl --vacuum-size=50M
sudo apt clean  # cached installer files on every update .

dpkg-query -Wf '${Installed-Size}\t${Package}\n' | sort -rn | head -20 # check installed packages by size . i.e old isntalled kernels /lib/modules lying around

# sudo apt remove --purge linux-modules-6.17.0-1012-aws linux-modules-6.17.0-1013-aws # IF REMOVING OLD KERNELS BE SURE TO CHECK LATEST ONE BEING USEDS with 
# uname -r
df -h /


# chop and dice check the memory slice
aws s3api list-objects-v2 --bucket "GLORIOUS BUCKET" \
  --query 'sort_by(Contents, &Size)[:10].{Key: Key, Size: Size}' \
  --output table