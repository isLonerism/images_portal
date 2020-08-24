import React from 'react';
import TagImage from './TagImage';

class HandleImages extends React.Component {
    constructor(props) {
        super(props);
        this.state = {};

        this.taggedImageList = []

        this.handleTagChange = this.handleTagChange.bind(this);
    };

    // Get old image name and new image name
    handleTagChange = (oldImage, newImage, tagged) => {
        var taggedImage = {
            old_image: {
                name: oldImage
            },
            new_image: {
                name: window.ENV.OSFT_REGISTRY_PREFIX + '/' + newImage
            }
        }

        // Add or remove image from list based on checkbox status
        if (!tagged) {
            this.taggedImageList.splice(this.taggedImageList.indexOf(taggedImage), 1)
        } else {
            this.taggedImageList.push(taggedImage)

            // Check if all images are tagged
            if (this.taggedImageList.length === this.props.imageList.length) {
                // Initialize tagged image list in parent component
                this.props.taggedImageList(this.taggedImageList);

                console.log('Finished tagging all elements');
                console.log(this.taggedImageList);
            }
        }
    }

    render() {
        return (
            <div>
                {this.props.imageList.map(image =>
                    <TagImage
                        key={image.name}
                        image={image.name}
                        projectList={window.USER.PROJECT_LIST}
                        handleTagChange={this.handleTagChange}>
                    </TagImage>)}
            </div>
        )
    }
}

export default HandleImages;