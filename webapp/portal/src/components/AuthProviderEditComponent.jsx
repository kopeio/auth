import React from 'react';
import {Card, CardActions, CardHeader, CardText} from 'material-ui/Card';
import RaisedButton from 'material-ui/RaisedButton';
import {TextField} from "material-ui";
import AuthProviders from "../api/AuthProviders";
import FormComponent from "./FormComponent.jsx";
import AuthProviderListComponent from "./AuthProviderListComponent.jsx";

class AuthProviderEditComponent extends FormComponent {

  constructor(props) {
    super(props, AuthProviders.namespace("kopeio-auth"), props.name);
  }

  render() {
    if (!this.state.data) {
      return <div>Loading</div>;
    }

    if (!this.state.data.oAuthConfig) {
      this.state.data.oAuthConfig = {};
    }

    return (
      <Card>
        <CardHeader
          title="Settings"
        />
        <CardText>
          <form className="container" onSubmit={this.handleFormSubmit}>

            <div>{this.state.data.metadata.name}</div>

            <div>
              <TextField
                name="data.description"
                hintText="Human-readable description"
                errorText={this.errorText('data.description')}
                floatingLabelText="Description"
                floatingLabelFixed={true}
                onChange={this.handleInputChange}
                value={this.state.data.description}
              /><br />

              <TextField
                name="data.oAuthConfig.clientID"
                hintText="ClientID"
                errorText={this.errorText('data.oAuthConfig.clientID')}
                floatingLabelText="Client ID"
                floatingLabelFixed={true}
                onChange={this.handleInputChange}
                value={this.state.data.oAuthConfig.clientID}
              /><br />

              <TextField
                name="data.oAuthConfig.clientSecret"
                hintText="clientSecret"
                errorText={this.errorText('data.oAuthConfig.clientSecret')}
                floatingLabelText="Client Secret"
                floatingLabelFixed={true}
                onChange={this.handleInputChange}
                value={this.state.data.oAuthConfig.clientSecret}
              /><br />


              <TextField
                name="data.permitEmails"
                hintText="permitEmails"
                errorText={this.errorText('data.permitEmails')}
                floatingLabelText="Emails permitted to authenticate"
                floatingLabelFixed={true}
                rows={2}
                onChange={this.handleInputChange}
                value={this.state.data.permitEmails}
              /><br />

            </div>
          </form>

          <AuthProviderListComponent />
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

export default AuthProviderEditComponent;