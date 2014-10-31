#!/bin/bash

#printf "** Building linux/386\n"
#go-linux-386 build -a -o bin/linux-386/masq github.com/zerklabs/masq

printf "** Building linux/amd64\n"
go-linux-amd64 build -a -o bin/linux-amd64/masq github.com/zerklabs/masq

#printf "** Building windows/386\n"
#go-windows-386 build -o bin/windows-386/masq.exe github.com/zerklabs/masq

#printf "** Building windows/amd64\n"
#go-windows-amd64 build -o bin/windows-amd64/masq.exe github.com/zerklabs/masq
