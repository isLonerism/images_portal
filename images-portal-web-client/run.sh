#!/bin/sh

envsubst < /usr/share/nginx/html/config.js > /tmp/config.js
mv -f /tmp/config.js /usr/share/nginx/html/config.js
nginx -g 'daemon off;'
