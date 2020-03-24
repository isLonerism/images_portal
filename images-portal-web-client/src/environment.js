
export default function initENVNamespace() {
    window.ENV = {
        "S3_BUCKET_NAME": process.env.REACT_APP_S3_BUCKET_NAME,
        "S3_ACCESS_KEY_ID": process.env.REACT_APP_S3_ACCESS_KEY_ID,
        "S3_SECRET_ACCESS_KEY": process.env.REACT_APP_S3_SECRET_ACCESS_KEY,
        "S3_ENDPOINT": process.env.REACT_APP_S3_ENDPOINT,
        "S3_UPLOAD_TIMEOUT_MS": process.env.REACT_APP_S3_UPLOAD_TIMEOUS_MS,
        "IMAGE_PUSH_ROUTE": process.env.REACT_APP_IMAGE_PUSH_ROUTE,
        "IMAGE_LOAD_ROUTE": process.env.REACT_APP_IMAGE_LOAD_ROUTE,
        "OSFT_REGISTRY_PREFIX": process.env.REACT_APP_OSFT_REGISTRY_PREFIX,
        "LOAD_TIMEOUT_MS": process.env.REACT_APP_LOAD_TIMEOUT_MS,
        "PUSH_TIMEOUT_MS": process.env.REACT_APP_PUSH_TIMEOUT_MS,
        "OAUTH_OPENSHIFT_ROUTE": process.env.REACT_APP_OAUTH_OPENSHIFT_ROUTE,
        "OAUTH_CLIENT_ID": process.env.REACT_APP_OAUTH_CLIENT_ID,
    }
}
