# Images Portal: Web-client

Interface for users to use to upload their images.

## Configuration

### Environment variables

#### Required

* S3_BUCKET_NAME: (use bucket name value from S3 data configmap)
* S3_ACCESS_KEY_ID: (use access key value from web client secret)
* S3_SECRET_ACCESS_KEY: (use secret key value from web client secret)
* S3_ENDPOINT: (use S3 host value from S3 data configmap)
* S3_UPLOAD_TIMEOUS_MS: timeout (ms) for S3 upload requests
* IMAGE_PUSH_ROUTE: route URL to send the push request to
* IMAGE_LOAD_ROUTE: route URL to send the load request to
* OSFT_REGISTRY_PREFIX: route/service URL to the internal image registry
* LOAD_TIMEOUT_MS: timeout (ms) for image load request
* PUSH_TIMEOUT_MS: timeout (ms) for image push request
* OPENSHIFT_API_ENDPOINT: this endpoint will be queried for user projects
* OPENSHIFT_OAUTH_ENDPOINT: this endpoint will be queried for authentication token
* PROJECTS_REQUEST_ROUTE: gRPC route to proxy the projects request through
* OAUTH_CLIENT_ID: client ID of the created OAuthClient object

## Procedure

1. Build and push the image from provided Dockerfile
2. Deploy image with the required configuration described above
3. Expose the web client using a route
4. Add the route URL to "redirectURIs" list within oauth_client.yml and create an OAuthClient object out of oauth_client.yml
