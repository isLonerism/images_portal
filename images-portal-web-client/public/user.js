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
        xhr.open('GET', window.ENV.OAUTH_OPENSHIFT_ROUTE + '/apis/project.openshift.io/v1/projects')
        xhr.setRequestHeader('Authorization', 'Bearer ' + window.USER.ACCESS_TOKEN)
        xhr.setRequestHeader('Accept', 'application/json')

        // API response callback
        xhr.onreadystatechange = function () {
            if (xhr.readyState == XMLHttpRequest.DONE) {
                projectList = JSON.parse(xhr.responseText)

                // map a list of user's projects
                window.USER.PROJECT_LIST = projectList['items'].map(function (project) {
                    return project['metadata']['name']
                })
            }
        }

        xhr.send()
    }

    // redirect to authentication screen if access token is not present
    else {
        let redirectUrl = window.ENV.OAUTH_OPENSHIFT_ROUTE + '/oauth/authorize?'

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