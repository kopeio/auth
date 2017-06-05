import React from 'react';
import {Card, CardActions, CardTitle} from 'material-ui/Card';
import RaisedButton from 'material-ui/RaisedButton';
import {Link} from "react-router-dom";
import './UserComponent.css';

class UserComponent extends React.Component {
  render() {
    if (this.props.user == null) {
      return (
        <Card className="UserComponent">
          <CardTitle title="Not logged in"/>
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
      <Card className="UserComponent">
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