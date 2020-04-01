# Images Portal

Load image archives to internal registry directly without the additional work.

## Pre-requisites

* Create a new project specifically for images portal
* Give privileged security context to the default service account on this project
* Create a configmap for S3 bucket data with the following keys:
  * S3 host
  * Bucket name
  * Region
* Create 2 separate secrets, one for web client S3 user and one for gRPC server S3 user, each containing the following keys:
  * Access key
  * Secret key

## Procedure

1. (OPTIONAL) [images-portal-minio/](Deploy an instance of minio) if you do not have any object storage services available
2. [images-portal-grpc-server](Take care of gRPC server pre-requisites)
3. [images-portal-grpc-client](Deploy the gRPC client)
4. [images-portal-web-client](Deploy the web client)
