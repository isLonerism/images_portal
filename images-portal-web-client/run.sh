#!/bin/sh

envsubst < /usr/share/nginx/html/config.js > /usr/share/nginx/html/config.js
nginx -g 'daemon off;'
