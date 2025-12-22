#!/usr/bin/env bash

if [ "$1" = "" ]; then
  echo "Usage: ./gen-certs.sh <cert CN>"
  exit 0
fi

mkdir -p "./$1_cert"

dir="./$1_cert"

openssl genrsa -out "./$dir/$1.key" 2048
openssl req -new -key "./$dir/$1.key" -out "./$dir/$1.csr" -subj "/CN=$1"
openssl x509 -req -in "./$dir/$1.csr" -CA "./docker/nginx/ssl/nginx.crt" -CAkey "./docker/nginx/ssl/nginx.key" \
  -CAcreateserial -out "./$dir/$1.crt" -days 365

echo "Client certificates generated in ./$1_cert/"
