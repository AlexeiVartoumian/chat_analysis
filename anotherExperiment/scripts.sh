#!/bin/bash


sudo apt update
sudo apt install -y apt-transport-https ca-certificates curl gnupg lsb-release

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io

sudo systemctl start docker
sudo usermod -aG docker $USER
newgrp docker


docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=yourpassword \
  -e POSTGRES_DB=yourdb \
  -p 5432:5432 \
  pgvector/pgvector:pg16

sudo apt install python3-pip
sudo apt install python3.12-venv
sudo python3 -m venv /venv
source /venv/bin/activate
sudo /venv/bin/pip install psycopg2-binary pgvector python-dotenv numpy

#pip install psycopg2-binary pgvector python-dotenv numpy


