#!/bin/bash

export MONGODB_URI="mongodb://localhost:27017/"
export MONGODB_DATABASE="post-database"
cd local && docker compose up -d && cd ..
go run cmd/main.go