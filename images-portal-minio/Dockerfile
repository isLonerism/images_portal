FROM minio/minio:RELEASE.2020-03-19T21-49-00Z
RUN apk add gettext
COPY images-portal /images-portal
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["/images-portal/docker-entrypoint.sh"]
