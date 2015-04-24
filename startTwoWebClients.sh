#/bin/bash

go build -o dibbler

./dibbler -port=8000 -logfile=port8000.log &
echo $! > ./port8000Dibbler-$!.pid
./dibbler -port=8001 -logfile=port8001.log &
echo $! > ./port8001Dibbler-$!.pid
