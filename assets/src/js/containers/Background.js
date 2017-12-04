/* eslint
react/forbid-prop-types: 'warn',
jsx-a11y/no-autofocus: 'warn',
*/

import React, { Component } from 'react';
import PropTypes from 'prop-types';

class Background extends Component {
  render() {
    return (
      <div id="background" style={this.props.style} />
    );
  }
}

Background.propTypes = {
  style: PropTypes.Object,
};

Background.defaultProps = {
  style: {
    backgroundImage: 'url("/assets/images/defaultbg.jpg")',
    minWidth: 800,
    minHeight: 500,
    backgroundSize: 'cover',
    backgroundRepeat: 'no-repeat',
    width: '100%',
  },
};

export default Background;
