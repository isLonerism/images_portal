window.USER = {
    // namespace for user related variables
}

window.addEventListener('load', () => {

    // access token is passed in the anchor field for some reason
    let queryString = window.location.hash.replace(/^#/, '?')
    let urlParams = new URLSearchParams(queryString)

    // populate USER namespace
    if (urlParams.has('access_token')) {
        window.USER.ACCESS_TOKEN = urlParams.get('access_token')

        // prepare API request for user's projects
        let xhr = new XMLHttpRequest()
        xhr.open('POST', window.ENV.PROJECTS_REQUEST_ROUTE)

        // API response callback
        xhr.onreadystatechange = function () {
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (xhr.status == 200) {
                    window.USER.PROJECT_LIST = JSON.parse(xhr.responseText)["ProjectList"]
                }
                else {
                    console.log("Could not get user's projects")
                }
            }
        }

        xhr.send(JSON.stringify({
            APIEndpoint: window.ENV.OPENSHIFT_API_ENDPOINT,
            Token: window.USER.ACCESS_TOKEN
        }))
    }

    // redirect to authentication screen if access token is not present
    else {
        let redirectUrl = window.ENV.OPENSHIFT_OAUTH_ENDPOINT + '/oauth/authorize?'

        // parameters passed to OpenShift OAuth server
        let redirectParams = new URLSearchParams({
            client_id: window.ENV.OAUTH_CLIENT_ID,
            redirect_uri: window.location.origin,
            response_type: 'token'
        })

        // redirect to OAuth server
        window.location.replace(redirectUrl + redirectParams.toString())
    }
})