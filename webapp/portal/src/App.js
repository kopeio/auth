import React, {Component} from 'react';
import UserComponent from "./components/UserComponent";
import {MuiThemeProvider} from "material-ui/styles";
import {BrowserRouter, Route} from "react-router-dom";

import './App.css';

class App extends Component {
  render() {
    return (
      <BrowserRouter>
        <div>
          <MuiThemeProvider>
            <div>
              <Route exact path="/" render={ () => <UserComponent /> }/>
            </div>
          </MuiThemeProvider>

        </div>
      </BrowserRouter>
    );
  }
}

export default App;
