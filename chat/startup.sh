#!/bin/sh
go build
echo "build OK."

source ./setenv.sh
echo "set env OK."

# app start
./chat
