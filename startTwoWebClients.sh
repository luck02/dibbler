#/bin/bash

go build -o dibbler

./dibbler -port=8000 -logfile=port8000.txt &

./dibbler -port=8001 -logfile=port8001.txt &
