import React, {Component} from 'react';
import UserComponent from "./components/UserComponent";
import {MuiThemeProvider} from "material-ui/styles";
import {BrowserRouter, Route} from "react-router-dom";

import './App.css';
import Header from "./app/components/header/index";
import Sidebar from "./app/components/sidebar/index";
import Config from "./pages/config/index";


class App extends Component {
  constructor(props) {
    super(props);
    this.state = {showSidebar: false};
  };

  openSidebar() {
    this.setState({showSidebar: true});
  };

  closeSidebar() {
    this.setState({showSidebar: false});
  };

  render() {
    return (
      <BrowserRouter>
        <div>
          <MuiThemeProvider>
            <div>

              <Header openSidebar={() => {this.openSidebar()}}/>

              <div className="App-main">
                <Sidebar open={this.state.showSidebar} closeSidebar={() => {this.closeSidebar()}} />

                <div className="App-body">
                  <Route exact path="/" render={ () => <UserComponent /> }/>
                  <Route path="/config" exact={true} render={ () => <Config /> }/>
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
