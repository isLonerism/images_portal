import React from 'react';
import Typography from '@material-ui/core/Typography';
import fetchWithTimeout from './fetchWithTimeout.js';
import CircularProgress from '@material-ui/core/CircularProgress';

class HandleTar extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null
        };

        this.loadTar = this.loadTar.bind(this);
    };

    // automatically start loading tar on step 2
    componentDidMount() {
        this.loadTar();
    }

    loadTar() {
        // const ENDPOINT = 'http://127.0.0.1:8786/load'
        // JSON: s3key: images.tar
        //const ENDPOINT = 'http://images-portal-grpc-client.router.myproject.svc.cluster.local/load';

        console.log('Sending request to load tar:');
        console.log(window.ENV.IMAGE_LOAD_ROUTE);
        /*         console.log({
                    s3key: this.props.s3key
                }); */
        //console.log(JSON.stringify({s3key:"statsd-exporter-vault.tar"}));

        fetchWithTimeout(window.ENV.IMAGE_LOAD_ROUTE, {
            method: 'POST',
            //headers: {    
            /* "Content-type": "application/json; charset=UTF-8;" */
            //"Content-type": "application/json;"
            //},
            //s3key: this.props.s3key
            //body: JSON.stringify('{"s3key":"statsd-exporter-vault.tar"}')
            body: JSON.stringify({
                s3key: this.props.s3key
            })
        }, window.ENV.LOAD_TIMEOUT_MS).then(
            response => response.json()
            //response => console.log(response)
        ).then(
            success => {
                // Temporary - delete after tests
                //var lala = JSON.parse('{"Token":"eyJQb2ROYW1lIjoiMjAxOTA2MjUtMTMxOTAwLTAiLCJQb2RJUCI6IjE3Mi4xNy4wLjgifQ==","Images":[{"name":"prom/statsd-exporter:v0.5.0"},{"name":"quay.io/coreos/vault:0.9.1-0"}]}');
                //this.props.podToken(lala.Token);
                this.props.podToken(success.Token);
                //this.props.tarImageList(lala.Images);
                this.props.tarImageList(success.Images);
                console.log(success);
            }
        ).catch(
            error => {
                this.setState({
                    error: error.message
                });
                this.props.failStep('Loading Tar to Docker Failed');
                console.log(error.message);
            }
        );
    };

    render() {
        if (this.state.error != null) {
            return (
                <Typography color="error">
                    AN ERROR OCCURED WHILE LOADING IMAGES
                    <div>
                        {this.state.error}
                    </div>
                </Typography>
            )
        }

        return (
            <div>
                <CircularProgress />
                <Typography variant="subtitle2" color='textSecondary'>
                    Loading your images...
                </Typography>
            </div>
        )
    }
}

export default HandleTar;