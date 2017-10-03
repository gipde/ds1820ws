#!/bin/sh


function compilecopy {
    cd $2
    env $1 go build -ldflags="-s -w"
    scp $2 $3
    cd ..
}

compilecopy "GOOS=linux GOARCH=arm GOARM=6" reader      pi@ds1820ws:/home/pi
#compilecopy "GOOS=linux GOARCH=amd64"       webservice  root@76b83848-66ad-479f-becf-603934bcdfaa.pub.cloud.scaleway.com:/root/ 
