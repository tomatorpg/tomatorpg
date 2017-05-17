import React, { Component } from 'react';
import PropTypes from 'prop-types';

class ChildHello extends Component {

  render() {
    const { message } = this.props;
    return (
      <span className="something">
        <strong>
          message: {message}
        </strong>
      </span>
    );
  }

}

ChildHello.propTypes = {
  message: PropTypes.string,
};

ChildHello.defaultProps = {
  message: 'no message',
};

export default ChildHello;
