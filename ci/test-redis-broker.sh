#!/bin/bash

export GOPATH=$GOPATH:$PWD
cd src/github.com/bcshuai/cf-redis-broker
./script/test
