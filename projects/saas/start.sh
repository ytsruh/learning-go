#!/bin/bash

clear
rm ytsruh.com/saas
go build --tags="mgo"
./ytsruh.com/saas -driver mongo -datasource "127.0.0.1"