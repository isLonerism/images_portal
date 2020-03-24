import React from 'react';
import Input from '@material-ui/core/Input';
import Checkbox from '@material-ui/core/Checkbox';
import Select from '@material-ui/core/Select';

class TagImage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            project: this.props.projectList[0],
            name: this.getNameFromFullName(this.props.image),
            tag: this.getTagFromFullName(this.props.image),
            isTagged: false,
        };

        this.handleTagChange = this.handleTagChange.bind(this)
        this.getTagFromFullName = this.getTagFromFullName.bind(this)
        this.getNameFromFullName = this.getNameFromFullName.bind(this)
    };

    getTagFromFullName = (fullname) => {
        return fullname.substring(fullname.lastIndexOf(":") + 1)
    }

    getNameFromFullName = (fullname) => {
        return fullname.substring(fullname.lastIndexOf("/") + 1, fullname.lastIndexOf(":"))
    }

    handleTagChange() {
        this.setState({
            isTagged: !this.state.isTagged
        }, () => {
            var newImage = this.state.project + "/" + this.state.name + ":" + this.state.tag
            this.props.handleTagChange(this.props.image, newImage, this.state.isTagged)
        })
    }

    render() {
        return (
            <div>
                <Select
                    native
                    style={{ width: '15%' }}
                    value={this.state.project}
                    disabled={this.state.isTagged}
                    onChange={(event) => {
                        this.setState({
                            project: event.target.value
                        })
                    }}
                >
                    {this.props.projectList.map(project =>
                        <option key={project} value={project}>{project}</option>)}
                </Select>
                <span style={{ marginLeft: "12px", marginRight: "12px" }}>/</span>
                <Input
                    style={{ width: '40%' }}
                    type='text'
                    disabled={this.state.isTagged}
                    placeholder={this.getNameFromFullName(this.props.image)}
                    onChange={(event) => { this.setState({ name: event.target.value }) }}
                    defaultValue={this.state.name}>
                </Input>
                <span style={{ marginLeft: "12px", marginRight: "12px" }}>:</span>
                <Input
                    style={{ width: '15%' }}
                    type='text'
                    disabled={this.state.isTagged}
                    placeholder={this.getTagFromFullName(this.props.image)}
                    onChange={(event) => { this.setState({ tag: event.target.value }) }}
                    defaultValue={this.state.tag}>
                </Input>
                <Checkbox
                    color="primary"
                    checked={this.state.isTagged}
                    onClick={this.handleTagChange}>
                </Checkbox>
            </div>
        )
    }
}

export default TagImage;