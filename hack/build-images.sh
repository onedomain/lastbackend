#!/bin/bash

declare -a arr=("lastbackend" "ingress" "discovery" "exporter")

if [[ $1 != "" ]]; then
  arr=($1)
fi

## now loop through the components array
for i in "${arr[@]}"
do
 echo "Build '$i' version '$VERSION'"
 #docker build -t "index.0xqi.com/lastbackend/$i" -f "./images/$i/Dockerfile" .
 docker build -t "kulado.azurecr.io/lastbackend/$i" -f "./images/$i/Dockerfile" .
done
