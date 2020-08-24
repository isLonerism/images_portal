# Images Portal: gRPC server

This image is used to deploy a DinD pod to load, tag and push images to the internal registry.

## Configuration

### Environment variables

#### Required

* DOCKER_HOST: host the gRPC server should query for docker commands (already present in the template and should not be changed)
* DOCKER_API_VERSION: api version for gRPC client to use to interact with the server (already present in the template and should not be changed)
* DELETE_OBJECT_AFTER_LOAD: should the server delete the object from object storage after loading it ('true' by default)

### Volumes

#### Required

* /var/lib/docker: ephemeral/persistent storage for docker images

## Procedure

1. Build and push the image from provided Dockerfile
2. Change the name of the image within the pod_template.yml to your uploaded image
3. (OPTIONAL) Change the configuration (described above) within the pod_template.yml to suit your needs
3. Create a configmap out of pod_template.yml
