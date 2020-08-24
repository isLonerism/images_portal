# Images Portal: Minio

Deploy an instance of minio with pre-configured users and permissions for images portal.

## Configuration

### Environment variables

#### Required

* BUCKET_NAME: name for your bucket (use value from S3 configmap)
* BUCKET_LIFECYCLE_DAYS: how long (in days) uploaded files stay in the bucket
* WEB_CLIENT_S3_ACCESS_KEY: use value from web client secret
* WEB_CLIENT_S3_SECRET_KEY: use value from web client secret
* GRPC_SERVER_S3_ACCESS_KEY: use value from gRPC server secret
* GRPC_SERVER_S3_SECRET_KEY: use value from gRPC server secret

#### Optional

* MINIO_ACCESS_KEY: access key for admin user
* MINIO_SECRET_KEY: secret key for admin user

### Volumes

#### Required

* /data: volume for S3 object storage

## Procedure

1. Build and push the image from provided Dockerfile
2. Deploy image with the required configuration described above
3. Change the value of S3 host within the S3 configmap to the service of this deployment
