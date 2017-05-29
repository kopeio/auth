// @flow

import {AppBar} from "material-ui";
import React from 'react';


class Header extends React.Component {
    render() {
        return <AppBar
            onLeftIconButtonTouchTap={this.props.openSidebar}
            title="Authentication"
        />;
    }

}

export default Header;