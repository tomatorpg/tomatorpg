/* eslint
react/forbid-prop-types: 'warn'
*/

import React, { Component } from 'react';
import PropTypes from 'prop-types';

class Notice extends Component {
    render() {
        return (
            <div className="notice-board">
            {this.props.children}
            </div>
        );
    }
}

Notice.propTypes = {
    visible: PropTypes.bool
};

Notice.defaultProps = {
    visible: true
};

export default Notice;
