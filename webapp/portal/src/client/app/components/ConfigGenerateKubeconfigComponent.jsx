import React from 'react';
import {Card, CardActions, CardHeader, CardTitle, CardText} from 'material-ui/Card';
import RaisedButton from 'material-ui/RaisedButton';
import {TextField} from "material-ui";

class ConfigGenerateKubeconfigComponent extends React.Component {
    constructor(props) {
        super(props);
        this.handleInputChange = this.handleInputChange.bind(this);

        this.state = this.props.config || {};
    }

    handleInputChange(event) {
        const target = event.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const name = target.name;

        this.setState({
            [name]: value
        }, () => this.props.onStateUpdate(this.state));
    }

    render() {
        return (
            <div>
                <TextField
                    name="server"
                    hintText="External URL for kubernetes API"
                    errorText={this.props.meta.server.errorText}
                    floatingLabelText="Kubernetes API URL"
                    floatingLabelFixed={true}
                    onChange={this.handleInputChange}
                    value={this.state.server}
                /><br />

                <TextField
                    name="name"
                    hintText="User friendly name for kubeconfig"
                    errorText={this.props.meta.name.errorText}
                    floatingLabelText="Kubeconfig profile name"
                    floatingLabelFixed={true}
                    onChange={this.handleInputChange}
                    value={this.state.name}
                />
            </div>
        );
    }

}

export default ConfigGenerateKubeconfigComponent;