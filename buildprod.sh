#! /bin/bash

# Build web and other services
cd D:/GoProject/src/stream_video_server/api
env GOOS=linux GOARCH=amd64 go build -o ../bin/api

cd D:/GoProject/src/stream_video_server/scheduler
env GOOS=linux GOARCH=amd64 go build -o ../bin/scheduler

cd D:/GoProject/src/stream_video_server/stream_server
env GOOS=linux GOARCH=amd64 go build -o ../bin/stream_server

cd D:/GoProject/src/stream_video_server/web
env GOOS=linux GOARCH=amd64 go build -o ../bin/web