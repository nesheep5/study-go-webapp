#!/bin/sh
go build
echo "build OK."

./setenv.sh
echo "set env OK."

# app start
./chat
