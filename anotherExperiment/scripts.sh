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
sudo apt install python3-pip -y
sudo apt install python3.12-venv -y
sudo python3 -m venv /venv

source /venv/bin/activate
sudo /venv/bin/pip install psycopg2-binary pgvector python-dotenv numpy boto3

# docker run -d \
#   --name postgres \
#   -e POSTGRES_PASSWORD=yourpassword \
#   -e POSTGRES_DB=yourdb \
#   -p 5432:5432 \
#   pgvector/pgvector:pg16
git clone --single-branch --branch model-insert https://github.com/AlexeiVartoumian/chat_analysis.git
sudo mv chat_analysis/anotherExperiment/ssmparam.py /home/ubuntu/ssmparam.py
python3 ssmparam.py
docker run -d \
  --name postgres \
  --env-file .env2 \
  -p 5432:5432 \
  pgvector/pgvector:pg16


sudo mv chat_analysis/anotherExperiment/model.py /home/ubuntu/model.py
sudo mv chat_analysis/anotherExperiment/update_search.py /home/ubuntu/update_search.py
sudo mv chat_analysis/anotherExperiment/listbucket.py /home/ubuntu/listbucket.py

sudo mv chat_analysis/anotherExperiment/hosted.py /home/ubuntu/hosted.py
my_ip=$(curl http://checkip.amazonaws.com)

python3 hosted.py $my_ip
python3 model.py

sudo apt install golang-go -y
sudo rm -rf chat_analysis/
git clone --single-branch --branch  cli-binary-code https://github.com/AlexeiVartoumian/chat_analysis.git
cd /home/ubuntu/chat_analysis/anotherExperiment/cli/
sudo GOOS=linux GOARCH=amd64 go build -o start main.go
sudo mv start /home/ubuntu/start && cd /home/ubuntu/
sudo chmod +x start

python3 listbucket.py #todo this should be triggered by something else
./start insert processedJobs.csv COMPANY
./start insert processedJobs.csv JOBS
./start insert company_data.csv COMPANY_METADATA
./start insert job_metadata.csv JOB_METADATA
./start insert job_description.csv JOB_DESCRIPTION
sudo rm -rf chat_analysis/
git clone --single-branch --branch  api-code-tls https://github.com/AlexeiVartoumian/chat_analysis.git
cd chat_analysis/anotherExperiment/api/cmd/
sudo go run main.go