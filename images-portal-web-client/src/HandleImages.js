import React from 'react';
import TagImage from './TagImage';
import { TextField } from '@material-ui/core';

class HandleImages extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            projectName: ''
        };

        this.taggedImageList = [];

        this.addTaggedImage = this.addTaggedImage.bind(this);
        this.handleProjctChange = this.handleProjctChange.bind(this);
    };

    // Get old image name and new image name
    addTaggedImage = (image, taggedImage) => {
        this.taggedImageList.push({
            old_image: {
                name: image
            },
            new_image: {
                name: window.ENV.OSFT_REGISTRY_PREFIX + '/' + taggedImage
            }
        });

        // Check if all images are tagged
        if (this.taggedImageList.length === this.props.imageList.length) {
            // Initialize tagged image list in parent component
            this.props.taggedImageList(this.taggedImageList);

            console.log('Finished tagging all elements');
            console.log (this.taggedImageList);
        };
    };

    handleProjctChange (e) {
        this.setState({
            projectName: e.target.value
        });
    }

    render() {
        var taggedImages = this.props.imageList.map(image =>
            <TagImage
                key={image.name}
                imageTag={image.name}
                projectName={this.state.projectName}
                taggedImage={this.addTaggedImage}>
            </TagImage> );
        return (
            <div>
                <TextField
                    required
                    id="outlined-required"
                    label="Project Name"
                    placeholder="Project Name"
                    margin="dense"
                    variant="outlined"
                    helperText="OpenShift Project Name (Namespace)"
                    onChange={this.handleProjctChange}
                />
                {taggedImages}
            </div>
        )
    }
}

export default HandleImages;