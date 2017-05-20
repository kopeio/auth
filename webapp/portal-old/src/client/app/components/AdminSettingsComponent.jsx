import React from 'react';
import {Card, CardActions, CardHeader, CardTitle} from 'material-ui/Card';
import RaisedButton from 'material-ui/RaisedButton';

class UserComponent extends React.Component {

    render() {


      ProviderSpec `json:"authProviders,omitempty""`

      GenerateKubeconfig *GenerateKubeconfig `json:"generateKubeconfig,omitempty"`
    }

  type AuthProviderSpec struct {
  // ID is a system-friendly identifier
  ID string `json:"id,omitempty"`

  // Name is a human-friendly name
  Name string `json:"name,omitempty"`

  OAuthConfig *OAuthConfig `json:"oAuthConfig,omitempty"`

  // Email addresses that are allowed to register using this provider
  PermitEmails []string `json:"permitEmails,omitempty"`
}

type OAuthConfig struct {
  ClientID string `json:"clientID,omitempty"`

  // TODO(authprovider-q): What do we do about secrets?  We presumably don't want this secret
  // in the configmap, because that might have a fairly permissive RBAC role.  But do we want to
  // do a layerable configuration?  Keep the secret in a second configuration object?  Have the
  // name of the secret here, and just runtime error until the secret is loaded?

  // ClientSecret is the OAuth secret
  ClientSecret string `json:"clientSecret,omitempty"`
}

type GenerateKubeconfig struct {
  Server string `json:"server,omitempty"`
  Name   string `json:"name,omitempty"`
}


if (this.props.user == null) {
            return (
                <Card>>
                    <CardTitle title="Not logged in" />
                    <CardActions>
                        <RaisedButton
                            label="Login with google"
                            primary={true}
                            href="/portal/actions/login?provider=google"
                        />
                    </CardActions>
                </Card>
            );
        }

        return (
            <Card>>
                <CardTitle title="Logged in" subtitle={this.props.user.username}/>
                <CardActions>
                    <RaisedButton
                        label="Download kubecfg file"
                        primary={true}
                        href="/portal/actions/kubeconfig"
                    />
                    <RaisedButton
                        label="Logout"
                        secondary={true}
                        href="/portal/actions/logout"
                    />
                </CardActions>
            </Card>
        );
    }

}

export default UserComponent;