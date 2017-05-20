import React from 'react';
import {Card, CardActions, CardHeader, CardTitle, CardText} from 'material-ui/Card';
import RaisedButton from 'material-ui/RaisedButton';
import {List, ListItem, TextField, Toggle} from "material-ui";
import AuthProviders from "../api/AuthProviders";
import {Link} from "react-router-dom";

class AuthProviderListComponent extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      data: null,
      open: true
    };

    this.handleNestedListToggle = this.handleNestedListToggle.bind(this);
  }

  componentDidMount() {
    AuthProviders.namespace("default").list().then(json => {
      this.setState({
        data: json,
      });
    });
  }

  handleNestedListToggle(item) {
    console.log("state", this.state);
    this.setState({
      open: !this.state.open,
    });
  };

  render() {
    console.log("A");
    if (!this.state.data) {
      return <div>Loading</div>;
    }

    return (
      <div>
        <List>
          <ListItem
            primaryText="Add"
          />
          <ListItem
            primaryText="Authentication Providers"
            open={this.state.open}
            onNestedListToggle={this.handleNestedListToggle}
            primaryTogglesNestedList={true}
            nestedItems={
              this.state.data.items.map(function (o, i) {
                return <ListItem primaryText={o.metadata.name}
                                 key={i}
                                 containerElement={<Link to={`config/authproviders/${o.metadata.name}`}/>}/>;
              })
            }
          />
        </List>
      </div>
    );
  }

}

export default AuthProviderListComponent;