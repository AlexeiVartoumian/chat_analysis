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