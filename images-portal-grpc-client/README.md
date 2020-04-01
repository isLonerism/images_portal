# Images Portal: gRPC client

gRPC client handles requests coming from the web client and interacts with the gRPC server.

## Configuration

### Environment variables

#### Required

* POD_TEMPLATE_PATH: path to the pod_template.yml (e.g. /data/pod_template.yml)
* PROJECT: name of the project gRPC client is deployed in (e.g. images-portal)
* GRPC_PORT: port exposed by the gRPC server (e.g. 7777)
* YAML_DECODER_BUFFER_SIZE: buffer size for the YAML decoder of pod_template.yml
* POD_STARTUP_TIME: time (in seconds) to wait for a pod to go up
* S3Region: (use S3 region value from S3 data configmap)
* S3Endpoint: (use S3 host value from S3 data configmap)
* S3Bucket: (use bucket name value from S3 data configmap)
* S3Accesskey: (use access key value from gRPC secret)
* S3Secretkey: (use secret key value from gRPC secret)

### Volumes

#### Required

* pod_template.yml: configmap which should be mounted at the directory POD_TEMPLATE_PATH points to (e.g. /data)

## Procedure

1. Get all go dependencies and build the server executable
2. Build and push the image from provided Dockerfile
3. Deploy image with the required configuration described above
