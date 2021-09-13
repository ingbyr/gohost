#!/bin/bash

#CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o gohost-windows-32bit.exe .
#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gohost-windows-64bit.exe .

CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o gohost-linux-32bit .
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gohost-linux-64bit .

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o gohost-darwin .