
import React from 'react';
import {Card, CardActions, CardHeader, CardText} from 'material-ui/Card';
import RaisedButton from 'material-ui/RaisedButton';
import {TextField} from "material-ui";
import AuthConfigurations from "../../../../api/AuthConfigurations";
import FormComponent from "../../../../components/FormComponent.jsx";

class AuthConfigurationEditComponent extends FormComponent {

  constructor(props) {
    super(props, AuthConfigurations.build(), "default");
  }

  render() {
    if (!this.state.data) {
      return <div>Loading</div>;
    }

    return (
      <Card>
        <CardHeader
          title="General Settings"
        />
        <CardText>
          <form className="container" onSubmit={this.handleFormSubmit}>

            <div>{this.state.data.metadata.name}</div>

            {this.state.data.generateKubeconfig ?
              <div>
                <TextField
                  name="data.generateKubeconfig.server"
                  hintText="External URL for kubernetes API"
                  errorText={this.errorText('server')}
                  floatingLabelText="Kubernetes API URL"
                  floatingLabelFixed={true}
                  onChange={this.handleInputChange}
                  value={this.state.data.generateKubeconfig.server}
                /><br />

                <TextField
                  name="data.generateKubeconfig.name"
                  hintText="User friendly name for kubeconfig"
                  errorText={this.errorText('name')}
                  floatingLabelText="Kubeconfig profile name"
                  floatingLabelFixed={true}
                  onChange={this.handleInputChange}
                  value={this.state.data.generateKubeconfig.name}
                />
              </div>
              :
              <div>
                Will use default settings - click to override
              </div>
            }
          </form>
        </CardText>
        <CardActions>
          <RaisedButton
            label="Save"
            primary={true}
            onClick={this.handleFormSubmit}
          />
        </CardActions>
      </Card>
    );
  }

}

export default AuthConfigurationEditComponent;