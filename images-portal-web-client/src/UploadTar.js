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

    };

    // File input change event, updates upload state and call UploadTarToS3 function
    uploadTar(e) {
        this.setState({
            uploadState: 'UPLOADING...',
            disableUpload: true
        }, () => {
            this.uploadTarToS3();
        })
    }

    uploadTarToS3(e) {

        // Assign chosen file to element
        let tarFile = this.uploadInput.files[0];

        // Import async module for asynchronous iterations
        let async = require('async')

        // Initialize new s3 service object using the s3 parameters and the chosen file
        let s3 = new AWS.S3({
            params: {
                Bucket: window.ENV.S3_BUCKET_NAME,
                Key: tarFile.name,
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

        // Initiate a multipart file upload
        s3.createMultipartUpload({}, function (err, multipart) {
            if (err) {
                console.log(err)
                this.setState({ uploadState: 'Upload Failed' });
                this.props.failStep('Upload Failed');
            }
            else {
                // Chunk size is 128 MB as per the minio standard
                let chunkSize = Math.min(1024 * 1024 * 128, tarFile.size)
                let parts = Math.ceil(tarFile.size / chunkSize)

                // Send each chunk in an async loop
                async.timesSeries(parts, (partNum, next) => {
                    let rangeStart = partNum * chunkSize
                    let rangeEnd = Math.min(rangeStart + chunkSize, tarFile.size)

                    // s3 expects the part index to start from 1
                    partNum++

                    s3.uploadPart({
                        Body: tarFile.slice(rangeStart, rangeEnd),
                        PartNumber: partNum,
                        UploadId: multipart.UploadId
                    }, function (err, data) {
                        next(err, {
                            ETag: data.ETag,    // each part upload returns its ETag
                            PartNumber: partNum
                        })
                    })
                }, function (err, dataPacks) {
                    s3.completeMultipartUpload({
                        MultipartUpload: {
                            Parts: dataPacks    // map of ETags and indices for each part
                        },
                        UploadId: multipart.UploadId
                    }, function (err, data) {
                        if (err) {
                            console.log(err)
                            this.setState({ uploadState: 'Upload Failed' });
                            this.props.failStep('Upload Failed');
                        }
                        else {
                            this.setState({ uploadState: tarFile.name });
                            this.props.onTarUpload(tarFile.name);
                            console.log(data);
                        }
                    }.bind(this))
                }.bind(this))
            }
        }.bind(this))
    }

    render() {
        return (
            <div>
                <input
                    id="uploadInput"
                    type="file"
                    accept=".tar"
                    onChange={this.uploadTar}
                    disabled={this.state.disableUpload}
                    ref={(ref) => { this.uploadInput = ref; }}
                    style={{ display: 'none' }}
                />
                <label htmlFor="uploadInput">
                    <Button
                        size="small"
                        variant="contained"
                        color="primary"
                        component="span">
                        {this.state.uploadState}
                    </Button>
                </label>
            </div>
        )
    }
}

export default UploadTar;