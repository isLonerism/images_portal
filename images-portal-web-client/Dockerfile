# build environment
FROM node:12.4.0-alpine as build
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY package.json /app/package.json
RUN npm install --silent
#RUN npm install
RUN npm install react-scripts@3.0.1 -g --silent
COPY . /app
RUN npm run build

# production environment
FROM nginx:1.16.0-alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY run.sh /tmp
COPY ca-certificates-20190108-r0.apk /tmp/ca-certificates-20190108-r0.apk
RUN sed -i s/"listen .*80"/"listen 3200"/ /etc/nginx/conf.d/default.conf \
    && chmod -R 777 /var/run /var/log/nginx /var/cache/nginx/ /etc/nginx/* \
    && touch repo.list \
    && apk add --repositories-file=repo.list --allow-untrusted --no-network --no-cache /tmp/ca-certificates-20190108-r0.apk \
    && update-ca-certificates
EXPOSE 3200
CMD ["/tmp/run.sh"]

