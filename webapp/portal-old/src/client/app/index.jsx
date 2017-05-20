import React from 'react';
import {render} from 'react-dom';
import UserComponent from './components/UserComponent.jsx';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import ConfigComponent from "./components/ConfigComponent.jsx";
import {BrowserRouter, Route, Link} from 'react-router-dom'
import ConfigAuthProviderComponent from "./components/ConfigAuthProviderComponent.jsx";

const initialProps = window.initialProps;

class App extends React.Component {
  render() {
    return (
      <BrowserRouter>
        <div>
          <MuiThemeProvider>
            <div>
                <Route exact path="/" render={ () => <UserComponent user={initialProps.user} /> }/>
                <Route path="/config" exact={true} render={ () => <ConfigComponent config={initialProps.config} meta={initialProps.configmeta} /> }/>
                <Route path="/config/authproviders/:name" render={ ({match}) => <ConfigAuthProviderComponent name={match.params.name} /> }/>
            </div>
          </MuiThemeProvider>

        </div>
      </BrowserRouter>
    );
  }
}


render(<App />, document.getElementById('app'));