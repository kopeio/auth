import React from 'react';
import {Card, CardActions, CardHeader, CardTitle} from 'material-ui/Card';
import RaisedButton from 'material-ui/RaisedButton';

class UserComponent extends React.Component {

    render() {
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
                <CardTitle title="Logged in" subtitle={this.props.user.name}/>
                <CardActions>
                    <RaisedButton
                        label="Download kubecfg file"
                        primary={true}
                        href="/portal/actions/kubecfg"
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