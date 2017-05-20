import React from 'react';
import {List, ListItem, Subheader} from "material-ui";
import AuthProviders from "../api/AuthProviders";
import {Link} from "react-router-dom";

class AuthProviderListComponent extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      data: null,
    };
  }

  componentDidMount() {
    AuthProviders.namespace("kopeio-auth").list().then(json => {
      this.setState({
        data: json,
      });
    });
  }


  render() {
    if (!this.state.data) {
      return <div>Loading</div>;
    }

    return (
      <div>
        <List>
          <Subheader>Authentication Providers</Subheader>
          <ListItem
            primaryText="Add"
          />
          {
            this.state.data.items.map(function (o, i) {
              return <ListItem primaryText={o.metadata.name}
                               key={i}
                               containerElement={<Link to={`config/authproviders/${o.metadata.name}`}/>}/>;
            })
          }
        </List>
      </div>
    );
  }

}

export default AuthProviderListComponent;