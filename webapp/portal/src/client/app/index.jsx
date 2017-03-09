import React from 'react';
import {render} from 'react-dom';
import UserComponent from './components/UserComponent.jsx';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';

class App extends React.Component {
    render () {
        return (
            <MuiThemeProvider>
                <div>
                    <UserComponent user={this.props.user} />
                </div>
            </MuiThemeProvider>
        );
    }
}

render(<App user={window.initialProps.user}/>, document.getElementById('app'));