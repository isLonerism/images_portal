import React from 'react';
import fetchWithTimeout from './fetchWithTimeout';
import Typography from '@material-ui/core/Typography';
import CircularProgress from '@material-ui/core/CircularProgress';

class PushImages extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null
        };

        this.pushImages = this.pushImages.bind(this);
    };

    componentDidMount() {
        this.pushImages()
    }

    pushImages() {

        let postBody = {
            PodToken: this.props.podToken,
            DockerToken: window.USER.ACCESS_TOKEN,
            images: {
                images: this.props.imageList
            }
        }

        console.log('Sending request to push images:');
        console.log(window.ENV.IMAGE_PUSH_ROUTE);
        console.log(postBody);

        fetchWithTimeout(window.ENV.IMAGE_PUSH_ROUTE, {
            method: 'POST',
            body: JSON.stringify(postBody)
        }, window.ENV.PUSH_TIMEOUT_MS).then(
            response => response.json()
        ).then(
            success => {
                this.props.pushSuccess();

                console.log('successful push!')
                console.log(success);
            }
        ).catch(
            error => {
                this.setState({
                    error: error.message
                });
                this.props.failStep('Failed to Push Images');
                console.log(error);
            }
        );
    };

    render() {
        if (this.state.error != null) {
            return (
                <Typography color="error">
                    AN ERROR OCCURED WHILE PUSHING IMAGES
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
                    Pushing your images...
                </Typography>
            </div>
        )
    }
}

export default PushImages;