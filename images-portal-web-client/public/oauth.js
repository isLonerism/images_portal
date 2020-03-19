window.addEventListener('load', () => {

    // access token is passed in the anchor field for some reason
    let queryString = window.location.hash.replace(/^#/, '?')
    let urlParams = new URLSearchParams(queryString)

    // store user's access token in ENV
    if (urlParams.has('access_token')) {
        window.ENV.USER_ACCESS_TOKEN = urlParams.get('access_token')
    }
    else {
        let redirectUrl = window.ENV.OAUTH_OPENSHIFT_ROUTE + ':8443/login?then=' + encodeURIComponent('/oauth/authorize?')

        // parameters passed to OpenShift OAuth server
        let redirectParams = {
            client_id: window.ENV.OAUTH_CLIENT_ID,
            redirect_uri: window.location.origin,
            response_type: 'token'
        }

        // encode entire queries, including '=' and '&'
        for (let param in redirectParams) {
            redirectUrl += encodeURIComponent(param + "=" + redirectParams[param] + '&')
        }

        // redirect to OAuth server (remove trailing ampersand)
        window.location.replace(redirectUrl.slice(0, -3))
    }
})