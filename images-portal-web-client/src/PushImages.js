import React from 'react';
import { TextField } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import fetchWithTimeout from './fetchWithTimeout';

class PushImages extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            pushState: 'PUSH IMAGES TO OPENSHIFT REGISTRY',
            disablePushBtn: false,
            osftToken: ''
        };

        this.handleTokenChange = this.handleTokenChange.bind(this);
        this.pushImages = this.pushImages.bind(this);
        this.handlePush = this.handlePush.bind(this);
    };

    handlePush (e) {
        this.setState({
            pushState: 'LOADING...',
            disablePushBtn: true
        }, () => {
            this.pushImages();
        })
    }

    pushImages () {
        // const ENDPOINT = 'http://127.0.0.1:8786/load'
        // JSON: s3key: images.tar
        //const ENDPOINT = 'https://images-portal-grpc-client.router.myproject.svc.cluster.local/push';

        let postBody = {
            PodToken: this.props.podToken,
            DockerToken: this.state.osftToken,
            images: {
                images: this.props.imageList
            }
        }

        console.log('Sending request to push images:');
        console.log(window.ENV.IMAGE_PUSH_ROUTE);
        console.log(postBody);


/*         let postBody = '{"PodToken":"eyJQb2ROYW1lIjoiMjAxOTA2MjUtMTQ1OTE4LTAiLCJQb2RJUCI6IjE3Mi4xNy4wLjExIn0=","DockerToken":"eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJteXByb2plY3QiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoiYnVpbGRlci10b2tlbi1tNTlxaiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJidWlsZGVyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiOTc4ODM3NzYtOTcyOS0xMWU5LWE3MGYtNTI1NDAwYjhiNjVhIiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50Om15cHJvamVjdDpidWlsZGVyIn0.jmueDWbRKmEDQdvpdaBz0zWa003WHJ2Er6kNTgFWlfixl0VRKKo4VPXbWST4knCra8WY2unAq-ShJCIeNBmLgBmS5aXREBIe4lV6MyYAxvLfJyTe75Ep3gRKzH4VMl_rQ8uyzDIDWtXYR3GAjBCwTeeQDjVHWQX683Zj-br3HNNpahNRZIqCqtIIL47UjRPG1HG1XTiHYQL6tSCNao65skoVvtbMsnRX53Hp_qQRE4JF7CZVA-7HYRpMWyYrDe_cOvxvJUF94zHC2ApytuBeyELEeYnhv-Zo2x1oJe2z2x_BltrlRF4hCK8AUQGW69FR9BryiR7kKFuuQa2D8lro_Q","images":{"images":[{"old_image":{"name":"prom/statsd-exporter:v0.5.0"},"new_image":{"name":"docker-registry.default.svc:5000/myproject/prom:prom-tag"}},{"old_image":{"name":"quay.io/coreos/vault:0.9.1-0"},"new_image":{"name":"docker-registry.default.svc:5000/myproject/vault:vault-tag"}}]}}'
        console.log(postBody); */

        fetchWithTimeout(window.ENV.IMAGE_PUSH_ROUTE, {
            method:'POST',
            //headers: {
                //"Content-type": "application/json; charset=UTF-8"
            //},
            body: JSON.stringify(postBody)
        }, window.ENV.PUSH_TIMEOUT_MS).then(
            response => response.json()
        ).then(
            success => {
                this.setState({
                    pushState: 'FINISHED PUSH',
                });
                this.props.pushSuccess();

                console.log('successful push!')
                console.log(success);
            }
        ).catch(
            error =>  {
                this.setState({
                    pushState: 'PUSHING IMAGES TO OPENSHIFT REGISTRY FAILED'
                });
                this.props.failStep('Failed to Push Images');
                console.log(error);
            }
        );
    };

    handleTokenChange (e) {
        this.setState({
            osftToken: e.target.value
        });
    }

    render() {
        return (
            <div>
                <TextField
                    variant="outlined"
                    multiline
                    rows="4"
                    required
                    label="OpenShift Token"
                    helperText="The OpenShift registry must be authenticated using either a user session or service account token."
                    placeholder="OpenShift Token"
                    margin="normal"
                    onChange={this.handleTokenChange}
                />
                <br></br>
                <br></br>
                <Button
                    size="small"
                    variant="contained"
                    color="primary"
                    component="span"
                    onClick = {this.handlePush}
                    disabled = {this.state.osftToken === '' || this.state.disablePushBtn}>
                    {this.state.pushState}
                </Button>
            </div>
        )
    }
}

export default PushImages;