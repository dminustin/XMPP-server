#!/bin/sh
if ps -ax|grep xmpp-server.go|grep -v grep ;
then echo "Already running" ;
else go run ./xmpp-server.go ;
fi