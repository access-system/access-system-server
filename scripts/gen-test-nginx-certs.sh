#!/usr/bin/env bash

openssl genrsa -out ./docker/nginx/ssl/nginx.key 2048
openssl req -new -key ./docker/nginx/ssl/nginx.key -out ./docker/nginx/ssl/nginx.csr
openssl x509 -req -in nginx.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out nginx.crt -days 3650 -sha256

echo "Nginx test certificates generated in ./docker/nginx/ssl/"
