import React from 'react';
import Input from '@material-ui/core/Input';
import Checkbox from '@material-ui/core/Checkbox';
//import { TextField } from '@material-ui/core';

class TagImage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            isTagged: false,
            newImageTag: this.props.imageTag.substring(this.props.imageTag.lastIndexOf("/")+1)
        };

        this.confirmTag = this.confirmTag.bind(this);
        this.handleChange = this.handleChange.bind(this);
    };

    confirmTag (e) {
        this.setState({
            isTagged: true
        }, () => {
            this.props.taggedImage(
                this.props.imageTag,
                this.props.projectName + '/' + this.state.newImageTag
            );
        });
    };

    handleChange (e) {
        this.setState({
            newImageTag: e.target.value
        });
    };

    render() {
        return (
            <div>
                {this.props.projectName}/
                <Input
                    style ={{width: '50%'}}
                    type='text'
                    defaultValue={this.props.imageTag.substring(this.props.imageTag.lastIndexOf("/")+1)}
                    //defaultValue={this.props.imageTag}
                    disabled={this.props.projectName === '' || this.state.isTagged}
                    onChange={this.handleChange}>
                </Input>
                <Checkbox
                    color="primary"
                    onClick={this.confirmTag}
                    disabled={this.props.projectName === '' || this.state.isTagged}>
                </Checkbox>
            </div>
        )
    }
}

export default TagImage;