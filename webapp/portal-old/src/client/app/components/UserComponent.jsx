import React from 'react';
import {Card, CardActions, CardHeader, CardTitle} from 'material-ui/Card';
import RaisedButton from 'material-ui/RaisedButton';
import {Link} from "react-router-dom";

class UserComponent extends React.Component {

    render() {
        if (this.props.user == null) {
            return (
                <Card>
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
            <Card>
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

                    <RaisedButton
                      label="Settings"
                      secondary={true}
                      containerElement={<Link to="/config"/>}
                    />

                  <li><Link to="/config">Config</Link></li>

                </CardActions>
            </Card>
        );
    }

}

export default UserComponent;