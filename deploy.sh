#! /bin/bash

cp -R -p ./templates ./bin/

mkdir ./bin/videos

cd bin

nohup ./api &
nohup ./scheduler &
nohup ./stream_server &
nohup ./web &

echo "deploy finished"