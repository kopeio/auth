import React from 'react';
import ReactDOM from 'react-dom';
import {
  Card, CardActions, CardHeader, CardText, Dialog, FlatButton, MuiThemeProvider
} from "material-ui";
import AuthProviders from "../../api/AuthProviders";
import AuthConfigurationEditComponent from "./components/authconfigurationeditor/index";
import AuthProviderEditComponent from "./components/authprovidereditor/index";

import {ActionDelete, ContentAddCircle} from "material-ui/svg-icons/index";

// Technique for ConfigDialog based on http://blog.arkency.com/2015/04/beautiful-confirm-window-with-react/
class ConfirmDialog extends React.Component {
  handleClose = (result) => {
    this.promiseResolve(result);
  };

  componentDidMount() {
    let comp = this;
    this.promise = new Promise(function (resolve, reject) {
      comp.promiseResolve = resolve;
      //comp.promiseReject = reject;
    });
  };

  render() {
    const actions = [
      <FlatButton
        label="Cancel"
        primary={true}
        onTouchTap={() => this.handleClose(false) }
      />,
      <FlatButton
        label="Delete"
        primary={true}
        keyboardFocused={true}
        onTouchTap={() => this.handleClose(true) }
      />,
    ];

    return (
      <MuiThemeProvider>
        <Dialog
          title="Confirm deletion"
          actions={actions}
          modal={false}
          open={true}
          onRequestClose={() => this.handleClose(false)}
        >
          {this.props.message}
        </Dialog>
      </MuiThemeProvider>
    );
  }
}
;

function confirm(props = {}) {
  let wrapper = document.body.appendChild(document.createElement('div'));
  let component = ReactDOM.render(<ConfirmDialog {...props} />, wrapper);
  // let component = React.ren(React.createElement(ConfirmDialog, props), wrapper);
  // // cleanup = ->
  // // React.unmountComponentAtNode(wrapper)
  // setTimeout -> wrapper.remove()
  let cleanup = () => {
    ReactDOM.unmountComponentAtNode(wrapper);
    setTimeout(() => wrapper.remove());
  };
  component.promise.then(cleanup, cleanup);
  return component.promise;
};

class ProviderCard extends React.Component {
  render() {
    let provider = this.props.provider;

    var initiallyExpanded = false;

    var name = provider.metadata.name;
    if (!name) {
      name = "(new)";
      initiallyExpanded = true;
    }

    return <Card initiallyExpanded={initiallyExpanded}>
      <CardHeader
        title={name}
        subtitle={provider.description}
        actAsExpander={true}
        showExpandableButton={true}
      />
      <CardActions>
        <FlatButton label="Delete" icon={<ActionDelete />} onClick={this.props.onDelete}/>
      </CardActions>
      <CardText expandable={true}>
        <AuthProviderEditComponent name={provider.metadata.name}/>
      </CardText>
    </Card>
  }
}

function generateUUID() { // Public Domain/MIT
  var d = new Date().getTime();
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    var r = (d + Math.random() * 16) % 16 | 0;
    d = Math.floor(d / 16);
    return (c === 'x' ? r : ((r & 0x3) | 0x8)).toString(16);
  });
}

class Config extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      data: null,
    };
  }

  componentDidMount() {
    AuthProviders.build().list().then(json => {
      this.setState({
        data: json,
      });
    });
  }

  addProvider() {
    var data = Object.assign({}, this.state.data);
    data.items = [{
      metadata: {
        // We generate a uid so we can use it as the React key
        uid: generateUUID(),
      },
    }].concat(data.items);
    this.setState({
      data: data,
    });
  }

  deleteProvider(p) {
    var data = this.state.data;

    var i = data.items.findIndex((o) => (o === p));
    if (i === -1) {
      console.log("unable to find item for deletion: ", p);
      return;
    }

    if (!p.metadata.name) {
      // A new item
      data.items.splice(i, 1);
      this.setState({data: data});
    } else {
        confirm({message: "Delete authentication provider " + p.metadata.name + "?"}).then(
        (confirmed) => {
          if (confirmed) {
            AuthProviders.build().delete(p.metadata.name).then(() => {
              data.items.splice(i, 1);
              this.setState({data: data});
            });
          }
        });
    }
  }

  render() {
    var comp = this;
    if (!this.state.data) {
      return <div>Loading</div>;
    }

    return (
      <div>
        <AuthConfigurationEditComponent />

        <h2>Authentication Providers</h2>
        <FlatButton label="Add provider" icon={<ContentAddCircle />} onClick={() => {
          this.addProvider()
        }}/>
        {
          this.state.data.items.map(function (o) {
            return <ProviderCard key={o.metadata.uid} provider={o} onDelete={() => comp.deleteProvider(o) }/>
          })
        }
      </div>
    );
  }

}

export default Config;