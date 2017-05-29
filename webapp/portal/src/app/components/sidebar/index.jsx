// @flow

import {ActionHome, ActionSettings} from "material-ui/svg-icons/index";

import {
  AppBar, Divider, Drawer, ListItem
} from "material-ui";
import React from 'react';
import {Link} from "react-router-dom";


class Sidebar extends React.Component {
  render() {
    return <Drawer
      docked={false}
      onRequestChange={this.props.closeSidebar}
      open={this.props.open}
    >
      <AppBar
        onLeftIconButtonTouchTap={this.props.closeSidebar}
        title="Authentication"
      />
      <div className="App-sidebar">
        <ListItem primaryText="Home" leftIcon={<ActionHome />}
                  onClick={ this.props.closeSidebar }
                  containerElement={<Link to="/">A</Link>}/>
        <Divider />
        <ListItem primaryText="Settings" leftIcon={<ActionSettings />}
                  onClick={ this.props.closeSidebar }
                  containerElement={<Link to="/config" ></Link>}/>
      </div>
    </Drawer>;
  }

}

export default Sidebar;