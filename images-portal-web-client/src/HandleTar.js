import React from 'react';
import Button from '@material-ui/core/Button';
import fetchWithTimeout from './fetchWithTimeout.js';

class HandleTar extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            loadState: 'LOAD TAR TO DOCKER',
            disableTarLoadBtn: false,
        };

        this.handleTar = this.handleTar.bind(this);
        this.loadTar = this.loadTar.bind(this);
    };

    handleTar (e) {
        this.setState({
            loadState: 'LOADING...',
            disableTarLoadBtn: true
        }, () => {
            this.loadTar();
        })
    }

    loadTar () {
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
            method:'POST',
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
                this.setState({
                    loadState: 'FINISHED LOADING',
                });
                //this.props.tarImageList(lala.Images);
                this.props.tarImageList(success.Images);
                console.log(success);
            }
        ).catch(
            error =>  {
                this.setState({
                    loadState: 'LOADING TAR TO DOCKER FAILED',
                });
                this.props.failStep('Loading Tar to Docker Failed');
                console.log(error.message);
            }
        );
    };

    render() {
        return (
            <div>
                <Button
                    size="small"
                    variant="contained"
                    color="primary"
                    component="span"
                    onClick = {this.handleTar}
                    disabled = {this.state.disableTarLoadBtn}>
                    {this.state.loadState}
                </Button>
            </div>
        )
    }
}

export default HandleTar;