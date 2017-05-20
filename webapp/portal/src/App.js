import React, {Component} from 'react';
import UserComponent from "./components/UserComponent";
import {MuiThemeProvider} from "material-ui/styles";
import {BrowserRouter, Route} from "react-router-dom";

import './App.css';
import AuthConfigurationEditComponent from "./components/AuthConfigurationEditComponent";
import AuthProviderEditComponent from "./components/AuthProviderEditComponent";

class App extends Component {
  render() {
    return (
      <BrowserRouter>
        <div>
          <MuiThemeProvider>
            <div>
              <Route exact path="/" render={ () => <UserComponent /> }/>
              <Route path="/config" exact={true} render={ () => <AuthConfigurationEditComponent /> }/>
              <Route path="/config/authproviders/:name" render={ ({match}) => <AuthProviderEditComponent name={match.params.name} /> }/>
            </div>
          </MuiThemeProvider>

        </div>
      </BrowserRouter>
    );
  }
}

export default App;
