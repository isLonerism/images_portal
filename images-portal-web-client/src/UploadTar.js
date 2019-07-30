import React from 'react';
import AWS from 'aws-sdk';
import Button from '@material-ui/core/Button';

class UploadTar extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            uploadState: 'UPLOAD TAR',
            disableUpload: false
        };

    this.uploadTarToS3 = this.uploadTarToS3.bind(this);
    this.uploadTar = this.uploadTar.bind(this);

/*     this.classes = makeStyles(theme => ({
        button: {
            margin: theme.spacing(1),
            color: 'blue',
        },
        input: {
            display: 'none',
        },
        })); */

    };

    // File input change event, updates upload state and call UploadTarToS3 function
    uploadTar (e) {
        this.setState({
            uploadState: 'UPLOADING...',
            disableUpload: true
        }, () => {
            this.uploadTarToS3();
        })
    }

    uploadTarToS3 (e) {

        // S3 necessary parameters
/*         const S3_BUCKET_NAME = 'testbucket';
        const S3_ACCESS_KEY_ID = 'ordavid';
        const S3_SECRET_ACCESS_KEY = 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY';
        const S3_ENDPOINT = 'http://172.17.0.2:9000'; */

        // Assign chosen file to element
        let tarFile = this.uploadInput.files[0];

        // Initialize new s3 service object using the s3 parameters and the chosen file
        let s3 = new AWS.S3 ({
            params: {
                Bucket: window.ENV.S3_BUCKET_NAME,
                Key: tarFile.name,
                Body: tarFile
            },
            credentials: {
                accessKeyId: window.ENV.S3_ACCESS_KEY_ID,
                secretAccessKey: window.ENV.S3_SECRET_ACCESS_KEY
            },
            endpoint: window.ENV.S3_ENDPOINT,
            s3ForcePathStyle: true,
            httpOptions: {
                timeout: window.ENV.S3_UPLOAD_TIMEOUT_MS
            }
        })

        // Upload the chosen file the bucket, render state according to outcome and update parent component upon success
        s3.putObject(s3.params, function(err, data) {
            if (err) {                
                this.setState({ uploadState: 'Upload Failed' });
                this.props.failStep('Upload Failed');
            }
            else {
                this.setState({ uploadState: tarFile.name });
                this.props.onTarUpload(tarFile.name);
                console.log(data);
            }
        }.bind(this));
    }

    render() {
        return (
            <div>
                <input
                    /* className={this.classes.input} */
                    id="uploadInput"
                    type="file"
                    accept=".tar"
                    onChange = {this.uploadTar}
                    disabled={this.state.disableUpload}
                    ref= { (ref) => { this.uploadInput = ref; }}
                    style={{ display: 'none' }}
                />
                <label htmlFor="uploadInput">
                    <Button 
                        size="small"
                        variant="contained"
                        color="primary"
                        component="span">
                    { this.state.uploadState }
                    </Button>
                </label>
            </div>
        )
    }
}

export default UploadTar;