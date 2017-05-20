import React, {Component} from 'react';
import UserComponent from "./components/UserComponent";
import {MuiThemeProvider} from "material-ui/styles";
import {BrowserRouter, Link, Route} from "react-router-dom";

import './App.css';
import AuthConfigurationEditComponent from "./components/AuthConfigurationEditComponent";
import AuthProviderEditComponent from "./components/AuthProviderEditComponent";
import {Divider, ListItem} from "material-ui";
import {ActionHome, ActionSettings} from "material-ui/svg-icons/index";
import AuthProviderListComponent from "./components/AuthProviderListComponent";


class App extends Component {
  render() {
    return (
      <BrowserRouter>
        <div>
          <MuiThemeProvider>
            <div>


              <div className="App-main">
                <div className="App-sidebar">
                  <ListItem primaryText="Home" leftIcon={<ActionHome />} containerElement={<Link to="/">A</Link>} />
                  <Divider />
                  <ListItem primaryText="Settings" leftIcon={<ActionSettings />} containerElement={<Link to="/config">B</Link>} />
                  <AuthProviderListComponent />
                </div>

                <div className="App-body">
                  <Route exact path="/" render={ () => <UserComponent /> }/>
                  <Route path="/config" exact={true} render={ () => <AuthConfigurationEditComponent /> }/>
                  <Route path="/config/authproviders/:name"
                         render={ ({match}) => <AuthProviderEditComponent name={match.params.name}/> }/>
                </div>
              </div>
            </div>

          </MuiThemeProvider>
        </div>
      </BrowserRouter>
    );
  }
}

export default App;
