/* eslint
react/forbid-prop-types: 'warn'
*/

import React, { Component } from 'react';
import PropTypes from 'prop-types';

class Notice extends Component {
  render() {
    return (this.props.visible) ? (
      <div id="notice">
        {this.props.children}
      </div>
    ) : (null);
  }
}

Notice.propTypes = {
  visible: PropTypes.bool,
  children: PropTypes.array,
};

Notice.defaultProps = {
  visible: true,
  children: [],
};

export default Notice;
