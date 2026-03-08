#!/bin/bash


psql -U postgres -c "CREATE DATABASE chat_analysis;"

psql -U postgres -d chat_analysis -f schema.sql
cd "C:\Program Files\PostgreSQL\18\bin"
psql -U postgres -c "CREATE EXTENSION IF NOT EXISTS vector;" chat_analysis